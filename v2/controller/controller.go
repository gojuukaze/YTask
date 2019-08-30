package controller

type TaskCtl struct {
	RetryCount int
	err        error
}

func NewTaskCtl() TaskCtl {
	return TaskCtl{
		RetryCount: 3,
	}
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
