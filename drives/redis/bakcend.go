package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/gojuukaze/YTask/v3/backends"
	"github.com/gojuukaze/YTask/v3/message"
	"github.com/gojuukaze/YTask/v3/util/yjson"
	"github.com/gojuukaze/YTask/v3/yerrors"
	"time"
)

type Backend struct {
	client   *Client
	host     string
	port     string
	password string
	db       int
	poolSize int
}

// NewRedisBackend
//  - poolSize: Maximum number of idle connections in the pool. If poolSize<=0 use default value.
//              default value is min(10, numWorkers) at SERVER
//              default value is 10 at CLIENT
//
func NewRedisBackend(host string, port string, password string, db int, poolSize int) Backend {
	return Backend{
		host:     host,
		port:     port,
		password: password,
		db:       db,
		poolSize: poolSize,
	}
}

func (r *Backend) Activate() {
	client := NewRedisClient(r.host, r.port, r.password, r.db, r.poolSize)
	r.client = &client
}

func (r *Backend) SetPoolSize(n int) {
	r.poolSize = n
}
func (r *Backend) GetPoolSize() int {
	return r.poolSize
}
func (r *Backend) SetResult(result message.Result, exTime int) error {
	b, err := yjson.YJson.Marshal(result)

	if err != nil {
		return err
	}
	err = r.client.Set(result.GetBackendKey(), b, time.Duration(exTime)*time.Second)
	return err
}
func (r *Backend) GetResult(key string) (message.Result, error) {
	var result message.Result

	b, err := r.client.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return result, yerrors.ErrNilResult{}
		}
		return result, err
	}

	err = yjson.YJson.Unmarshal(b, &result)
	return result, err
}

func (r Backend) Clone() backends.BackendInterface {
	return &Backend{
		host:     r.host,
		port:     r.port,
		password: r.password,
		db:       r.db,
		poolSize: r.poolSize,
	}
}
