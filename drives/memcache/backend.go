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
	server   []string
	poolSize int
}

func NewMemCacheBackend(server []string, poolSize int) MemCacheBackend {
	return MemCacheBackend{
		server:   server,
		poolSize: poolSize,
	}
}

func (r *MemCacheBackend) Activate() {
	r.client = NewMemCacheClient(r.server, r.poolSize)
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
		server:   r.server,
		poolSize: r.poolSize,
	}
}
