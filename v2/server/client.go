package server

import (
	"time"

	"github.com/vua/YTask/v2/config"
	"github.com/vua/YTask/v2/controller"
	"github.com/vua/YTask/v2/message"
	"github.com/vua/YTask/v2/yerrors"
)

type ctlKeyChoices struct {
	RetryCount int
	RunAt      int
	RunAfter   int
}

var ctlKey = ctlKeyChoices{
	RetryCount: 0,
	RunAt:      1,
	RunAfter:   2,
}

type Client struct {
	sUtils  *serverUtils
	isClone bool
	ctl     controller.TaskCtl

	// ctl name
	ctlKeyChoices
}

func NewClient(c config.Config) Client {
	su := newServerUtils(c.Broker, c.Backend, c.StatusExpires, c.ResultExpires)
	client := Client{
		sUtils:        &su,
		ctl:           controller.NewTaskCtl(),
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
			ctl:           c.ctl,
			ctlKeyChoices: ctlKey,
		}
	}
}
func (c *Client) SetTaskCtl(name int, value interface{}) *Client {
	cloneC := c.Clone()
	switch name {
	case ctlKey.RetryCount:
		cloneC.ctl.RetryCount = value.(int)
	case ctlKey.RunAfter:
		n := time.Now()
		cloneC.ctl.SetRunTime(n.Add(value.(time.Duration)))
	case ctlKey.RunAt:

		cloneC.ctl.SetRunTime(value.(time.Time))

	}
	return cloneC
}

//
// return: taskId, err
//
func (c *Client) Send(groupName string, workerName string, args ...interface{}) (string, error) {
	if !c.ctl.IsZeroRunTime() {
		groupName = c.sUtils.GetDelayGroupName(groupName)
	}
	return c.sUtils.Send(groupName, workerName, c.ctl, args...)
}

// taskId:
// timeout:
// sleepDuration:
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

type Promise struct {
	done   chan struct{}
	result interface{}
	err    error
}

type InvarParamFunc func(message.Result) (interface{}, error)
type VarParamFunc func(interface{}) (interface{}, error)

func (c *Client) NewPromise(taskId string, handle InvarParamFunc, timeout time.Duration, sleepTime time.Duration) *Promise {
	promise := Promise{done: make(chan struct{})}
	go func() {
		defer close(promise.done)
		msg, err := c.GetResult(taskId, timeout, sleepTime)
		if err != nil {
			promise.result, promise.err = nil, err
			return
		}
		promise.result, promise.err = handle(msg)
	}()
	return &promise
}

func (p *Promise) Then(handle VarParamFunc) *Promise {
	promise := Promise{done: make(chan struct{})}
	go func() {
		defer close(promise.done)
		res, err := p.Done()
		if err != nil {
			promise.result, promise.err = nil, p.err
		}
		promise.result, promise.err = handle(res)
	}()
	return &promise
}

func (p *Promise) Done() (interface{}, error) {
	<-p.done
	return p.result, p.err
}
