package brokers

import (
	"github.com/gojuukaze/YTask/v1.1/message"
)

type BrokerInterface interface {
	Get(queryName string) (message.Message, error)
	Send(queryName string, msg message.Message) error
}


