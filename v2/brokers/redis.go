package brokers

import (
	"encoding/json"
	"github.com/gojuukaze/YTask/v2/drive"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/yerrors"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisBroker struct {
	client   *drive.RedisClient
	host     string
	port     string
	password string
	db       int
	poolSize int
}

func NewRedisBroker(host string, port string, password string, db int, poolSize int) RedisBroker {
	return RedisBroker{
		host:     host,
		port:     port,
		password: password,
		db:       db,
		poolSize: poolSize,
	}
}

func (r *RedisBroker) Activate() {
	client := drive.NewRedisClient(r.host, r.port, r.password, r.db, r.poolSize)
	r.client = &client
}

func (r *RedisBroker) SetPoolSize(n int) {
	r.poolSize = n
}
func (r *RedisBroker) GetPoolSize() int {
	return r.poolSize
}

func (r *RedisBroker) Next(queryName string) (message.Message, error) {
	var msg message.Message
	values, err := redis.Values(r.client.BLPop(queryName, 2*time.Second))
	if err != nil {
		if err == redis.ErrNil {
			return msg, yerrors.ErrEmptyQuery{}
		}
		return msg, err
	}
	b, err := redis.Bytes(values[1], err)

	if err != nil {
		return msg, err
	}
	err = json.Unmarshal(b, &msg)
	return msg, err
}

func (r *RedisBroker) Send(queryName string, msg message.Message) error {
	s, err := json.Marshal(msg)

	if err != nil {
		return err
	}
	_, err = r.client.RPush(queryName, s)
	return err
}
