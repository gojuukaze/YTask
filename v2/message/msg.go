package message

import (
	"github.com/gojuukaze/YTask/v2/util"
	"github.com/google/uuid"
)

type Message struct {
	Id         string `json:"uuid"`
	WorkerName string `json:"worker_name"`
	JsonArgs   string `json:"json_args"` // [ {"type":"int", "value":123 } , ...]
}

func NewMessage() Message {
	id := uuid.New().String()
	return Message{
		Id: id,
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
