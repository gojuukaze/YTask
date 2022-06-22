package util

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrConnPoolGetTimeout = errors.New("ErrConnPoolGetTimeout")
	ErrNoIdleConn         = errors.New("NoIdleConn")
)

type ConnInterface interface {
	Dail(timeout int) error
	Close()
	GetId() int
	// SetId 注意，这个id要存下来。(这个id是为了调试用)
	SetId(id int)
	// Clone 复制配置（host，pass等），然后返回指针（注意，必须返回【指针】！！！）
	Clone() interface{}
}

type ConnPoolConfig struct {
	Size             int
	IdleTimeout      int // 链接超时时间. Second
	GetConnTimeout   int
	GetConnSleepTime time.Duration
	ReDailTimes      int // Dail报错的重试次数
}

type ConnPool struct {
	Conf      ConnPoolConfig
	Conn      ConnInterface
	Mu        sync.Mutex
	IdleConn  map[int]interface{} // 用slice会频繁修改，chan比较适合跨协程场景，因此用map
	Locked    bool
	NumOpen   int
	ConnCount int
}

func NewConnPoolConfig() ConnPoolConfig {
	return ConnPoolConfig{
		Size:             10,
		IdleTimeout:      60,
		GetConnTimeout:   20,
		GetConnSleepTime: 100 * time.Millisecond,
		ReDailTimes:      2,
	}
}
func NewConnPool(conn ConnInterface, conf ConnPoolConfig) ConnPool {
	return ConnPool{
		Conf:     conf,
		Conn:     conn,
		IdleConn: make(map[int]interface{}),
		NumOpen:  0,
		Locked:   false,
	}
}
func (p *ConnPool) Lock() {
	p.Mu.Lock()
	p.Locked = true
}
func (p *ConnPool) UnLock() {
	p.Locked = false
	p.Mu.Unlock()
}

func (p *ConnPool) Get() (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(p.Conf.GetConnTimeout)*time.Second)
	for {
		conn, err := p.get()
		if err == nil {
			return conn, err
		} else if err != ErrNoIdleConn {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, ErrConnPoolGetTimeout
		case <-time.After(p.Conf.GetConnSleepTime):

		}
	}
}

func (p *ConnPool) get() (interface{}, error) {
	p.Lock()
	defer p.UnLock()

	if p.NumOpen == p.Conf.Size {
		return nil, ErrNoIdleConn
	}
	for id, conn := range p.IdleConn {
		p.IdleConn[id] = nil
		return conn, nil
	}
	newConn := p.Conn.Clone()
	var err error
	for num := p.Conf.ReDailTimes; num > 0; num-- {
		err = newConn.(ConnInterface).Dail(p.Conf.IdleTimeout - 1)
		time.Sleep(p.Conf.GetConnSleepTime)
	}
	if err != nil {
		return nil, err
	}
	p.ConnCount++
	newConn.(ConnInterface).SetId(p.ConnCount)
	return &newConn, err

}

func (p *ConnPool) NewConn() {
	if !p.Locked {
		panic("Need to lock before use NewConn")
	}

}
