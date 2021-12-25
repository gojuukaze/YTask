package drive

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemCacheClient struct {
	memCacheConn *memcache.Client
}

func NewMemCacheClient(host, port string, poolSize int) MemCacheClient {
	client := memcache.New(fmt.Sprintf("%s:%s", host, port))
	client.MaxIdleConns = poolSize
	err := client.Ping()
	if err != nil {
		panic("YTask: connect memCached error : " + err.Error())
	}

	return MemCacheClient{
		memCacheConn: client,
	}
}

// =======================
// high api
// =======================
func (c *MemCacheClient) Get(key string) (string, error) {
	msg, err := c.memCacheConn.Get(key)
	if err != nil {
		return "", err
	}
	return string(msg.Value), nil
}

func (c *MemCacheClient) Set(key string, value interface{}, exTime int) error {
	err := c.memCacheConn.Set(&memcache.Item{Key: key, Value: value.([]byte), Expiration: int32(exTime)})
	if err != nil {
		return err
	}
	return nil
}

func (c *MemCacheClient) Flush() error {
	return c.memCacheConn.FlushAll()
}

func (c *MemCacheClient) Ping() error {
	return c.memCacheConn.Ping()
}

func (c *MemCacheClient) Close() {

}
