package message

import (
	"github.com/gojuukaze/YTask/v2/util"
)

type Message struct {
	WorkerName string `json:"worker_name"`
	JsonArgs   string `json:"json_args"` // [ {"type":"int", "value":123 } , ...]
}

func (m *Message) SetArgs(args ...interface{}) error {
	s, err := util.GoArgsToJson(args...)
	if err != nil {
		return err
	}
	m.JsonArgs = s
	return nil
}
