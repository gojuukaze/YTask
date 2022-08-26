package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/gojuukaze/YTask/v3/core/brokers"
	"github.com/gojuukaze/YTask/v3/core/message"
	"github.com/gojuukaze/YTask/v3/core/util/yjson"
	"github.com/gojuukaze/YTask/v3/core/yerrors"
	"time"
)

type Broker struct {
	client   *Client
	host     string
	port     string
	password string
	db       int
	poolSize int
}

// NewRedisBroker
//  - poolSize: Maximum number of idle connections in client pool.
//              If clientPoolSize<=0, clientPoolSize=10
func NewRedisBroker(host string, port string, password string, db int, poolSize int) Broker {
	return Broker{
		host:     host,
		port:     port,
		password: password,
		db:       db,
		poolSize: poolSize,
	}
}

func (r *Broker) Activate() {
	client := NewRedisClient(r.host, r.port, r.password, r.db, r.poolSize)
	r.client = &client
}

func (r *Broker) SetPoolSize(n int) {
	r.poolSize = n
}
func (r *Broker) GetPoolSize() int {
	return r.poolSize
}

func (r *Broker) Next(queueName string) (message.Message, error) {
	var msg message.Message
	values, err := r.client.BLPop(queueName, 2*time.Second).Result()
	if err != nil {
		if err == redis.Nil {
			return msg, yerrors.ErrEmptyQueue{}
		}
		return msg, err
	}

	err = yjson.YJson.UnmarshalFromString(values[1], &msg)
	return msg, err
}

func (r *Broker) Send(queueName string, msg message.Message) error {
	b, err := yjson.YJson.Marshal(msg)

	if err != nil {
		return err
	}
	err = r.client.RPush(queueName, b)
	return err
}

func (r *Broker) LSend(queueName string, msg message.Message) error {
	b, err := yjson.YJson.Marshal(msg)

	if err != nil {
		return err
	}
	err = r.client.LPush(queueName, b)
	return err
}

func (r Broker) Clone() brokers.BrokerInterface {

	return &Broker{
		host:     r.host,
		port:     r.port,
		password: r.password,
		db:       r.db,
		poolSize: r.poolSize,
	}
}
