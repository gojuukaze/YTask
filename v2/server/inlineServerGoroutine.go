package server

import (
	"sync"

	"github.com/vua/YTask/v2/log"
	"github.com/vua/YTask/v2/message"
	"github.com/vua/YTask/v2/worker"
	"github.com/vua/YTask/v2/yerrors"
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

			// 这里应该直接创建result就行，不知道之前为啥会从backend中获取
			result := message.NewResult(msg.Id)
			//result, err := t.GetResult(msg.Id)
			//log.YTaskLog.WithField("goroutine", "worker").
			//	Info("rrr= ", result, err)
			//if err != nil {
			//	if yerrors.IsEqual(err, yerrors.ErrTypeNilResult) {
			//		result = message.NewResult(msg.Id)
			//	} else {
			//		log.YTaskLog.WithField("goroutine", "worker").
			//			Error("get result error ", err)
			//		result = message.NewResult(msg.Id)
			//	}
			//}

			t.workerGoroutine_RunWorker(w, &msg, &result)

		}(msg)
	}

	waitWorkerWG.Wait()
	t.workerGoroutineStopChan <- struct{}{}
	log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "worker").Debug("stop")

}

func (t *InlineServer) workerGoroutine_RunWorker(w worker.WorkerInterface, msg *message.Message, result *message.Result) {

	ctl := msg.TaskCtl

RUN:

	result.SetStatusRunning()
	t.workerGoroutine_SaveResult(*result)

	err := w.Run(&ctl, msg.FuncArgs, result)
	if err == nil {
		result.Status = message.ResultStatus.Success
		t.workerGoroutine_SaveResult(*result)

		return
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

}

func (t *InlineServer) workerGoroutine_SaveResult(result message.Result) {
	log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "worker").
		Debugf("save result %+v", result)

	err := t.SetResult(result)
	if err != nil {
		log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "worker").Errorf("save result error ", err)
	}
}

func (t *InlineServer) GetDelayMessageGoroutine() {
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
				log.YTaskLog.WithField("server", t.groupName).
					WithField("goroutine", "get_next_message").
					Error("get msg error, ", err)
			}
			continue
		}
		log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "get_next_message").Infof("new msg %+v", msg)
		t.msgChan <- msg
	}

	t.getMessageGoroutineStopChan <- struct{}{}
	log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "get_next_message").Debug("stop")

}
