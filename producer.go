package ibmmq

import (
	"context"
	"fmt"
	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

type Producer struct {
	Conn *ibmmq.MQQueueManager

	Topic string
}

func NewProducer(conn *ibmmq.MQQueueManager, topic string) *Producer {
	return &Producer{conn, topic}
}

var qObjectForP ibmmq.MQObject

func (p *Producer) Produce(ctx context.Context, data []byte) (string, error) {

	openOptions := ibmmq.MQOO_OUTPUT
	mqod := ibmmq.NewMQOD()
	mqod.ObjectType = ibmmq.MQOT_TOPIC
	mqod.ObjectString = p.Topic

	topicObject, err := p.Conn.Open(mqod, openOptions)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Publishing message to", p.Topic)
	}
	putMQMD := ibmmq.NewMQMD()
	pmo := ibmmq.NewMQPMO()

	// The default options are OK, but it's always
	// a good idea to be explicit about transactional boundaries as
	// not all platforms behave the same way.
	pmo.Options = ibmmq.MQPMO_NO_SYNCPOINT

	// Tell MQ what the message body format is. In this case, a text string
	putMQMD.Format = ibmmq.MQFMT_STRING

	// Now put the message to the queue
	err = topicObject.Put(putMQMD, pmo, data)
	return "", err
}
