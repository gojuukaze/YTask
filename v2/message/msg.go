package message

import (
	"github.com/gojuukaze/YTask/v2/controller"
	"github.com/gojuukaze/YTask/v2/util"
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Id         string             `json:"id"`
	WorkerName string             `json:"worker_name"`
	FuncArgs   []string           `json:"func_args"` //yjson string slice
	TaskCtl    controller.TaskCtl `json:"task_ctl"`
}

func NewMessage(ctl controller.TaskCtl) Message {
	id := uuid.New().String()
	return Message{
		Id:      id,
		TaskCtl: ctl,
	}
}

func (m *Message) SetArgs(args ...interface{}) error {
	r, err := util.GoVarsToYJsonSlice(args...)
	if err != nil {
		return err
	}
	m.FuncArgs = r
	return nil
}

func (m Message) IsDelayMessage() bool {
	return m.TaskCtl.IsZeroRunTime()
}

func (m Message) IsRunTime() bool {
	n := time.Now().Unix()
	return n >= m.TaskCtl.GetRunTime().Unix()
}

func (m Message) RunTimeAfter(t time.Time) bool {
	return m.TaskCtl.GetRunTime().Unix() > t.Unix()
}

func (m Message) RunTimeAfterOrEqual(t time.Time) bool {
	return m.TaskCtl.GetRunTime().Unix() >= t.Unix()
}

func (m Message) RunTimeBefore(t time.Time) bool {
	return m.TaskCtl.GetRunTime().Unix() < t.Unix()
}

func (m Message) RunTimeBeforeOrEqual(t time.Time) bool {
	return m.TaskCtl.GetRunTime().Unix() <= t.Unix()
}

func (m Message) RunTimeEqual(t time.Time) bool {
	return m.TaskCtl.GetRunTime().Unix() == t.Unix()
}
