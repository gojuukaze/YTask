package server

import (
	"context"
	"github.com/gojuukaze/YTask/v3/config"
	"github.com/gojuukaze/YTask/v3/log"
	"github.com/gojuukaze/YTask/v3/message"
	"github.com/gojuukaze/YTask/v3/util"
	"github.com/gojuukaze/YTask/v3/worker"
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
)

type InlineServer struct {
	sync.Map
	serverUtils

	groupName string
	workerMap map[string]worker.WorkerInterface // [workerName]worker

	workerReadyChan chan struct{}
	msgChan         chan message.Message

	getMessageGoroutineStopChan chan struct{}
	workerGoroutineStopChan     chan struct{}

	safeStopChan chan struct{}
}

func NewInlineServer(groupName string, c config.Config) InlineServer {

	wm := make(map[string]worker.WorkerInterface)
	if c.Debug {
		log.YTaskLog.SetLevel(logrus.DebugLevel)
	}

	return InlineServer{
		groupName:                   groupName,
		workerMap:                   wm,
		serverUtils:                 newServerUtils(c.Broker, c.Backend, c.StatusExpires, c.ResultExpires),
		safeStopChan:                make(chan struct{}),
		getMessageGoroutineStopChan: make(chan struct{}),
		workerGoroutineStopChan:     make(chan struct{}),
	}
}

func (t *InlineServer) IsRunning() bool {
	_, ok := t.Load("isRunning")
	return ok
}

func (t *InlineServer) SetRunning() {
	t.Store("isRunning", struct{}{})

}

func (t *InlineServer) IsStop() bool {
	_, ok := t.Load("isStop")
	return ok
}

func (t *InlineServer) SetStop() {
	t.Store("isStop", struct{}{})

}

func (t *InlineServer) MakeWorkerReady() {
	defer func() {
		recover()
	}()
	t.workerReadyChan <- struct{}{}
}

func (t *InlineServer) Run(numWorkers int) {

	if t.IsRunning() {
		panic("inlineServer " + t.groupName + " is running")
	}
	t.SetRunning()

	// 初始化Broker, Backend
	if t.GetBrokerPoolSize() <= 0 {
		t.SetBrokerPoolSize(3)
	} else {
		t.SetBrokerPoolSize(t.GetBrokerPoolSize())
	}
	t.BrokerActivate()

	if t.backend != nil {
		if t.GetBackendPoolSize() <= 0 {
			t.SetBackendPoolSize(util.Min(10, numWorkers))
		}
		t.BackendActivate()
	}

	log.YTaskLog.WithField("server", t.groupName).Infof("Start server[%s] numWorkers=%d", t.groupName, numWorkers)

	log.YTaskLog.WithField("server", t.groupName).Info("group workers:")
	for name := range t.workerMap {
		log.YTaskLog.WithField("server", t.groupName).Info("  - " + name)
	}

	// 初始化chan，运行协程
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
	t.SetStop()
	<-t.getMessageGoroutineStopChan

	// stop worker goroutine
	close(t.msgChan)
	<-t.workerGoroutineStopChan

}

func (t *InlineServer) Shutdown(ctx context.Context) error {

	go func() {
		t.safeStop()
		close(t.safeStopChan)

	}()

	select {
	case <-t.safeStopChan:
	case <-ctx.Done():
		return ctx.Err()
	}

	log.YTaskLog.WithField("server", t.groupName).Info("Shutdown!")
	return nil
}

// Add worker to group
// w : worker func
// callbackFunc : callbackFunc func
func (t *InlineServer) Add(workerName string, w interface{}, callbackFunc ...interface{}) {

	var cFunc interface{} = nil

	cType := "func"
	if len(callbackFunc) > 0 {
		cFunc = callbackFunc[0]
		cType = reflect.TypeOf(cFunc).Kind().String()
	}

	wType := reflect.TypeOf(w).Kind().String()
	if wType == "func" && cType == "func" {
		funcWorker := worker.FuncWorker{
			Name:         workerName,
			Func:         w,
			CallbackFunc: cFunc,
		}
		t.workerMap[workerName] = funcWorker
	} else {
		panic("worker and callbackFunc must be func")
	}

}
