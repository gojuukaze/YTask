package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type Client interface {
	Exists(key string) (bool, error)
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, exTime time.Duration) error
	RPush(key string, value interface{}) error
	LPush(key string, value interface{}) error
	BLPop(key string, timeout time.Duration) *redis.StringSliceCmd
	Do(args ...interface{}) *redis.Cmd
	Flush() error
	Ping() error
	Close()
}
