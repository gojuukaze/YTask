package server

import (
	"github.com/gojuukaze/YTask/v3/log"
	"github.com/gojuukaze/YTask/v3/message"
	"github.com/gojuukaze/YTask/v3/worker"
	"github.com/gojuukaze/YTask/v3/yerrors"
	"sync"
)

// get next message if worker is ready
func (t *InlineServer) GetNextMessageGoroutine() {
	log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "get_next_message").Debug("start")
	var msg message.Message
	var err error
	for range t.workerReadyChan {
		if t.IsStop() {
			break
		}
		msg, err = t.Next(t.groupName)

		if err != nil {
			go t.MakeWorkerReady()
			if !yerrors.IsEqual(err, yerrors.ErrTypeEmptyQuery) {
				log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "get_next_message").Error("get msg error, ", err)
			}
			continue
		}
		log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "get_next_message").Infof("new msg %+v", msg)
		t.msgChan <- msg
	}

	t.getMessageGoroutineStopChan <- struct{}{}
	log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "get_next_message").Debug("stop")

}

// start worker to run
func (t *InlineServer) WorkerGoroutine() {
	log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "worker").Debug("start")

	waitWorkerWG := sync.WaitGroup{}

	for msg := range t.msgChan {
		go func(msg message.Message) {
			defer func() {
				e := recover()
				if e != nil {
					log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "worker").Errorf("run worker[%s] panic %v", msg.WorkerName, e)
				}
			}()

			defer func() { go t.MakeWorkerReady() }()

			w, ok := t.workerMap[msg.WorkerName]
			if !ok {
				log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "worker").Error("not found worker", msg.WorkerName)
				return
			}

			waitWorkerWG.Add(1)
			defer waitWorkerWG.Done()

			result := message.NewResult(msg.Id)

			t.workerGoroutine_RunWorker(w, &msg, &result)

		}(msg)
	}

	waitWorkerWG.Wait()
	t.workerGoroutineStopChan <- struct{}{}
	log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "worker").Debug("stop")

}

func (t *InlineServer) workerGoroutine_RunWorker(w worker.WorkerInterface, msg *message.Message, result *message.Result) {

	var err error
	ctl := msg.TaskCtl

RUN:

	if ctl.IsExpired() {
		result.Status = message.ResultStatus.Expired
		t.workerGoroutine_SaveResult(*result)
		goto AFTER
	}

	result.SetStatusRunning()
	t.workerGoroutine_SaveResult(*result)

	err = w.Run(&ctl, msg.FuncArgs, result)

	if err == nil {
		result.Status = message.ResultStatus.Success
		t.workerGoroutine_SaveResult(*result)
		goto AFTER
	}
	log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "worker").Errorf("run worker[%s] error %s", msg.WorkerName, err)

	if ctl.CanRetry() {
		result.Status = message.ResultStatus.WaitingRetry
		ctl.RetryCount -= 1
		msg.TaskCtl = ctl
		log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "worker").Infof("retry task %s", msg)
		ctl.SetError(nil)

		goto RUN
	} else {
		result.Status = message.ResultStatus.Failure
		t.workerGoroutine_SaveResult(*result)

	}

AFTER:

	err = w.After(&ctl, msg.FuncArgs, result)
	if err != nil {
		log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "worker").Errorf("run worker[%s] callback error %s", msg.WorkerName, err)
	}

}

func (t *InlineServer) workerGoroutine_SaveResult(result message.Result) {
	log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "worker").
		Debugf("save result %+v", result)

	err := t.SetResult(result)
	if err != nil {
		log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "worker").Errorf("save result error ", err)
	}
}
