package memcache

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gojuukaze/YTask/v3/backends"
	"github.com/gojuukaze/YTask/v3/message"
	"github.com/gojuukaze/YTask/v3/util/yjson"
	"github.com/gojuukaze/YTask/v3/yerrors"
)

type MemCacheBackend struct {
	client   *Client
	host     string
	port     string
	poolSize int
}

func NewMemCacheBackend(host, port string, poolSize int) MemCacheBackend {
	return MemCacheBackend{
		host:     host,
		port:     port,
		poolSize: poolSize,
	}
}

func (r *MemCacheBackend) Activate() {
	client := NewMemCacheClient(r.host, r.port, r.poolSize)
	r.client = &client
}

func (r *MemCacheBackend) SetPoolSize(n int) {
	r.poolSize = n
}

func (r *MemCacheBackend) GetPoolSize() int {
	return r.poolSize
}

func (r *MemCacheBackend) SetResult(result message.Result, exTime int) error {
	b, err := yjson.YJson.Marshal(result)

	if err != nil {
		return err
	}
	err = r.client.Set(result.GetBackendKey(), b, exTime)
	return err
}

func (r *MemCacheBackend) GetResult(key string) (message.Result, error) {
	var result message.Result

	b, err := r.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return result, yerrors.ErrNilResult{}
		}
		return result, err
	}

	err = yjson.YJson.Unmarshal(b, &result)
	return result, err
}

func (r MemCacheBackend) Clone() backends.BackendInterface {
	return &MemCacheBackend{
		host:     r.host,
		port:     r.port,
		poolSize: r.poolSize,
	}
}
