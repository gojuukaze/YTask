package server

import (
	"github.com/gojuukaze/YTask/v2/backends"
	"github.com/gojuukaze/YTask/v2/brokers"
	"github.com/gojuukaze/YTask/v2/controller"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/yerrors"
)

// serverUtils用于把 delayServer，inlineServer，client 用到的方法抽离出来
type serverUtils struct {
	broker  brokers.BrokerInterface
	backend backends.BackendInterface

	// config
	statusExpires int // second, -1:forever
	resultExpires int // second, -1:forever
}

func newServerUtils(broker brokers.BrokerInterface, backend backends.BackendInterface, statusExpires int, resultExpires int) serverUtils {
	return serverUtils{broker: broker, backend: backend, statusExpires: statusExpires, resultExpires: resultExpires}
}

func (b *serverUtils) GetBrokerPoolSize() int {
	return b.broker.GetPoolSize()
}

func (b *serverUtils) SetBrokerPoolSize(num int) {
	b.broker.SetPoolSize(num)
}

func (b *serverUtils) BrokerActivate() {
	b.broker.Activate()
}

func (b serverUtils) GetQueueName(groupName string) string {
	// 这个key的名称拼错了，为了不影响已在运行的程序，只能这样了 = =
	return "YTask:Query:" + groupName
}

func (b *serverUtils) Next(groupName string) (message.Message, error) {
	return b.broker.Next(b.GetQueueName(groupName))
}

// send msg to queue
// t.Send("groupName", "workerName" , 1,"hi",1.2)
//
func (b *serverUtils) Send(groupName string, workerName string, ctl controller.TaskCtl, args ...interface{}) (string, error) {
	var msg = message.NewMessage(ctl)
	msg.WorkerName = workerName
	err := msg.SetArgs(args...)
	if err != nil {
		return "", err
	}

	return msg.Id, b.broker.Send(b.GetQueueName(groupName), msg)

}

func (b *serverUtils) GetBackendPoolSize() int {
	if b.backend == nil {
		return 0
	}
	return b.backend.GetPoolSize()
}

func (b *serverUtils) SetBackendPoolSize(num int) {
	if b.backend == nil {
		return
	}
	b.backend.SetPoolSize(num)
}

func (b *serverUtils) BackendActivate() {
	if b.backend == nil {
		return
	}
	b.backend.Activate()
}

func (b *serverUtils) SetResult(result message.Result) error {
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

func (b *serverUtils) GetResult(id string) (message.Result, error) {
	if b.backend == nil {
		return message.Result{}, yerrors.ErrNilResult{}
	}
	result := message.NewResult(id)
	return b.backend.GetResult(result.GetBackendKey())
}
