package drive

import (
	"context"
	//"github.com/gomodule/redigo/redis"
	"github.com/go-redis/redis/v7"

	"time"
)

type RedisClient struct {
	redisPool *redis.Client
}

func NewRedisClient(host string, port string, password string, db int, poolSize int) RedisClient {
	client := RedisClient{
		redisPool: redis.NewClient(&redis.Options{
			Addr:        host + ":" + port,
			Password:    password,
			DB:          db,
			PoolSize:    poolSize,
			PoolTimeout: 10 * time.Second,
		}),
	}
	err := client.Ping()
	if err != nil {
		panic("YTask: connect redisBroker error : " + err.Error())
	}
	return client

}

func HideKey(key string) string {
	// todo
	//i := strings.Index(key, "::")
	//if i >= 0 {
	//	key = fmt.Sprintf("%s:%s", key[:i], util.GetStrMd5(key[i+2:]))
	//}
	return key
}

// =======================
// high api
// =======================
func (c *RedisClient) Exists(key string) (bool, error) {
	key = HideKey(key)
	r, err := c.redisPool.Exists(key).Result()
	return r == 1, err
}
func (c *RedisClient) Get(key string) *redis.StringCmd {
	key = HideKey(key)

	return c.redisPool.Get(key)
}

func (c *RedisClient) Set(key string, value interface{}, exTime time.Duration) error {
	key = HideKey(key)

	if exTime <= 0 {
		return c.redisPool.Set(key, value, 0).Err()
	} else {
		return c.redisPool.Set(key, value, exTime).Err()
	}
}

func (c *RedisClient) RPush(key string, value interface{}) error {
	key = HideKey(key)
	return c.redisPool.RPush(key, value).Err()
}

func (c *RedisClient) LPush(key string, value interface{}) error {
	key = HideKey(key)
	return c.redisPool.LPush(key, value).Err()
}

func (c *RedisClient) BLPop(key string, timeout time.Duration) *redis.StringSliceCmd {
	key = HideKey(key)

	return c.redisPool.BLPop(timeout, key)
}

func (c *RedisClient) Do(args ...interface{}) *redis.Cmd {
	var ctx = context.Background()

	return c.redisPool.Do(ctx,args)
}

func (c *RedisClient) Flush() error {

	return c.redisPool.FlushDB().Err()
}

func (c *RedisClient) Ping() error {

	return c.redisPool.Ping().Err()
}

func (c *RedisClient) Close() {
	c.redisPool.Close()
}
