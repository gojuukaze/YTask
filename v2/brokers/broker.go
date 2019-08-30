package brokers

import (
	"github.com/gojuukaze/YTask/v2/message"
)

type BrokerInterface interface {
	Next(queryName string) (message.Message, error)
	Send(queryName string, msg message.Message) error
	// 调用Activate后才真正建立连接
	Activate()
	SetPoolSize(int)
	GetPoolSize()int
}
