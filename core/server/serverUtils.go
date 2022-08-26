package server

import (
	"github.com/gojuukaze/YTask/v3/core/backends"
	"github.com/gojuukaze/YTask/v3/core/brokers"
	"github.com/gojuukaze/YTask/v3/core/consts"
	"github.com/gojuukaze/YTask/v3/core/log"
	"github.com/gojuukaze/YTask/v3/core/message"
	"github.com/gojuukaze/YTask/v3/core/yerrors"
	"time"
)

// serverUtils用于把 delayServer，inlineServer，client 用到的方法抽离出来
type ServerUtils struct {
	broker  brokers.BrokerInterface
	backend backends.BackendInterface
	logger  log.LoggerInterface

	// config
	statusExpires int // second, -1:forever
	resultExpires int // second, -1:forever
}

func newServerUtils(broker brokers.BrokerInterface, backend backends.BackendInterface, logger log.LoggerInterface, statusExpires int, resultExpires int) ServerUtils {
	return ServerUtils{broker: broker, backend: backend, logger: logger, statusExpires: statusExpires, resultExpires: resultExpires}
}

func (b ServerUtils) GetQueueName(groupName string) string {
	if consts.UserV2Name {
		return "YTask:Query:" + groupName
	} else {
		return "YTask:Queue:" + groupName
	}
}

func (b ServerUtils) GetDelayGroupName(groupName string) string {
	return "Delay:" + groupName
}

func (b *ServerUtils) GetBrokerPoolSize() int {
	return b.broker.GetPoolSize()
}

func (b *ServerUtils) SetBrokerPoolSize(num int) {
	b.broker.SetPoolSize(num)
}

func (b *ServerUtils) BrokerActivate() {
	b.broker.Activate()
}

func (b *ServerUtils) Next(groupName string) (message.Message, error) {
	return b.broker.Next(b.GetQueueName(groupName))
}

// send msg to Queue
// t.Send("groupName", "workerName" , 1,"hi",1.2)
//
func (b *ServerUtils) Send(groupName string, workerName string, msgArgs message.MessageArgs, args ...interface{}) (string, error) {
	var msg = message.NewMessage(msgArgs)
	msg.WorkerName = workerName
	err := msg.SetArgs(args...)
	if err != nil {
		return "", err
	}
	return msg.Id, b.SendMsg(groupName, msg)

}

func (b *ServerUtils) SendMsg(groupName string, msg message.Message) error {
	var err error
	for i := 0; i < 3; i++ {
		err = b.broker.Send(b.GetQueueName(groupName), msg)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	return err

}

func (b *ServerUtils) LSendMsg(groupName string, msg message.Message) error {

	var err error
	for i := 0; i < 3; i++ {
		err = b.broker.LSend(b.GetQueueName(groupName), msg)
		if err == nil {
			break
		}
	}
	return err

}

func (b *ServerUtils) GetBackendPoolSize() int {
	if b.backend == nil {
		return 0
	}
	return b.backend.GetPoolSize()
}

func (b *ServerUtils) SetBackendPoolSize(num int) {
	if b.backend == nil {
		return
	}
	b.backend.SetPoolSize(num)
}

func (b *ServerUtils) BackendActivate() {
	if b.backend == nil {
		return
	}
	b.backend.Activate()
}

func (b *ServerUtils) SetResult(result message.Result) error {
	if b.backend == nil {
		return nil
	}
	var exTime int
	if result.IsFinish() {
		exTime = b.resultExpires
	} else {
		exTime = b.statusExpires
	}
	if exTime == 0 {
		return nil
	}
	return b.backend.SetResult(result, exTime)
}

func (b *ServerUtils) GetResult(id string) (message.Result, error) {
	if b.backend == nil {
		return message.Result{}, yerrors.ErrNilResult{}
	}
	result := message.NewResult(id)
	return b.backend.GetResult(result.GetBackendKey())
}

// - exTime : 过期时间，秒
func (b *ServerUtils) AbortTask(id string, exTime int) error {
	if b.backend == nil {
		return yerrors.ErrNilBackend{}
	}
	return b.backend.SetResult(message.NewAbortResult(id), exTime)
}

func (b *ServerUtils) IsAbort(id string) (bool, error) {
	if b.backend == nil {
		return false, yerrors.ErrNilBackend{}
	}
	_, err := b.backend.GetResult(message.NewAbortResult(id).GetBackendKey())
	if err == nil {
		return true, err
	}
	if yerrors.IsEqual(err, yerrors.ErrTypeNilResult) {
		return false, nil
	}
	return false, err
}
