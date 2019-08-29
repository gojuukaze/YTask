package drive

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisClient struct {
	redisPool redis.Pool
	host string
	port string
	password string
	db int
}

func NewRedisClient(host string, port string, password string, db int, numConns int) RedisClient {

	client := RedisClient{
		redisPool: redis.Pool{
			MaxIdle:     numConns,
			MaxActive:   numConns,
			IdleTimeout: 1 * time.Hour,
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				con, err := redis.Dial("tcp",
					host+":"+port,
					redis.DialPassword(password),
					redis.DialDatabase(db),
					redis.DialConnectTimeout(5*time.Second),
					redis.DialReadTimeout(5*time.Second),
					redis.DialWriteTimeout(5*time.Second),
				)
				if err != nil {
					return nil, err
				}
				return con, nil
			},
		},
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
	return redis.Bool(c.Do("EXISTS", key))
}
func (c *RedisClient) Get(key string) (interface{}, error) {
	key = HideKey(key)
	return c.Do("GET", key)
}

func (c *RedisClient) Set(key string, value interface{}, exTime time.Duration) (interface{}, error) {
	key = HideKey(key)
	if exTime <= 0 {
		return c.Do("SET", key, value)
	} else {
		return c.Do("SET", key, value, "EX", int(exTime.Seconds()))
	}
}

func (c *RedisClient) Expire(key string, t time.Duration) error {
	key = HideKey(key)
	_, err := c.Do("EXPIRE", key, int(t.Seconds()))
	return err
}
func (c *RedisClient) Del(key string) error {
	key = HideKey(key)
	_, err := c.Do("DEL", key)
	return err
}

func (c *RedisClient) LPush(key string, value interface{}) (interface{}, error) {
	key = HideKey(key)
	return c.Do("lpush", key, value)
}

func (c *RedisClient) RPush(key string, value interface{}) (interface{}, error) {
	key = HideKey(key)
	return c.Do("rpush", key, value)
}

func (c *RedisClient) LPop(key string) (interface{}, error) {
	key = HideKey(key)
	return c.Do("lpop", key)
}

func (c *RedisClient) BLPop(key string, timeout time.Duration) (interface{}, error) {
	key = HideKey(key)

	return c.Do("blpop", key, int(timeout.Seconds()))
}

func (c *RedisClient) BRPop(key string, timeout time.Duration) (interface{}, error) {
	key = HideKey(key)
	return c.Do("brpop", key, int(timeout.Seconds()))
}

// =======================
// lower api
// =======================
func (c *RedisClient) Do(commandName string, args ...interface{}) (interface{}, error) {
	client := c.redisPool.Get()
	//fmt.Printf("%v , %p\n", c.redisPool, &c.redisPool)
	defer client.Close()
	return client.Do(commandName, args...)
}

func (c *RedisClient) Send(commandName string, args ...interface{}) error {
	client := c.redisPool.Get()
	defer client.Close()
	return client.Send(commandName, args...)
}

func (c *RedisClient) Flush() error {
	client := c.redisPool.Get()
	defer client.Close()
	return client.Flush()
}

func (c *RedisClient) Receive() (interface{}, error) {
	client := c.redisPool.Get()
	defer client.Close()
	return client.Receive()
}

func (c *RedisClient) Ping() error {
	client := c.redisPool.Get()
	defer client.Close()
	_, err := client.Do("ping")
	return err
}

func (c *RedisClient) Close() {
	c.redisPool.Close()
}
