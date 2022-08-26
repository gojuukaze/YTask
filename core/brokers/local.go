package brokers

import (
	"github.com/gojuukaze/YTask/v3/core/drive"
	"github.com/gojuukaze/YTask/v3/core/message"
	"github.com/gojuukaze/YTask/v3/core/util/yjson"
	"github.com/gojuukaze/YTask/v3/core/yerrors"
)

// LocalBroker
// ！！！只用于本地测试！！！
// !!! Only for local test !!!
type LocalBroker struct {
	client drive.LocalDrive
}

func NewLocalBroker() LocalBroker {
	return LocalBroker{}
}
func (l *LocalBroker) Activate() {
	l.client = drive.NewLocalDrive(true)
}
func (l *LocalBroker) Next(queueName string) (message.Message, error) {
	var msg message.Message
	b, err := l.client.LPop(queueName)
	if err != nil {
		if err == drive.EmptyQueueError {
			return msg, yerrors.ErrEmptyQueue{}
		}
		return msg, err
	}
	err = yjson.YJson.Unmarshal(b, &msg)
	return msg, err
}

func (l *LocalBroker) Send(queueName string, msg message.Message) error {
	b, err := yjson.YJson.Marshal(msg)

	if err != nil {
		return err
	}
	err = l.client.RPush(queueName, b)
	return err
}

func (l *LocalBroker) LSend(queueName string, msg message.Message) error {
	b, err := yjson.YJson.Marshal(msg)

	if err != nil {
		return err
	}
	err = l.client.LPush(queueName, b)
	return err
}

func (l *LocalBroker) SetPoolSize(i int) {

}

func (l *LocalBroker) GetPoolSize() int {
	return 0
}

func (l *LocalBroker) Clone() BrokerInterface {
	return &LocalBroker{}
}
