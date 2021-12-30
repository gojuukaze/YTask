package brokers

import (
	"github.com/go-redis/redis/v7"
	"github.com/gojuukaze/YTask/v2/drive"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/util/yjson"
	"github.com/gojuukaze/YTask/v2/yerrors"
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

func (r *RedisBroker) Next(queueName string) (message.Message, error) {
	var msg message.Message
	values, err := r.client.BLPop(queueName, 2*time.Second).Result()
	if err != nil {
		if err == redis.Nil {
			return msg, yerrors.ErrEmptyQuery{}
		}
		return msg, err
	}

	err = yjson.YJson.UnmarshalFromString(values[1], &msg)
	return msg, err
}

func (r *RedisBroker) Send(queueName string, msg message.Message) error {
	b, err := yjson.YJson.Marshal(msg)

	if err != nil {
		return err
	}
	err = r.client.RPush(queueName, b)
	return err
}

func (r *RedisBroker) LSend(queueName string, msg message.Message) error {
	b, err := yjson.YJson.Marshal(msg)

	if err != nil {
		return err
	}
	err = r.client.LPush(queueName, b)
	return err
}

func (r RedisBroker) Clone() BrokerInterface {

	return &RedisBroker{
		host:     r.host,
		port:     r.port,
		password: r.password,
		db:       r.db,
		poolSize: r.poolSize,
	}
}
