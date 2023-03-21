package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/gojuukaze/YTask/v3/brokers"
	"github.com/gojuukaze/YTask/v3/message"
	"github.com/gojuukaze/YTask/v3/util/yjson"
	"github.com/gojuukaze/YTask/v3/yerrors"
	"time"
)

type Broker struct {
	client     Client
	clientType int // 0 standalone  1 cluster
	hosts      []string
	password   string
	db         int
	poolSize   int
}

// NewRedisBroker
//   - host: redis url. If standalone,hosts[0] is 127.0.0.1:6379.
//   - poolSize: Maximum number of idle connections in client pool.
//     If clientPoolSize<=0, clientPoolSize=10
//   - clientType: redis server version
//     default value 0 is standalone, value 1 is clustered
func NewRedisBroker(hosts []string, password string, db int, poolSize int, clientType int) Broker {
	return Broker{
		hosts:    hosts,
		password: password,
		db:       db,
		poolSize: poolSize,
	}
}

func (r *Broker) Activate() {
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
		hosts:    r.hosts,
		password: r.password,
		db:       r.db,
		poolSize: r.poolSize,
	}
}
