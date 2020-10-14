package ibmmq

import (
	"context"
	"github.com/common-go/mq"
	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

type Consumer struct {
	QueueManager *ibmmq.MQQueueManager
	QueueName    string
}

var qObjectForC ibmmq.MQObject

func NewConsumer(manager *ibmmq.MQQueueManager, queueName string) *Consumer {
	return &Consumer{manager, queueName}
}
func NewConsumerByConfig(c QueueConfig, auth MQAuth) (*Consumer, error) {
	mgr, err := NewQueueManagerByConfig(c, auth)
	if err != nil {
		return nil, err
	}
	return &Consumer{QueueManager: mgr, QueueName: c.QueueName}, nil
}
func (c *Consumer) Consume(ctx context.Context, caller mq.ConsumerCaller) {
	// Create the Object Descriptor that allows us to give the topic name
	sd := ibmmq.NewMQSD()
	sd.Options = ibmmq.MQSO_CREATE |
		ibmmq.MQSO_NON_DURABLE |
		ibmmq.MQSO_MANAGED

	sd.ObjectString = c.QueueName

	// The qObject is filled in with a reference to the queue created automatically
	// for publications. It will be used in a moment for the Get operations
	_, er0 := c.QueueManager.Sub(sd, &qObjectForC)

	msgAvail := true
	for msgAvail == true && er0 == nil {
		// The GET requires control structures, the Message Descriptor (MQMD)
		// and Get Options (MQGMO). Create those with default values.
		md := ibmmq.NewMQMD()
		gmo := ibmmq.NewMQGMO()

		// The default options are OK, but it's always
		// a good idea to be explicit about transactional boundaries as
		// not all platforms behave the same way.
		gmo.Options = ibmmq.MQGMO_NO_SYNCPOINT
		// Set options to wait for a maximum of 3 seconds for any new message to arrive
		gmo.Options |= ibmmq.MQGMO_WAIT
		gmo.WaitInterval = 3 * 1000 // The WaitInterval is in milliseconds

		// Create a buffer for the message data. This one is large enough
		// for the messages put by the amqsput sample.
		buffer := make([]byte, 1024)
		_, er1 := qObjectForC.Get(md, gmo, buffer)

		if er1 != nil {
			msgAvail = false
			mqReturn := er1.(*ibmmq.MQReturn)
			if mqReturn.MQRC != ibmmq.MQRC_NO_MSG_AVAILABLE {
				caller.Call(ctx, nil, er1)
			} else {
				// If there's no message available, then I won't treat that as a real error as
				// it's an expected situation
				er1 = nil
			}
		} else {
			msg := mq.Message{
				Data:       buffer,
			}
			caller.Call(ctx, &msg, nil)
		}
	}
}
