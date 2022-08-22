package server

import (
	"github.com/gojuukaze/YTask/v3/log"
	"github.com/gojuukaze/YTask/v3/message"
	"github.com/gojuukaze/YTask/v3/worker"
	"github.com/gojuukaze/YTask/v3/yerrors"
	"sync"
	"time"
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

// return : current Workflow index
func (t *InlineServer) workerGoroutine_UpdateWorkflowResult(msg *message.Message, result *message.Result) int {
	workflowIndex := 0
	result.Workflow = make([][2]string, len(msg.TaskCtl.Workflow))
	for i, w := range msg.TaskCtl.Workflow {
		if w.WorkerName == msg.WorkerName {
			workflowIndex = i
		}
		result.Workflow[i] = [2]string{w.WorkerName, message.WorkflowStatus.Waiting}
	}
	for i := 0; i < workflowIndex; i++ {
		result.Workflow[i][1] = message.WorkflowStatus.Success
	}
	return workflowIndex
}

func (t *InlineServer) workerGoroutine_RunWorker(w worker.WorkerInterface, msg *message.Message, result *message.Result) {

	var err error
	ctl := msg.TaskCtl
	workflowIndex := -1
	if len(ctl.Workflow) > 0 {
		workflowIndex = t.workerGoroutine_UpdateWorkflowResult(msg, result)
	}

RUN:

	if ctl.IsExpired() {
		t.workerGoroutine_UpdateResultStatus(message.ResultStatus.Expired, workflowIndex, result)
		t.workerGoroutine_SaveResult(*result)
		goto AFTER
	}

	result.SetStatusRunning()
	t.workerGoroutine_UpdateResultStatus(result.Status, workflowIndex, result)
	t.workerGoroutine_SaveResult(*result)

	err = w.Run(&ctl, msg.FuncArgs, result)

	if err == nil {
		t.workerGoroutine_UpdateResultStatus(message.ResultStatus.Success, workflowIndex, result)
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
		t.workerGoroutine_UpdateResultStatus(message.ResultStatus.Failure, workflowIndex, result)
		t.workerGoroutine_SaveResult(*result)
	}

AFTER:
	// 为了逻辑更简单，工作流和回调暂不兼容
	if workflowIndex >= 0 && workflowIndex+1 < len(ctl.Workflow) {
		t.workerGoroutine_NextWorkflow(workflowIndex+1, *msg, *result)
	} else {
		err = w.After(&ctl, msg.FuncArgs, result)
		if err != nil {
			log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "worker").Errorf("run worker[%s] callback error %s", msg.WorkerName, err)
		}
	}
}

func (t *InlineServer) workerGoroutine_UpdateResultStatus(status int, workflowIndex int, result *message.Result) {
	result.Status = status
	if workflowIndex >= 0 {
		result.Workflow[workflowIndex][1] = message.StatusToWorkflowStatus[status]
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

func (t *InlineServer) workerGoroutine_NextWorkflow(nextIndex int, msg message.Message, result message.Result) {

	next := msg.TaskCtl.Workflow[nextIndex]

	log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "worker").Debugf("send next workflow [id=%s, next=%s]", msg.Id, next.WorkerName)

	msg.FuncArgs = result.FuncReturn
	msg.TaskCtl.SetRetryCount(next.RetryCount)
	if next.RunAfter != 0 {
		n := time.Now()
		msg.TaskCtl.SetRunTime(n.Add(next.RunAfter))
	}
	if !next.ExpireTime.IsZero() {
		msg.TaskCtl.SetExpireTime(next.ExpireTime)
	}
	groupName := next.GroupName
	if !msg.TaskCtl.IsZeroRunTime() {
		groupName = t.GetDelayGroupName(groupName)
	}
	msg.WorkerName = next.WorkerName
	err := t.SendMsg(groupName, msg)

	if err != nil {
		log.YTaskLog.WithField("server", t.groupName).WithField("goroutine", "worker").Errorf("send next workflow error %s [id=%s]", err, msg.Id)
		result.Err = yerrors.ErrSendMsg{Msg: err.Error()}
		t.workerGoroutine_UpdateResultStatus(message.ResultStatus.Failure, nextIndex, &result)
	} else {
		t.workerGoroutine_UpdateResultStatus(message.ResultStatus.Sent, nextIndex, &result)
	}

	t.workerGoroutine_SaveResult(result)

}
