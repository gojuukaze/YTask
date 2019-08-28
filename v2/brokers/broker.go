package brokers

import (
	"github.com/gojuukaze/YTask/v2/message"
)

type BrokerInterface interface {
	Next(queryName string) (message.Message, error)
	Send(queryName string, msg message.Message) error
	Activate()
	SetPoolSize(int)
	GetPoolSize()int
}
