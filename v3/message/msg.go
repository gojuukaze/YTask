package message

import (
	"github.com/gojuukaze/YTask/v3/util"
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Id         string   `json:"id"`
	WorkerName string   `json:"worker_name"`
	FuncArgs   []string `json:"func_args"` //yjson string slice

	MsgArgs MessageArgs `v2JsonName:"task_ctl"` // 这里面的参数client端send时可通过setTaskArgs修改，因此单独放在一个结构体里

}

type MessageArgs struct {
	RetryCount int
	RunTime    time.Time
	ExpireTime time.Time
	Workflow   []MessageWorkflowArgs `json:"workflow"`
}

type MessageWorkflowArgs struct {
	GroupName  string
	WorkerName string
	RetryCount int
	RunAfter   time.Duration
	ExpireTime time.Time
}

func NewMsgArgs() MessageArgs {
	return MessageArgs{RetryCount: 3}
}
func NewMessage(msgArgs MessageArgs) Message {
	id := uuid.New().String()
	return Message{
		Id:      id,
		MsgArgs: msgArgs,
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

func (m Message) IsRunTime() bool {
	n := time.Now().Unix()
	return n >= m.MsgArgs.GetRunTime().Unix()
}

func (m Message) RunTimeAfter(t time.Time) bool {
	return m.MsgArgs.GetRunTime().Unix() > t.Unix()
}

func (m Message) RunTimeAfterOrEqual(t time.Time) bool {
	return m.MsgArgs.GetRunTime().Unix() >= t.Unix()
}

func (m Message) RunTimeBefore(t time.Time) bool {
	return m.MsgArgs.GetRunTime().Unix() < t.Unix()
}

func (m Message) RunTimeBeforeOrEqual(t time.Time) bool {
	return m.MsgArgs.GetRunTime().Unix() <= t.Unix()
}

func (m Message) RunTimeEqual(t time.Time) bool {
	return m.MsgArgs.GetRunTime().Unix() == t.Unix()
}

func (m MessageArgs) IsDelayMessage() bool {
	return !m.RunTime.IsZero()
}

func (t *MessageArgs) AppendWorkflow(work MessageWorkflowArgs) {
	t.Workflow = append(t.Workflow, work)
}

func (t MessageArgs) GetRunTime() time.Time {
	return t.RunTime
}
