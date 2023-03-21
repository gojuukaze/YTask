package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type ClusterClient struct {
	redisCluster *redis.ClusterClient
}

func NewRedisClusterClient(hosts []string, password string, poolSize int) ClusterClient {
	client := ClusterClient{
		redisCluster: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    hosts,
			Password: password,
			PoolSize: poolSize,
		}),
	}
	fmt.Println(hosts)
	err := client.Ping()
	if err != nil {
		panic("YTask: connect redisBroker error : " + err.Error())
	}
	return client

}

// =======================
// high api
// =======================
func (c *ClusterClient) Exists(key string) (bool, error) {
	r, err := c.redisCluster.Exists(context.Background(), key).Result()
	return r == 1, err
}
func (c *ClusterClient) Get(key string) *redis.StringCmd {

	return c.redisCluster.Get(context.Background(), key)
}

func (c *ClusterClient) Set(key string, value interface{}, exTime time.Duration) error {

	if exTime <= 0 {
		exTime = 0
	}
	return c.redisCluster.Set(context.Background(), key, value, exTime).Err()

}

func (c *ClusterClient) RPush(key string, value interface{}) error {
	return c.redisCluster.RPush(context.Background(), key, value).Err()
}

func (c *ClusterClient) LPush(key string, value interface{}) error {
	return c.redisCluster.LPush(context.Background(), key, value).Err()
}

func (c *ClusterClient) BLPop(key string, timeout time.Duration) *redis.StringSliceCmd {

	return c.redisCluster.BLPop(context.Background(), timeout, key)
}

func (c *ClusterClient) Do(args ...interface{}) *redis.Cmd {
	var ctx = context.Background()

	return c.redisCluster.Do(ctx, args)
}

func (c *ClusterClient) Flush() error {

	return c.redisCluster.FlushDB(context.Background()).Err()
}

func (c *ClusterClient) Ping() error {

	return c.redisCluster.Ping(context.Background()).Err()
}

func (c *ClusterClient) Close() {
	c.redisCluster.Close()
}
