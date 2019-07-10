package ytask

import (
	"YTask/v1/brokers"
	"YTask/v1/config"
	"YTask/v1/yerrors"
	"YTask/v1/ylog"
	"YTask/v1/ymsg"
	"YTask/v1/yworker"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
)

// [workerName]woker
type workerMap map[string]yworker.WorkerInterface

type YTask struct {
	workerGroup map[string]workerMap // [groupName]workerMap

	broker brokers.BrokerInterface

	cancel       context.CancelFunc
	wg           sync.WaitGroup
	shutdownFlag bool
	waitFlag     chan struct{}
	lock         sync.RWMutex
}

func NewYTask(c config.Config) YTask {

	g := make(map[string]workerMap)
	if c.Debug {
		ylog.YTaskLog.SetLevel(logrus.DebugLevel)
	}
	return YTask{
		workerGroup: g,
		broker:      c.Broker,
		waitFlag:    make(chan struct{}),
	}
}

func (t *YTask) Run(groupName string, numWorkers int) {

	workerMap, ok := t.workerGroup[groupName]
	if !ok {
		panic("not find group '" + groupName + "'")
	}
	ylog.YTaskLog.Infof("Start group[%s]", groupName)

	ylog.YTaskLog.Info("group workers:")
	for name := range workerMap {
		ylog.YTaskLog.Info("  - " + name)
	}

	t.wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func(workerId int) {
			defer t.wg.Done()

			for {
				t.lock.RLock()
				if t.shutdownFlag {
					t.lock.RUnlock()
					ylog.YTaskLog.Debugf("worker[%d] stop", workerId)
					return
				}
				t.lock.RUnlock()
				msg, err := t.Get(groupName)

				if err != nil {
					if err == yerrors.ErrEmptyQuery {
						continue
					}
					ylog.YTaskLog.Error("get ymsg error: ", err)
				}
				ylog.YTaskLog.Infof("New ymsg %+v", msg)
				worker, ok := workerMap[msg.WorkerName]
				if ok {
					err = worker.Run(msg)
					if err != nil {
						ylog.YTaskLog.Errorf("Run yworker[%s] error %s", msg.WorkerName, err)
					}
				} else {
					ylog.YTaskLog.Error("not found yworker", msg.WorkerName)
				}

			}
		}(i)
	}
}

func (t *YTask) waitTask() {
	ylog.YTaskLog.Info("waiting for incomplete tasks ")
	t.wg.Wait()

	close(t.waitFlag)
	fmt.Println(111)

}

func (t *YTask) Shutdown(ctx context.Context) error {

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

	ylog.YTaskLog.Info("Shutdown!")
	return nil

}

func (t *YTask) Add(group string, w yworker.WorkerInterface) {

	_, ok := t.workerGroup[group]
	if !ok {
		t.workerGroup[group] = make(workerMap)
	}
	t.workerGroup[group][w.Name()] = w
}

func (t *YTask) Get(groupName string) (ymsg.Message, error) {
	return t.broker.Get(groupName)
}

// send task
//
// t.Send("xx", ymsg.Message{} )
//
// t.Send("xx", "workerName1" , User{1,"name"} )
//
// t.Send("xx", "workerName1" , "userJsonString" )
//
func (t *YTask) Send(groupName string, values ...interface{}) error {
	if len(values) == 1 {
		msg, ok := values[0].(ymsg.Message)
		if ok {
			return t.broker.Send(groupName, msg)
		}
		return errors.New("values must be ymsg.Message")
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

		return t.broker.Send(groupName, ymsg.Message{wName, jsonArgs})
	}

	return errors.New("too many values")

}
