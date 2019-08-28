package server

import (
	"github.com/gojuukaze/YTask/v2/log"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/yerrors"
	"sync"
)

// get next message if worker is ready
func (t *Server) GetNextMessageGoroutine(groupName string) {
	log.YTaskLog.WithField("goroutine", "GetNextMessage").Debug("start")
	var msg message.Message
	var err error
	for range t.workerReadyChan {

		msg, err = t.Next(groupName)

		if err != nil {
			go t.MakeWorkerReady()
			ytaskErr, ok := err.(yerrors.YTaskError)
			if ok && ytaskErr.Type() != yerrors.ErrTypeEmptyQuery {
				log.YTaskLog.WithField("goroutine", "GetNextMessage").Error("get msg error, ", err)
			}
			continue
		}
		log.YTaskLog.WithField("goroutine", "GetNextMessage").Infof("new msg %+v", msg)
		t.msgChan <- msg
	}

	t.getMessageGoroutineStopChan <- struct{}{}
	log.YTaskLog.WithField("goroutine", "GetNextMessage").Debug("stop")

}

// start worker to run
func (t *Server) WorkerGoroutine(groupName string) {
	log.YTaskLog.WithField("goroutine", "worker").Debug("start")

	workerMap, _ := t.workerGroup[groupName]
	waitWorkerWG := sync.WaitGroup{}

	for msg := range t.msgChan {
		go func(msg message.Message) {
			var result = message.NewResult(msg.Id)
			result.SetStatusRunning()
			//t.resultChan <- result

			log.YTaskLog.WithField("goroutine", "worker").
				Debugf("save result %+v", result)
			err := t.SetResult(result)
			if err != nil {
				log.YTaskLog.WithField("goroutine", "worker").Errorf("save result error ", err)
			}

			defer func() {
				log.YTaskLog.WithField("goroutine", "worker").
					Debugf("save result %+v", result)
				err := t.SetResult(result)
				if err != nil {
					log.YTaskLog.WithField("goroutine", "worker").Errorf("save result error ", err)
				}
			}()
			defer func() {
				e := recover()
				if e != nil {
					log.YTaskLog.WithField("goroutine", "worker").Errorf("run worker[%s] panic %v", msg.WorkerName, e)
				}
			}()

			defer func() {
				go t.MakeWorkerReady()
			}()

			waitWorkerWG.Add(1)
			defer waitWorkerWG.Done()

			w, ok := workerMap[msg.WorkerName]
			if ok {
				err := w.Run(msg, &result)
				if err != nil {
					log.YTaskLog.WithField("goroutine", "worker").Errorf("run worker[%s] error %s", msg.WorkerName, err)
				}
			} else {
				log.YTaskLog.WithField("goroutine", "worker").Error("not found worker", msg.WorkerName)
			}
		}(msg)
	}

	waitWorkerWG.Wait()
	t.workerGoroutineStopChan <- struct{}{}
	log.YTaskLog.WithField("goroutine", "worker").Debug("stop")

}

// save result
func (t *Server) SaveResultGoroutine() {

	log.YTaskLog.WithField("goroutine", "SaveResult").Debug("start")

	wg := sync.WaitGroup{}
	if t.backend != nil {
		for r := range t.resultChan {
			wg.Add(1)
			go func(result message.Result) {
				defer wg.Done()
				log.YTaskLog.WithField("goroutine", "SaveResult").
					Debugf("save result %+v", result)
				err := t.SetResult(result)
				if err != nil {
					log.YTaskLog.WithField("goroutine", "SaveResult").Errorf("save result error ", err)
				}

			}(r)
		}
	} else {
		log.YTaskLog.WithField("goroutine", "SaveResult").Debug("backend is nil")
		for range t.resultChan {
		}
	}

	wg.Wait()
	t.saveResultGoRoutineStopChan <- struct{}{}
	log.YTaskLog.WithField("goroutine", "SaveResult").Debug("stop")
}
