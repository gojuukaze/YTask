package drive

import (
	"context"
	//"github.com/gomodule/redigo/redis"
	"github.com/go-redis/redis/v8"

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
	var ctx = context.Background()
	r, err := c.redisPool.Exists(ctx,key).Result()
	return r == 1, err
}
func (c *RedisClient) Get(key string) *redis.StringCmd {
	key = HideKey(key)
	var ctx = context.Background()

	return c.redisPool.Get(ctx,key)
}

func (c *RedisClient) Set(key string, value interface{}, exTime time.Duration) error {
	key = HideKey(key)
	var ctx = context.Background()

	if exTime <= 0 {
		return c.redisPool.Set(ctx,key, value, 0).Err()
	} else {
		return c.redisPool.Set(ctx,key, value, exTime).Err()
	}
}

func (c *RedisClient) RPush(key string, value interface{}) error {
	key = HideKey(key)
	var ctx = context.Background()

	return c.redisPool.RPush(ctx,key, value).Err()
}

func (c *RedisClient) BLPop(key string, timeout time.Duration) *redis.StringSliceCmd {
	key = HideKey(key)
	var ctx = context.Background()

	return c.redisPool.BLPop(ctx,timeout, key)
}

func (c *RedisClient) Do(args ...interface{}) *redis.Cmd {
	var ctx = context.Background()

	return c.redisPool.Do(ctx,args)
}

func (c *RedisClient) Flush() error {
	var ctx = context.Background()

	return c.redisPool.FlushDB(ctx).Err()
}

func (c *RedisClient) Ping() error {
	var ctx = context.Background()

	return c.redisPool.Ping(ctx).Err()
}

func (c *RedisClient) Close() {
	c.redisPool.Close()
}
