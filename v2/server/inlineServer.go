package server

import (
	"context"
	"fmt"
	"github.com/gojuukaze/YTask/v2/config"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/util"
	"github.com/gojuukaze/YTask/v2/worker"
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
		//log.YTaskLog.SetLevel(logrus.DebugLevel)
		c.Logger.SetLevel("debug")
	}

	return InlineServer{
		groupName:                   groupName,
		workerMap:                   wm,
		serverUtils:                 newServerUtils(c.Broker, c.Backend, c.Logger, c.StatusExpires, c.ResultExpires),
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

	t.SetBrokerPoolSize(1)
	t.BrokerActivate()

	if t.backend != nil {
		if t.GetBackendPoolSize() <= 0 {
			t.SetBackendPoolSize(util.Min(10, numWorkers))
		}
		t.BackendActivate()
	}

	//log.YTaskLog.WithField("server", t.groupName).Infof("Start server[%s] numWorkers=%d", t.groupName, numWorkers)
	t.logger.InfoWithField(fmt.Sprintf("Start server[%s] numWorkers=%d", t.groupName, numWorkers), "server", t.groupName)

	//log.YTaskLog.WithField("server", t.groupName).Info("group workers:")
	t.logger.InfoWithField("group workers:", "server", t.groupName)

	for name := range t.workerMap {
		//log.YTaskLog.WithField("server", t.groupName).Info("  - " + name)
		t.logger.InfoWithField("  - " + name, "server", t.groupName)
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
	//log.YTaskLog.WithField("server", t.groupName).Info("waiting for incomplete tasks ")
	t.logger.InfoWithField("waiting for incomplete tasks ", "server", t.groupName)

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

	//log.YTaskLog.WithField("server", t.groupName).Info("Shutdown!")
	t.logger.InfoWithField("Shutdown!", "server", t.groupName)
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
		funcWorker := &worker.FuncWorker{
			Name:         workerName,
			Func:         w,
			CallbackFunc: cFunc,
			Logger:       t.logger,
		}
		t.workerMap[workerName] = funcWorker
	} else {
		panic("worker and callbackFunc must be func")
	}
}
