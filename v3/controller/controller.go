package controller

import (
	"errors"
	"github.com/gojuukaze/YTask/v3/server"
	"github.com/gojuukaze/YTask/v3/yerrors"
	"time"
)

type TaskCtlWorkflowArgs struct {
	GroupName  string
	WorkerName string
	RetryCount int
	RunAfter   time.Duration
	ExpireTime time.Time
}

type TaskCtl struct {
	RetryCount int
	RunTime    time.Time
	ExpireTime time.Time
	err        error
	Workflow   []TaskCtlWorkflowArgs `json:"workflow"`
	id         string
	su         *server.ServerUtils
}

func NewTaskCtl() TaskCtl {
	return TaskCtl{
		RetryCount: 3,
	}
}

func (t *TaskCtl) GetTaskId() string {
	return t.id
}

func (t *TaskCtl) Retry(err error) {
	t.err = err
}

func (t *TaskCtl) SetRetryCount(c int) {
	t.RetryCount = c
}

func (t TaskCtl) CanRetry() bool {
	return t.RetryCount > 0
}

func (t TaskCtl) GetError() error {
	return t.err
}

func (t *TaskCtl) SetError(err error) {
	t.err = err
}

func (t *TaskCtl) SetRunTime(_t time.Time) {
	t.RunTime = _t
}

func (t *TaskCtl) GetRunTime() time.Time {
	return t.RunTime
}

func (t *TaskCtl) IsZeroRunTime() bool {
	return t.RunTime.IsZero()
}

func (t *TaskCtl) SetExpireTime(_t time.Time) {
	t.ExpireTime = _t
}

func (t *TaskCtl) IsExpired() bool {
	return !t.ExpireTime.IsZero() && time.Now().After(t.ExpireTime)
}

func (t *TaskCtl) AppendWorkflow(work TaskCtlWorkflowArgs) {
	t.Workflow = append(t.Workflow, work)
}

func (t *TaskCtl) Abort(msg string) {
	t.err = yerrors.ErrAbortTask{msg}
	t.RetryCount = 0
}

func (t *TaskCtl) IsAbort() error {
	if t.su == nil {
		return errors.New("IsAbort() can only be called on the server side")
	}
	return t.su.IsAbort(t.id)
}
