package server

import (
	"github.com/gojuukaze/YTask/v2/log"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/worker"
	"github.com/gojuukaze/YTask/v2/yerrors"
	"sync"
)

// get next message if worker is ready
func (t *Server) GetNextMessageGoroutine(groupName string) {
	log.YTaskLog.WithField("goroutine", "GetNextMessage").Debug("start")
	var msg message.Message
	var err error
	for range t.workerReadyChan {
		_, ok := t.Load("isStop")
		if ok {
			break
		}
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
			defer func() {
				e := recover()
				if e != nil {
					log.YTaskLog.WithField("goroutine", "worker").Errorf("run worker[%s] panic %v", msg.WorkerName, e)
				}
			}()

			defer func() { go t.MakeWorkerReady() }()

			w, ok := workerMap[msg.WorkerName]
			if !ok {
				log.YTaskLog.WithField("goroutine", "worker").Error("not found worker", msg.WorkerName)
				return
			}

			waitWorkerWG.Add(1)
			defer waitWorkerWG.Done()

			result, err := t.GetResult(msg.Id)
			if err != nil {
				if yerrors.Compare(err, yerrors.ErrTypeNilResult) {
					result = message.NewResult(msg.Id)
				} else {
					log.YTaskLog.WithField("goroutine", "worker").
						Error("get result error ", err)
					result = message.NewResult(msg.Id)
				}
			}

			t.workerGoroutine_RunWorker(w, &msg, &result)

		}(msg)
	}

	waitWorkerWG.Wait()
	t.workerGoroutineStopChan <- struct{}{}
	log.YTaskLog.WithField("goroutine", "worker").Debug("stop")

}

func (t *Server) workerGoroutine_RunWorker(w worker.WorkerInterface, msg *message.Message, result *message.Result) {

	ctl := msg.TaskCtl

RUN:

	result.SetStatusRunning()
	t.workerGoroutine_SaveResult(*result)

	err := w.Run(&ctl, msg.JsonArgs, result)
	if err == nil {
		result.Status = message.ResultStatus.Success
		t.workerGoroutine_SaveResult(*result)

		return
	}
	log.YTaskLog.WithField("goroutine", "worker").Errorf("run worker[%s] error %s", msg.WorkerName, err)

	if ctl.CanRetry() {
		result.Status = message.ResultStatus.WaitingRetry
		ctl.RetryCount -= 1
		msg.TaskCtl = ctl
		log.YTaskLog.WithField("goroutine", "worker").Errorf("retry task %s", msg)
		ctl.SetError(nil)

		goto RUN
	} else {
		result.Status = message.ResultStatus.Failure
		t.workerGoroutine_SaveResult(*result)

	}

}

func (t *Server) workerGoroutine_SaveResult(result message.Result) {
	log.YTaskLog.WithField("goroutine", "worker").
		Debugf("save result %+v", result)
	err := t.SetResult(result)
	if err != nil {
		log.YTaskLog.WithField("goroutine", "worker").Errorf("save result error ", err)
	}
}
