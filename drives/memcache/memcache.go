package memcache

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type Client struct {
	Server       []string
	PoolSize     int
	memCacheConn *memcache.Client
}

// NewMemCacheClient
//  - server: ["host:port"]
func NewMemCacheClient(server []string, poolSize int) *Client {
	client := Client{server, poolSize, nil}
	client.Dail()
	err := client.Ping()
	if err != nil {
		panic("YTask: connect memCached error : " + err.Error())
	}

	return &client
}

func (c *Client) Dail() {
	c.memCacheConn = memcache.New(c.Server...)
	c.memCacheConn.MaxIdleConns = c.PoolSize
	// 为什么默认的Timeout是100millisecond？？这么短就没必要用连接池了。所以这里改长一点。
	// 写这个注释是因为不熟悉memcached，不知道设这么短是否有特殊的原因，后人看到有问题的话请反馈
	c.memCacheConn.Timeout = time.Second * 10
}

func (c *Client) Check() error {
	err := c.Ping()
	if err != nil && err.Error() == "EOF" {
		c.Dail()
		return nil
	}
	return err
}

func (c *Client) Get(key string) ([]byte, error) {
	err := c.Check()
	if err != nil {
		return nil, err
	}
	item, err := c.memCacheConn.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

func (c *Client) Set(key string, value []byte, exTime int) error {
	err := c.Check()
	if err != nil {
		return err
	}
	err = c.memCacheConn.Set(&memcache.Item{Key: key, Value: value, Expiration: int32(exTime)})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Flush() error {
	return c.memCacheConn.FlushAll()
}

func (c *Client) Ping() error {
	return c.memCacheConn.Ping()
}

func (c *Client) Close() {

}
