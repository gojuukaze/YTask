package redisBroker

import (
	"encoding/json"
	"github.com/gojuukaze/YTask/v1.1/errors"
	"github.com/gojuukaze/YTask/v1.1/message"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisBroker struct {
	client *RedisClient
}

func NewRedisBroker(host string, port string, password string, db int, numConns int) RedisBroker {
	client := NewRedisClient(host, port, password, db, numConns)
	return RedisBroker{&client}
}

func (r RedisBroker) Get(queryName string) (message.Message, error) {
	var msg message.Message
	values, err := redis.Values(r.client.BLPop(queryName, 2*time.Second))
	if err != nil {
		if err == redis.ErrNil {
			return msg, errors.ErrEmptyQuery
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

func (r RedisBroker) Send(queryName string, msg message.Message) error {
	s, err := json.Marshal(msg)

	if err != nil {
		return err
	}
	_, err = r.client.RPush(queryName, s)
	return err
}
