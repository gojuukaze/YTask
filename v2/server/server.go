package server

import (
	"context"
	"fmt"
	"github.com/gojuukaze/YTask/v2/brokers"
	"github.com/gojuukaze/YTask/v2/config"
	yerrors "github.com/gojuukaze/YTask/v2/errors"
	"github.com/gojuukaze/YTask/v2/log"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/worker"
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
)

// [workerName]worker
type workerMap map[string]worker.WorkerInterface

type Server struct {
	workerGroup map[string]workerMap // [groupName]workerMap

	broker brokers.BrokerInterface

	workerReadyChan chan struct{}
	msgChan         chan message.Message

	getMessageGoroutineStopChan chan struct{}
	workerGoroutineStopChan     chan struct{}

	safeStopChan chan struct{}
}

func NewServer(c config.Config) Server {

	g := make(map[string]workerMap)
	if c.Debug {
		log.YTaskLog.SetLevel(logrus.DebugLevel)
	}
	return Server{
		workerGroup:  g,
		broker:       c.Broker,
		safeStopChan: make(chan struct{}),
	}
}

// get next message if worker is ready
func (t *Server) GetMessageGoroutine(groupName string) {
	log.YTaskLog.Debug("goroutine[GetMessage]: start")
	var msg message.Message
	var err error
	for range t.workerReadyChan {

		msg, err = t.Get(groupName)

		if err != nil {
			go t.MakeWorkerReady()
			ytaskErr, ok := err.(yerrors.YtaskError)
			if ok && ytaskErr.Type() != yerrors.ErrTypeEmptyQuery {
				log.YTaskLog.Error("goroutine[GetMessage]: get msg error, ", err)
			}
			continue
		}
		log.YTaskLog.Infof("goroutine[GetMessage]: New msg %+v", msg)
		t.msgChan <- msg
	}

	t.getMessageGoroutineStopChan <- struct{}{}
	log.YTaskLog.Debug("goroutine[GetMessage]: stop")

}

func (t *Server) MakeWorkerReady() {
	defer func() {
		recover()
	}()
	t.workerReadyChan <- struct{}{}
}

// start worker to run
func (t *Server) WorkerGoroutine(groupName string) {
	log.YTaskLog.Debug("goroutine[worker]: start")

	workerMap, _ := t.workerGroup[groupName]
	waitWorkerWG := sync.WaitGroup{}

	for msg := range t.msgChan {
		go func(msg message.Message) {
			defer func() {
				e := recover()
				if e != nil {
					log.YTaskLog.Errorf("Run worker[%s] panic %v", msg.WorkerName, e)
				}
			}()

			defer func() {
				go t.MakeWorkerReady()
			}()

			waitWorkerWG.Add(1)
			defer waitWorkerWG.Done()

			w, ok := workerMap[msg.WorkerName]
			if ok {
				err := w.Run(msg)
				if err != nil {
					log.YTaskLog.Errorf("goroutine[worker]: Run worker[%s] error %s", msg.WorkerName, err)
				}
			} else {
				log.YTaskLog.Error("goroutine[worker]: not found worker", msg.WorkerName)
			}
		}(msg)
	}

	waitWorkerWG.Wait()
	t.workerGoroutineStopChan <- struct{}{}
	log.YTaskLog.Debug("goroutine[worker]: stop")

}

func (t *Server) Run(groupName string, numWorkers int) {

	workerMap, ok := t.workerGroup[groupName]
	if !ok {
		panic("not find group '" + groupName + "'")
	}
	log.YTaskLog.Infof("Start group[%s] numWorkers=%d", groupName, numWorkers)

	log.YTaskLog.Info("group workers:")
	for name := range workerMap {
		log.YTaskLog.Info("  - " + name)
	}

	t.workerReadyChan = make(chan struct{}, numWorkers)
	t.msgChan = make(chan message.Message, numWorkers)

	t.getMessageGoroutineStopChan = make(chan struct{}, 1)
	go t.GetMessageGoroutine(groupName)
	t.workerGoroutineStopChan = make(chan struct{}, 1)

	go t.WorkerGoroutine(groupName)

	for i := 0; i < numWorkers; i++ {
		t.workerReadyChan <- struct{}{}
	}

}

func (t *Server) safeStop() {
	log.YTaskLog.Info("waiting for incomplete tasks ")

	// stop get message goroutine
	close(t.workerReadyChan)
	<-t.getMessageGoroutineStopChan

	// stop worker goroutine
	close(t.msgChan)
	<-t.workerGoroutineStopChan

	close(t.safeStopChan)

}

func (t *Server) Shutdown(ctx context.Context) error {

	go func() {
		t.safeStop()
	}()

	select {
	case <-t.safeStopChan:
	case <-ctx.Done():
		return ctx.Err()
	}

	log.YTaskLog.Info("Shutdown!")
	return nil
}

// add worker to group
// w : worker func
func (t *Server) Add(groupName string, workerName string, w interface{}) {
	_, ok := t.workerGroup[groupName]
	if !ok {
		t.workerGroup[groupName] = make(workerMap)
	}
	wType := reflect.TypeOf(w).Kind().String()
	if wType == "func" {
		funcWorker := worker.FuncWorker{
			Name: workerName,
			F:    w,
		}
		t.workerGroup[groupName][workerName] = funcWorker
	} else {
		panic(fmt.Sprintf("worker must be func"))
	}

}

func (t *Server) Get(groupName string) (message.Message, error) {
	return t.broker.Get(groupName)
}

// send msg to queue
// t.Send("groupName", "workerName" , 1,"hi",1.2)
//
func (t *Server) Send(groupName string, workerName string, args ...interface{}) error {
	var msg = message.Message{
		WorkerName: workerName,
	}
	err := msg.SetArgs(args...)
	if err != nil {
		return err
	}

	return t.broker.Send(groupName, msg)

}
