package util

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrConnPoolGetTimeout = errors.New("ErrConnPoolGetTimeout")
	ErrConnPoolClosed     = errors.New("ErrConnPoolClosed")
	ErrNoIdleConn         = errors.New("NoIdleConn")

	IdleTimeout      = 60
	GetConnTimeout   = 60
	GetConnSleepTime = 100 * time.Millisecond
	ReDailTimes      = 2
)

type ConnInterface interface {
	Dail(timeout int) error
	Close()
	GetId() string
	// SetId 注意，这个id要存下来。这个id其实是上次使用的时间戳，用于判断是否过期。格式是: {random_int}-{time}
	SetId(id string)
	// Clone 复制配置（host，pass等），然后返回指针（注意，必须返回【指针】！！！）
	Clone() interface{}
	// IsActive 判断链接是否存活，通常是调用ping。
	// 这个操作会浪费一点时间，若觉得不必要或是没有ping方法 也可以直接返回true
	// 如果直接返回true，使用pool.Get() 返回的链接时注意判断错误，因为链接可能因各种情况关闭。
	//     注意！！就算链接已经关闭，仍然需要调用pool.Put() (此时isBad为true)，然后再调用pool.Get()重新获取链接
	IsActive() bool
}

type ConnPoolConfig struct {
	Size             int
	IdleTimeout      int // 链接超时时间. Second
	GetConnTimeout   int
	GetConnSleepTime time.Duration
	ReDailTimes      int // Dail报错的重试次数
}

type ConnPool struct {
	Conf     ConnPoolConfig
	Conn     ConnInterface
	Mu       sync.Mutex
	IdleConn map[string]interface{} // 用slice会频繁修改，chan比较适合跨协程场景，因此用map
	Locked   bool
	NumOpen  int
	IsClose  bool
}

func NewConnPoolConfig() ConnPoolConfig {
	return ConnPoolConfig{
		Size:             10,
		IdleTimeout:      IdleTimeout,
		GetConnTimeout:   GetConnTimeout,
		GetConnSleepTime: GetConnSleepTime,
		ReDailTimes:      ReDailTimes,
	}
}
func NewConnPool(conn ConnInterface, conf ConnPoolConfig) ConnPool {
	return ConnPool{
		Conf:     conf,
		Conn:     conn,
		IdleConn: make(map[string]interface{}),
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
	if p.IsClose {
		return nil, ErrConnPoolClosed
	}
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

	for id, conn := range p.IdleConn {
		isBad := p.IsConnTimeout(conn) || !conn.(ConnInterface).IsActive()
		delete(p.IdleConn, id)
		if !isBad {
			p.UpdateConnTime(conn)
			return conn, nil
		} else {
			p.CloseConn(conn)
		}
	}
	// 走到这说明要创建新的链接
	if p.NumOpen == p.Conf.Size {
		return nil, ErrNoIdleConn
	}
	return p.NewConn()

}

func (p *ConnPool) IsConnTimeout(conn interface{}) bool {
	_, t := p.DecodeId(conn.(ConnInterface).GetId())
	// +1提前一秒过期
	return time.Now().UnixMilli()+1000 > t+int64(p.Conf.IdleTimeout)*1000

}

func (p *ConnPool) NewConn() (interface{}, error) {
	if !p.Locked {
		panic("Need to lock before use NewConn")
	}
	newConn := p.Conn.Clone()
	var err error
	for num := p.Conf.ReDailTimes; num > 0; num-- {
		// Timeout + 1 防止在时间临界点时pool.Get()获取的链接无法使用
		err = newConn.(ConnInterface).Dail(p.Conf.IdleTimeout + 1)
		if err == nil {
			break
		}
		time.Sleep(p.Conf.GetConnSleepTime)
	}
	if err != nil {
		return nil, err
	}
	newConn.(ConnInterface).SetId(p.GenId(-1))
	p.NumOpen++
	return newConn, err

}

func (p *ConnPool) GenId(randomNum int) string {
	if !p.Locked {
		panic("Need to lock before use GenId")
	}
	if randomNum == -1 {
		randomNum = p.NumOpen
	}
	return fmt.Sprintf("%v-%v", randomNum, time.Now().UnixMilli())

}

func (p *ConnPool) DecodeId(id string) (n int, t int64) {
	fmt.Sscanf(id, "%d-%d", &n, &t)
	return
}

func (p *ConnPool) CloseConn(conn interface{}) {
	if !p.Locked {
		panic("Need to lock before use GenId")
	}
	delete(p.IdleConn, conn.(ConnInterface).GetId())
	p.NumOpen--
	conn.(ConnInterface).Close()
}

func (p *ConnPool) UpdateConnTime(conn interface{}) {
	n, _ := p.DecodeId(conn.(ConnInterface).GetId())
	newId := p.GenId(n)
	conn.(ConnInterface).SetId(newId)
}

// Put
// isBadConn: 用户Get获取的链接时，如果链接已关闭或者因其他原因不能在用，则isBadConn为true
func (p *ConnPool) Put(conn interface{}, isBadConn bool) {
	p.Lock()
	defer p.UnLock()
	if isBadConn || p.IsConnTimeout(conn) {
		p.CloseConn(conn)
	} else {
		p.UpdateConnTime(conn)
		p.IdleConn[conn.(ConnInterface).GetId()] = conn
	}

}

func (p *ConnPool) Close() {
	p.Lock()
	defer p.UnLock()
	for _, conn := range p.IdleConn {
		p.CloseConn(conn)
	}

}
