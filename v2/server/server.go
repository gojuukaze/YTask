package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gojuukaze/YTask/v1.1/brokers"
	"github.com/gojuukaze/YTask/v1.1/config"
	yerrors "github.com/gojuukaze/YTask/v1.1/errors"
	"github.com/gojuukaze/YTask/v1.1/log"
	"github.com/gojuukaze/YTask/v1.1/message"
	"github.com/gojuukaze/YTask/v1.1/worker"
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
)

// [workerName]woker
type workerMap map[string]worker.WorkerInterface

type Server struct {
	workerGroup map[string]workerMap // [groupName]workerMap

	broker brokers.BrokerInterface

	wg           sync.WaitGroup
	shutdownFlag bool
	waitFlag     chan struct{}
	lock         sync.RWMutex
}

func NewServer(c config.Config) Server {

	g := make(map[string]workerMap)
	if c.Debug {
		log.YTaskLog.SetLevel(logrus.DebugLevel)
	}
	return Server{
		workerGroup: g,
		broker:      c.Broker,
		waitFlag:    make(chan struct{}),
	}
}

func (t *Server) Run(groupName string, numWorkers int) {

	workerMap, ok := t.workerGroup[groupName]
	if !ok {
		panic("not find group '" + groupName + "'")
	}
	log.YTaskLog.Infof("Start group[%s]", groupName)

	log.YTaskLog.Info("group workers:")
	for name := range workerMap {
		log.YTaskLog.Info("  - " + name)
	}

	t.wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func(workerId int) {
			defer t.wg.Done()

			for {
				t.lock.RLock()
				if t.shutdownFlag {
					t.lock.RUnlock()
					log.YTaskLog.Debugf("worker[%d] stop", workerId)
					return
				}
				t.lock.RUnlock()
				msg, err := t.Get(groupName)

				if err != nil {
					if err == yerrors.ErrEmptyQuery {
						continue
					}
					log.YTaskLog.Error("get msg error: ", err)
				}
				log.YTaskLog.Infof("New msg %+v", msg)
				worker, ok := workerMap[msg.WorkerName]
				if ok {
					err = worker.Run(msg)
					if err != nil {
						log.YTaskLog.Errorf("Run worker[%s] error %s", msg.WorkerName, err)
					}
				} else {
					log.YTaskLog.Error("not found worker", msg.WorkerName)
				}

			}
		}(i)
	}
}

func (t *Server) waitTask() {
	log.YTaskLog.Info("waiting for incomplete tasks ")
	t.wg.Wait()

	close(t.waitFlag)

}

func (t *Server) Shutdown(ctx context.Context) error {

	t.lock.Lock()
	t.shutdownFlag = true
	t.lock.Unlock()

	go func() {
		t.waitTask()
	}()

	select {
	case <-t.waitFlag:
	case <-ctx.Done():
		return ctx.Err()

	}

	log.YTaskLog.Info("Shutdown!")
	return nil

}

// add worker to group
// worker : func or struct
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
		if wType == "struct" || wType == "ptr" {
			structWorker := worker.StructWorker{
				S:    w,
				Name: workerName,
			}
			t.workerGroup[groupName][workerName] = structWorker

		} else {
			s := fmt.Sprintf("worker must be func, struct, *struct")
			panic(s)
		}
	}

}

func (t *Server) Get(groupName string) (message.Message, error) {
	return t.broker.Get(groupName)
}

// send msg to queue
//
// t.Send("xx", message.Message{} )
//
// t.Send("xx", "workerName1" , User{1,"name"} )
//
// t.Send("xx", "workerName1" , `{"id":1,"name":"xx"}` )
//
func (t *Server) Send(groupName string, values ...interface{}) error {
	if len(values) == 1 {
		msg, ok := values[0].(message.Message)
		if ok {
			return t.broker.Send(groupName, msg)
		}
		return errors.New("values must be msg.Message")
	}
	if len(values) == 2 {
		wName, ok := values[0].(string)
		if !ok {
			return errors.New("values[0] must be string")
		}

		var jsonArgs string
		switch reflect.TypeOf(values[1]).Kind().String() {
		case "string":
			jsonArgs = values[1].(string)
		case "struct":
			b, err := json.Marshal(values[1])
			if err != nil {
				return nil
			}
			jsonArgs = string(b)
		}

		return t.broker.Send(groupName, message.Message{wName, jsonArgs})
	}

	return errors.New("too many values")

}
