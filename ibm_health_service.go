package ibmmq

import (
	"context"
	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
	"time"
)

type IBMHealthService struct {
	name       string
	connection *ibmmq.MQQueueManager
	topic      string
	timeout    time.Duration
}

var qObject ibmmq.MQObject

func NewDefaultIBMHealthService(connection *ibmmq.MQQueueManager, topic string) *IBMHealthService {
	return &IBMHealthService{"IBM", connection, topic, 5 * time.Second}
}

func NewIBMHealthService(name string, connection *ibmmq.MQQueueManager, topic string, timeout time.Duration) *IBMHealthService {
	return &IBMHealthService{name, connection, topic, timeout}
}

func (s *IBMHealthService) Name() string {
	return s.name
}

func (s *IBMHealthService) Check(ctx context.Context) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	mqsd := ibmmq.NewMQSD()
	mqsd.Options = ibmmq.MQSO_CREATE |
		ibmmq.MQSO_NON_DURABLE |
		ibmmq.MQSO_MANAGED
	mqsd.ObjectString = s.topic
	subscriptionObject, err := s.connection.Sub(mqsd, &qObject)
	if err != nil {
		return nil, err
	}
	err = subscriptionObject.Close(0)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *IBMHealthService) Build(ctx context.Context, data map[string]interface{}, err error) map[string]interface{} {
	if err == nil {
		return data
	}
	data["error"] = err.Error()
	return data
}
