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
	client     Client
	clientType int // 0 standalone  1 cluster
	hosts      []string
	password   string
	db         int
	poolSize   int
}

// NewRedisBackend
//   - host: redis url. If standalone,hosts[0] is 127.0.0.1:6379.
//   - poolSize: Maximum number of idle connections in the pool. If poolSize<=0 use default value.
//     default value is min(10, numWorkers) at SERVER
//     default value is 10 at CLIENT
//   - clientType: redis server version
//     default value 0 is standalone, value 1 is clustered
func NewRedisBackend(hosts []string, password string, db int, poolSize int, clientType int) Backend {
	return Backend{
		hosts:      hosts,
		password:   password,
		db:         db,
		poolSize:   poolSize,
		clientType: clientType,
	}
}

func (r *Backend) Activate() {
	switch r.clientType {
	case 0:
		client := NewRedisClient(r.hosts[0], r.password, r.db, r.poolSize)
		r.client = &client
	case 1:
		client := NewRedisClusterClient(r.hosts, r.password, r.poolSize)
		r.client = &client
	default:
		panic("YTask: check clientType!!")
	}

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
		hosts:    r.hosts,
		password: r.password,
		db:       r.db,
		poolSize: r.poolSize,
	}
}
