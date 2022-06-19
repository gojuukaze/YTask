package memcache

import (
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type Client struct {
	memCacheConn *memcache.Client
}

func NewMemCacheClient(host, port string, poolSize int) Client {
	client := memcache.New(fmt.Sprintf("%s:%s", host, port))
	client.MaxIdleConns = poolSize
	// 为什么默认的Timeout是100millisecond？？这么短就没必要用连接池了。所以这里改长一点。
	// 写这个注释是因为不熟悉memcached，不知道设这么短是否有特殊的原因，后人看到有问题的话请反馈
	client.Timeout = time.Second * 10
	err := client.Ping()
	if err != nil {
		panic("YTask: connect memCached error : " + err.Error())
	}

	return Client{
		memCacheConn: client,
	}
}

func (c *Client) Get(key string) ([]byte, error) {
	item, err := c.memCacheConn.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

func (c *Client) Set(key string, value []byte, exTime int) error {
	err := c.memCacheConn.Set(&memcache.Item{Key: key, Value: value, Expiration: int32(exTime)})
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
