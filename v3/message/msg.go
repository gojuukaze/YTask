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

	//MsgArgs    taskMessage.TaskCtl `json:"task_ctl"`
	MsgArgs MessageArgs `json:"task_ctl"` // 这个MsgArgs就是之前的task_ctl，之所以改了个类型为了解决循环引用问题。
	// 具体说明见 server/taskMsg.go 中的注释
	//

}

type MessageWorkflowArgs struct {
	GroupName  string
	WorkerName string
	RetryCount int
	RunAfter   time.Duration
	ExpireTime time.Time
}

type MessageArgs struct {
	RetryCount int
	RunTime    time.Time
	ExpireTime time.Time
	Workflow   []MessageWorkflowArgs `json:"workflow"`
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

func (m Message) IsDelayMessage() bool {
	return m.MsgArgs.RunTime.IsZero()
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
