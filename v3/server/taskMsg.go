package server

import (
	"errors"
	"github.com/gojuukaze/YTask/v3/yerrors"
	"time"
)

/*
Message 与 TaskMessage 是差不多的（TaskCtl中的参数有细微差别），之所以又写了个TaskMessage是为了解决循环引用问题

因为任务函数中可以使用TaskCtl控制一些东西，但TaskCtl会用到broker,backend等，broker,backend中又用到了message这样就造成了循环引用
（这里必须fuck一下go的循环引用！！！）


 - Message用于client端send时，只保存任务参数；
 - server端获取Message后会转为TaskMessage

 Message中为了便于区分，task_ctl改名为MsgArgs
*/

type TaskMessage struct {
	Id         string   `json:"id"`
	WorkerName string   `json:"worker_name"`
	FuncArgs   []string `json:"func_args"` //yjson string slice
	Ctl        TaskCtl  `json:"task_ctl"`
}
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
	su         *ServerUtils
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
