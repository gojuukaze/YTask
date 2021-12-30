package redisBroker

import (
	"github.com/gojuukaze/YTask/v1/yerrors"
	"github.com/gojuukaze/YTask/v1/ymsg"
	"encoding/json"
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

func (r RedisBroker) Get(queryName string) (ymsg.Message, error) {
	var msg ymsg.Message
	values, err := redis.Values(r.client.BLPop(queryName, 2*time.Second))
	if err != nil {
		if err == redis.ErrNil {
			return msg, yerrors.ErrEmptyQuery
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

func (r RedisBroker) Send(queryName string, msg ymsg.Message) error {
	s, err := json.Marshal(msg)

	if err != nil {
		return err
	}
	_, err = r.client.RPush(queryName, s)
	return err
}
