package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type StandaloneClient struct {
	redisPool *redis.Client
}

func NewRedisClient(host string, password string, db int, poolSize int) StandaloneClient {
	client := StandaloneClient{
		redisPool: redis.NewClient(&redis.Options{
			Addr:        host,
			Password:    password,
			DB:          db,
			PoolSize:    poolSize,
			PoolTimeout: 60 * time.Second,
		}),
	}
	err := client.Ping()
	if err != nil {
		panic("YTask: connect redisBroker error : " + err.Error())
	}
	return client

}

// =======================
// high api
// =======================
func (c *StandaloneClient) Exists(key string) (bool, error) {
	r, err := c.redisPool.Exists(context.Background(), key).Result()
	return r == 1, err
}
func (c *StandaloneClient) Get(key string) *redis.StringCmd {

	return c.redisPool.Get(context.Background(), key)
}

func (c *StandaloneClient) Set(key string, value interface{}, exTime time.Duration) error {

	if exTime <= 0 {
		exTime = 0
	}
	return c.redisPool.Set(context.Background(), key, value, exTime).Err()

}

func (c *StandaloneClient) RPush(key string, value interface{}) error {
	return c.redisPool.RPush(context.Background(), key, value).Err()
}

func (c *StandaloneClient) LPush(key string, value interface{}) error {
	return c.redisPool.LPush(context.Background(), key, value).Err()
}

func (c *StandaloneClient) BLPop(key string, timeout time.Duration) *redis.StringSliceCmd {

	return c.redisPool.BLPop(context.Background(), timeout, key)
}

func (c *StandaloneClient) Do(args ...interface{}) *redis.Cmd {
	var ctx = context.Background()

	return c.redisPool.Do(ctx, args)
}

func (c *StandaloneClient) Flush() error {

	return c.redisPool.FlushDB(context.Background()).Err()
}

func (c *StandaloneClient) Ping() error {

	return c.redisPool.Ping(context.Background()).Err()
}

func (c *StandaloneClient) Close() {
	c.redisPool.Close()
}
