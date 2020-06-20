package server

import (
	"context"
	"github.com/gojuukaze/YTask/v2/backends"
	"github.com/gojuukaze/YTask/v2/brokers"
	"github.com/gojuukaze/YTask/v2/config"
	"github.com/gojuukaze/YTask/v2/controller"
	"github.com/gojuukaze/YTask/v2/log"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/worker"
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
)

type InlineServer struct {
	sync.Map
	groupName string

	workerMap map[string]worker.WorkerInterface // [workerName]worker

	broker  brokers.BrokerInterface
	backend backends.BackendInterface

	workerReadyChan chan struct{}
	msgChan         chan message.Message

	getMessageGoroutineStopChan chan struct{}
	workerGoroutineStopChan     chan struct{}

	safeStopChan chan struct{}

	// config
	StatusExpires int // second, -1:forever
	ResultExpires int // second, -1:forever
}

func NewInlineServer(groupName string, c config.Config) InlineServer {

	wm := make(map[string]worker.WorkerInterface)
	if c.Debug {
		log.YTaskLog.SetLevel(logrus.DebugLevel)
	}
	return InlineServer{
		groupName:     groupName,
		workerMap:     wm,
		broker:        c.Broker,
		backend:       c.Backend,
		safeStopChan:  make(chan struct{}),
		StatusExpires: c.StatusExpires,
		ResultExpires: c.ResultExpires,
	}
}

func (t *InlineServer) GetQueryName(name ...string) string {
	if len(name) > 0 {
		return "YTask:Query:" + name[0]
	}
	return "YTask:Query:" + t.groupName
}

func (t *InlineServer) MakeWorkerReady() {
	defer func() {
		recover()
	}()
	t.workerReadyChan <- struct{}{}
}

func (t *InlineServer) Run(numWorkers int) {

	_, ok := t.Load("isRunning")
	if ok {
		panic("inlineServer " + t.groupName + " is running")
	}
	t.Store("isRunning", struct{}{})

	if t.broker != nil {
		if t.broker.GetPoolSize() <= 0 {
			t.broker.SetPoolSize(1)
		}
		t.broker.Activate()
	}
	if t.backend != nil {
		if t.backend.GetPoolSize() <= 0 {
			t.backend.SetPoolSize(numWorkers)
		}
		t.backend.Activate()
	}

	log.YTaskLog.WithField("server", t.groupName).Infof("Start server[%s] numWorkers=%d", t.groupName, numWorkers)

	log.YTaskLog.WithField("server", t.groupName).Info("group workers:")
	for name := range t.workerMap {
		log.YTaskLog.WithField("server", t.groupName).Info("  - " + name)
	}

	t.workerReadyChan = make(chan struct{}, numWorkers)
	t.msgChan = make(chan message.Message, numWorkers)

	t.getMessageGoroutineStopChan = make(chan struct{}, 1)
	go t.GetNextMessageGoroutine()

	t.workerGoroutineStopChan = make(chan struct{}, 1)
	go t.WorkerGoroutine()

	for i := 0; i < numWorkers; i++ {
		t.MakeWorkerReady()
	}

}

func (t *InlineServer) safeStop() {
	log.YTaskLog.WithField("server", t.groupName).Info("waiting for incomplete tasks ")

	// stop get message goroutine
	close(t.workerReadyChan)
	t.Store("isStop", struct{}{})
	<-t.getMessageGoroutineStopChan

	// stop worker goroutine
	close(t.msgChan)
	<-t.workerGoroutineStopChan

	close(t.safeStopChan)

}

func (t *InlineServer) Shutdown(ctx context.Context) error {

	go func() {
		t.safeStop()
	}()

	select {
	case <-t.safeStopChan:
	case <-ctx.Done():
		return ctx.Err()
	}

	log.YTaskLog.WithField("server", t.groupName).Info("Shutdown!")
	return nil
}

func (t *InlineServer) IsRunning() bool {
	_, ok := t.Load("isRunning")
	return ok
}

// add worker to group
// w : worker func
func (t *InlineServer) Add(workerName string, w interface{}) {

	wType := reflect.TypeOf(w).Kind().String()
	if wType == "func" {
		funcWorker := worker.FuncWorker{
			Name: workerName,
			F:    w,
		}
		t.workerMap[workerName] = funcWorker
	} else {
		panic("worker must be func")
	}

}

func (t *InlineServer) Next() (message.Message, error) {
	return t.broker.Next(t.GetQueryName())
}

func (t *InlineServer) SetResult(result message.Result) error {
	var exTime int
	if result.IsFinish() {
		exTime = t.ResultExpires
	} else {
		exTime = t.StatusExpires
	}
	if exTime == 0 {
		return nil
	}
	return t.backend.SetResult(result, exTime)
}

// send msg to queue
// t.Send("groupName", "workerName" , 1,"hi",1.2)
//
func (t *InlineServer) Send(groupName string, workerName string, ctl controller.TaskCtl, args ...interface{}) (string, error) {
	var msg = message.NewMessage(ctl)
	msg.WorkerName = workerName
	err := msg.SetArgs(args...)
	if err != nil {
		return "", err
	}

	return msg.Id, t.broker.Send(t.GetQueryName(groupName), msg)

}

func (t *InlineServer) GetResult(id string) (message.Result, error) {
	result := message.NewResult(id)
	return t.backend.GetResult(result.GetBackendKey())
}

func (t *InlineServer) GetClient() Client {
	if t.broker.GetPoolSize() <= 0 {
			t.broker.SetPoolSize(10)
		}
		t.broker.Activate()

	if t.backend != nil {
		if t.backend.GetPoolSize() <= 0 {
			t.backend.SetPoolSize(10)
		}
		t.backend.Activate()
	}
	return Client{
		server: t,
		ctl:    controller.NewTaskCtl(),
	}
}
