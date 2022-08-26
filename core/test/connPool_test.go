package test

import (
	"github.com/gojuukaze/YTask/v3/core/util"
	"testing"
	"time"
)

type testConn struct {
	id       string
	user     string
	password string
	conn     *struct{}
}

func (t *testConn) Dail(timeout int) error {
	t.conn = &struct{}{}
	return nil
}
func (t *testConn) Close() {
	// close
}
func (t *testConn) GetId() string {
	return t.id
}

func (t *testConn) SetId(id string) {
	t.id = id
}

// Clone 复制配置（host，pass等），然后返回指针（注意，必须返回【指针】！！！）
func (t *testConn) Clone() interface{} {
	return &testConn{
		user:     t.user,
		password: t.password,
	}
}

func (t *testConn) IsActive() bool {
	// ping ...
	return true
}

func TestConnPool(t *testing.T) {
	conf := util.NewConnPoolConfig()
	conf.IdleTimeout = 3
	conf.GetConnTimeout = 1
	conf.Size = 2
	pool := util.NewConnPool(&testConn{}, conf)

	// test1
	c1, err := pool.Get()
	if err != nil {
		t.Fatal("pool get err", err)
	}
	c2, err := pool.Get()
	if err != nil {
		t.Fatal("pool get err", err)
	}
	_, err = pool.Get()
	// 连接池为2，所以这里应该timeout
	if err != util.ErrConnPoolGetTimeout {
		t.Fatal("pool does not timeout ", err)
	}

	pool.Put(c1, false)
	pool.Put(c2, false)
	// 此时应该能顺利放回
	if len(pool.IdleConn) != 2 {
		t.Fatal("len(pool.IdleConn)!=2")
	}
	time.Sleep(3 * time.Second)
	// test2
	c3, err := pool.Get()
	// 原来的c1,c2应该过期了，c3应该是新的
	if c1 == c3 || c2 == c3 {
		t.Fatal("c1==c3 || c2==c3")
	}
	pool.Put(c3, true)

	// 因为c1,c2超时， c3设为bad，此时应该没链接
	if len(pool.IdleConn) > 0 {
		t.Fatal("len(pool.IdleConn)>0")
	}

	// test3
	c1, err = pool.Get()
	pool.Put(c1, false)
	c2, err = pool.Get()
	pool.Put(c1, false)
	// 拿出又放回，此时应该是同一个链接
	if c1 != c2 {
		t.Fatal("c1!=c2")
	}

}
