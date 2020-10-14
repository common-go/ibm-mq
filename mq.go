package ibmmq

import "github.com/ibm-messaging/mq-golang/ibmmq"

type MQAuth struct {
	UserId   string `mapstructure:"user_id"`
	Password string `mapstructure:"password"`
}

func NewMQCDByChannelAndConnection(channelName string, connectionName string) *ibmmq.MQCD {
	cd := ibmmq.NewMQCD()
	cd.ChannelName = channelName
	cd.ConnectionName = connectionName
	return cd
}
func NewMQCSPByConfig(auth MQAuth) *ibmmq.MQCSP {
	csp := ibmmq.NewMQCSP()
	csp.AuthenticationType = ibmmq.MQCSP_AUTH_USER_ID_AND_PWD
	csp.UserId = auth.UserId
	csp.Password = auth.Password
	return csp
}

type QueueConfig struct {
	ManagerName    string `mapstructure:"manager_name"`
	ChannelName    string `mapstructure:"channel_name"`
	ConnectionName string `mapstructure:"connection_name"`
	QueueName      string `mapstructure:"queue_name"`
}

func NewQueueManagerByConfig(c QueueConfig, auth MQAuth) (*ibmmq.MQQueueManager, error) {
	cd := NewMQCDByChannelAndConnection(c.ChannelName, c.ConnectionName)
	csp := NewMQCSPByConfig(auth)

	cno := ibmmq.NewMQCNO()
	cno.ClientConn = cd
	cno.Options = ibmmq.MQCNO_CLIENT_BINDING + ibmmq.MQCNO_RECONNECT + ibmmq.MQCNO_HANDLE_SHARE_BLOCK
	cno.SecurityParms = csp

	mgr, err := ibmmq.Connx(c.ManagerName, cno)
	return &mgr, err
}
