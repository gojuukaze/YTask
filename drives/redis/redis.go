package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Client struct {
	redisPool *redis.Client
}

func NewRedisClient(host string, port string, password string, db int, poolSize int) Client {
	client := Client{
		redisPool: redis.NewClient(&redis.Options{
			Addr:        host + ":" + port,
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
func (c *Client) Exists(key string) (bool, error) {
	r, err := c.redisPool.Exists(context.Background(), key).Result()
	return r == 1, err
}
func (c *Client) Get(key string) *redis.StringCmd {

	return c.redisPool.Get(context.Background(), key)
}

func (c *Client) Set(key string, value interface{}, exTime time.Duration) error {

	if exTime <= 0 {
		exTime = 0
	}
	return c.redisPool.Set(context.Background(), key, value, exTime).Err()

}

func (c *Client) RPush(key string, value interface{}) error {
	return c.redisPool.RPush(context.Background(), key, value).Err()
}

func (c *Client) LPush(key string, value interface{}) error {
	return c.redisPool.LPush(context.Background(), key, value).Err()
}

func (c *Client) BLPop(key string, timeout time.Duration) *redis.StringSliceCmd {

	return c.redisPool.BLPop(context.Background(), timeout, key)
}

func (c *Client) Do(args ...interface{}) *redis.Cmd {
	var ctx = context.Background()

	return c.redisPool.Do(ctx, args)
}

func (c *Client) Flush() error {

	return c.redisPool.FlushDB(context.Background()).Err()
}

func (c *Client) Ping() error {

	return c.redisPool.Ping(context.Background()).Err()
}

func (c *Client) Close() {
	c.redisPool.Close()
}
