package message

import (
	"github.com/gojuukaze/YTask/v2/controller"
	"github.com/gojuukaze/YTask/v2/util"
	"github.com/google/uuid"
)

type Message struct {
	Id         string             `json:"uuid"`
	WorkerName string             `json:"worker_name"`
	JsonArgs   string             `json:"json_args"` // [ {"type":"int", "value":123 } , ...]
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
	s, err := util.GoArgsToJson(args...)
	if err != nil {
		return err
	}
	m.JsonArgs = s
	return nil
}
