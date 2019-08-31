package message

import (
	"github.com/gojuukaze/YTask/v2/controller"
	"github.com/gojuukaze/YTask/v2/util"
	"github.com/google/uuid"
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
