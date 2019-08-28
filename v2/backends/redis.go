package backends

import (
	"encoding/json"
	"github.com/gojuukaze/YTask/v2/drive"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/yerrors"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisBackend struct {
	client   *drive.RedisClient
	host     string
	port     string
	password string
	db       int
	poolSize int
}

func NewRedisBackend(host string, port string, password string, db int, poolSize int) RedisBackend {
	return RedisBackend{
		host:     host,
		port:     port,
		password: password,
		db:       db,
		poolSize: poolSize,
	}
}

func (r *RedisBackend) Activate() {
	client := drive.NewRedisClient(r.host, r.port, r.password, r.db, r.poolSize)
	r.client = &client
}

func (r *RedisBackend) SetPoolSize(n int) {
	r.poolSize = n
}
func (r *RedisBackend) GetPoolSize() int {
	return r.poolSize
}
func (r *RedisBackend) SetResult(result message.Result, exTime int) error {
	s, err := json.Marshal(result)

	if err != nil {
		return err
	}
	_, err = r.client.Set(result.GetBackendKey(), s, time.Duration(exTime)*time.Second)
	return err
}
func (r *RedisBackend) GetResult(key string) (message.Result, error) {
	var result message.Result
	b, err := redis.Bytes(r.client.Get(key))
	if err != nil {
		if err == redis.ErrNil {
			return result, yerrors.ErrNilResult{}
		}
		return result, err
	}

	err = json.Unmarshal(b, &result)
	return result, err
}
