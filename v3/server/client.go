package server

import (
	"github.com/gojuukaze/YTask/v3/config"
	"github.com/gojuukaze/YTask/v3/message"
	"github.com/gojuukaze/YTask/v3/yerrors"
	"time"
)

type ctlKeyChoices struct {
	RetryCount int
	RunAt      int
	RunAfter   int
	ExpireTime int
}

var ctlKey = ctlKeyChoices{
	RetryCount: 0,
	RunAt:      1,
	RunAfter:   2,
	ExpireTime: 3,
}

type Client struct {
	sUtils  *ServerUtils
	isClone bool
	msgArgs message.MessageArgs

	// msgArgs name
	ctlKeyChoices
}

func NewClient(c config.Config) Client {
	su := newServerUtils(c.Broker, c.Backend, c.Logger, c.StatusExpires, c.ResultExpires)
	client := Client{
		sUtils:        &su,
		msgArgs:       message.NewMsgArgs(),
		ctlKeyChoices: ctlKey,
	}

	if client.sUtils.GetBrokerPoolSize() <= 0 {
		client.sUtils.SetBrokerPoolSize(10)
	}
	client.sUtils.BrokerActivate()

	if client.sUtils.backend != nil {
		if client.sUtils.GetBackendPoolSize() <= 0 {
			client.sUtils.SetBackendPoolSize(10)
		}
		client.sUtils.BackendActivate()
	}
	return client
}

func (c *Client) Clone() *Client {
	if c.isClone {
		return c
	} else {
		return &Client{
			sUtils:        c.sUtils,
			isClone:       true,
			msgArgs:       c.msgArgs,
			ctlKeyChoices: ctlKey,
		}
	}
}

func (c *Client) SetTaskCtl(name int, value interface{}) *Client {
	cloneC := c.Clone()
	switch name {
	case ctlKey.RetryCount:
		cloneC.msgArgs.RetryCount = value.(int)
	case ctlKey.RunAfter:
		n := time.Now()
		cloneC.msgArgs.RunTime = n.Add(value.(time.Duration))
	case ctlKey.RunAt:
		cloneC.msgArgs.RunTime = value.(time.Time)
	case ctlKey.ExpireTime:
		cloneC.msgArgs.ExpireTime = value.(time.Time)
	}
	return cloneC
}

// Send
// return: taskId, err
func (c *Client) Send(groupName string, workerName string, args ...interface{}) (string, error) {
	if c.msgArgs.IsDelayMessage() {
		groupName = c.sUtils.GetDelayGroupName(groupName)
	}
	return c.sUtils.Send(groupName, workerName, c.msgArgs, args...)
}

// Workflow
// start a workflow
func (c *Client) Workflow() *ClientWithWorkflow {
	cloneC := c.Clone()
	return &ClientWithWorkflow{client: cloneC}
}

// GetResult
// Only return the result when the task is over
// 只有任务结束才返回结果（任务失败也算结束）
//  - taskId:
//  - timeout:
//  - sleepDuration:
func (c *Client) GetResult(taskId string, timeout time.Duration, sleepTime time.Duration) (message.Result, error) {
	if c.sUtils.backend == nil {
		return message.Result{}, yerrors.ErrNilResult{}
	}
	n := time.Now()
	for {
		if time.Now().Sub(n) >= timeout {
			return message.Result{}, yerrors.ErrTimeOut{}
		}
		r, err := c.sUtils.GetResult(taskId)
		if err == nil && r.IsFinish() {
			return r, nil
		}
		time.Sleep(sleepTime)
	}
}

// GetResult2
//  Return the result whether the task is finished or not
//  无论任务是否结束都返回结果（此结果只要任务开始运行就有）
func (c *Client) GetResult2(taskId string, timeout time.Duration, sleepTime time.Duration) (message.Result, error) {
	if c.sUtils.backend == nil {
		return message.Result{}, yerrors.ErrNilResult{}
	}
	n := time.Now()
	for {
		if time.Now().Sub(n) >= timeout {
			return message.Result{}, yerrors.ErrTimeOut{}
		}
		r, err := c.sUtils.GetResult(taskId)
		if err == nil {
			return r, nil
		}
		time.Sleep(sleepTime)
	}
}

// taskId:
// timeout:
// sleepDuration:
func (c *Client) GetStatus(taskId string, timeout time.Duration, sleepTime time.Duration) (int, error) {
	if c.sUtils.backend == nil {
		return 0, yerrors.ErrNilResult{}
	}
	n := time.Now()
	for {
		if time.Now().Sub(n) >= timeout {
			return 0, yerrors.ErrTimeOut{}
		}
		r, err := c.sUtils.GetResult(taskId)
		if err == nil {
			return r.Status, nil
		}
		time.Sleep(sleepTime)
	}
}

// AbortTask
//  <exTime>: 过期时间，秒。<=0表示不过期
func (c *Client) AbortTask(taskID string, exTime int) error {
	return c.sUtils.AbortTask(taskID, exTime)
}

type ClientWithWorkflow struct {
	client       *Client
	WorkflowArgs message.MessageWorkflowArgs
	args         []interface{}
}

func (c *ClientWithWorkflow) SetTaskCtl(name int, value interface{}) *ClientWithWorkflow {
	switch name {
	case ctlKey.RetryCount:
		c.WorkflowArgs.RetryCount = value.(int)
	case ctlKey.RunAfter:
		c.WorkflowArgs.RunAfter = value.(time.Duration)
	case ctlKey.ExpireTime:
		c.WorkflowArgs.ExpireTime = value.(time.Time)

	}
	return c
}

// Send
//  - args : 只有第一个任务才能填！！！后续任务的参数固定为第一个任务的返回值
func (c *ClientWithWorkflow) Send(groupName string, workerName string, args ...interface{}) *ClientWithWorkflow {
	if len(args) > 0 {
		c.args = args
	}
	c.WorkflowArgs.GroupName = groupName
	c.WorkflowArgs.WorkerName = workerName

	c.client.msgArgs.AppendWorkflow(c.WorkflowArgs)
	c.WorkflowArgs = message.MessageWorkflowArgs{}
	return c

}

// SendWorkflow
// return: taskId, err
func (c *ClientWithWorkflow) Done() (string, error) {
	first := c.client.msgArgs.Workflow[0]
	c.client.SetTaskCtl(ctlKey.RetryCount, first.RetryCount)
	if first.RunAfter != 0 {
		c.client.SetTaskCtl(ctlKey.RunAfter, first.RunAfter)
	}
	if !first.ExpireTime.IsZero() {
		c.client.SetTaskCtl(ctlKey.ExpireTime, first.ExpireTime)
	}
	return c.client.Send(first.GroupName, first.WorkerName, c.args...)

}
