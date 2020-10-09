package ibmmq

import (
	"context"
	"fmt"
	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
	"strings"
)

type Consumer struct {
	Conn   *ibmmq.MQQueueManager
	Topic  string
	Header bool
}

var qObjectForC ibmmq.MQObject

func NewConsumer(conn *ibmmq.MQQueueManager, topic string, header bool) *Consumer {
	return &Consumer{conn, topic, header}
}

func (c *Consumer) Consume(ctx context.Context) {

	// Create the Object Descriptor that allows us to give the topic name
	mqsd := ibmmq.NewMQSD()
	mqsd.Options = ibmmq.MQSO_CREATE |
		ibmmq.MQSO_NON_DURABLE |
		ibmmq.MQSO_MANAGED

	mqsd.ObjectString = c.Topic

	// The qObject is filled in with a reference to the queue created automatically
	// for publications. It will be used in a moment for the Get operations
	_, err := c.Conn.Sub(mqsd, &qObjectForC)

	msgAvail := true
	for msgAvail == true && err == nil {
		var dataLen int

		// The GET requires control structures, the Message Descriptor (MQMD)
		// and Get Options (MQGMO). Create those with default values.
		getMQMD := ibmmq.NewMQMD()
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
		dataLen, err = qObjectForC.Get(getMQMD, gmo, buffer)

		if err != nil {
			msgAvail = false
			fmt.Println(err)
			mqReturn := err.(*ibmmq.MQReturn)
			if mqReturn.MQRC == ibmmq.MQRC_NO_MSG_AVAILABLE {
				// If there's no message available, then I won't treat that as a real error as
				// it's an expected situation
				err = nil
			}
		} else {
			fmt.Printf("Got message of length %d: ", dataLen)
			fmt.Println(strings.TrimSpace(string(buffer[:dataLen])))
		}
	}
}