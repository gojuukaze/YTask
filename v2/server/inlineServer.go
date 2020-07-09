package server

import (
	"context"
	"github.com/gojuukaze/YTask/v2/config"
	"github.com/gojuukaze/YTask/v2/log"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/util"
	"github.com/gojuukaze/YTask/v2/worker"
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
		groupName:     groupName,
		workerMap:     wm,
		serverUtils:   newServerUtils(c.Broker,c.Backend,c.StatusExpires,c.ResultExpires),
		safeStopChan:  make(chan struct{}),

	}
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
	t.Store("isRunning", struct{}{})

	t.SetBrokerPoolSize(1)
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

