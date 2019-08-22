package brokers

import (
	"github.com/gojuukaze/YTask/v2/message"
)

type BrokerInterface interface {
	Get(queryName string) (message.Message, error)
	Send(queryName string, msg message.Message) error
}


