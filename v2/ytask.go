package ytask

import (
	"github.com/gojuukaze/YTask/v2/backends"
	"github.com/gojuukaze/YTask/v2/brokers"
	"github.com/gojuukaze/YTask/v2/config"
	"github.com/gojuukaze/YTask/v2/log"
	"github.com/gojuukaze/YTask/v2/server"
)

var (
	Server  = iServer{}
	Broker  = iBroker{}
	Logger  = iLogger{}
	Config  = iConfig{}
	Backend = iBackend{}
)

type iServer struct {
}

func (i iServer) NewServer(setConfigFunc ...config.SetConfigFunc) server.Server {
	c := config.NewConfig(setConfigFunc...)
	return server.NewServer(c)
}

type iBroker struct {
}

//
// clientPoolSize: Maximum number of idle connections in client pool.
//                 If clientPoolSize<=0, clientPoolSize=10
//
func (i iBroker) NewRedisBroker(host string, port string, password string, db int, clientPoolSize int) brokers.RedisBroker {
	return brokers.NewRedisBroker(host, port, password, db, clientPoolSize)
}

func (i iBroker) NewRabbitMqBroker(host, port, user, password, vhost string) brokers.RabbitMqBroker {
	return brokers.NewRabbitMqBroker(host, port, user, password, vhost)
}

func (i iBroker) NewRocketMqBroker(namesrvAddr []string,brokerAddr... []string) brokers.RocketMqBroker {
	return brokers.NewRocketMqBroker(namesrvAddr,brokerAddr...)
}

type iConfig struct {
}

func (i iConfig) Broker(b brokers.BrokerInterface) config.SetConfigFunc {
	return config.Broker(b)
}

func (i iConfig) Backend(b backends.BackendInterface) config.SetConfigFunc {
	return config.Backend(b)
}

func (i iConfig) Logger(l log.LoggerInterface) config.SetConfigFunc {
	return config.Logger(l)
}

func (i iConfig) Debug(debug bool) config.SetConfigFunc {
	return config.Debug(debug)
}

func (i iConfig) EnableDelayServer(enable bool) config.SetConfigFunc {
	return config.EnableDelayServer(enable)
}

func (i iConfig) DelayServerQueueSize(size int) config.SetConfigFunc {
	return config.DelayServerQueueSize(size)
}

// default: 1 day
// task status expires in ex seconds, -1:forever,
func (i iConfig) StatusExpires(ex int) config.SetConfigFunc {
	return config.StatusExpires(ex)
}

// default: 1day
// task result expires in ex seconds, -1:forever,
func (i iConfig) ResultExpires(ex int) config.SetConfigFunc {
	return config.ResultExpires(ex)
}

type iBackend struct {
}

//
// poolSize: Maximum number of idle connections in the pool. If poolSize<=0 use default value
//           default value is min(10, numWorkers) at server
//           default value is 10 at client
//
func (i iBackend) NewRedisBackend(host string, port string, password string, db int, poolSize int) backends.RedisBackend {
	return backends.NewRedisBackend(host, port, password, db, poolSize)
}

func (i iBackend) NewMemCacheBackend(host, port string, poolSize int) backends.MemCacheBackend {
	return backends.NewMemCacheBackend(host, port, poolSize)
}

func (i iBackend) NewMongoBackend(host, port, user, password, db, collection string) backends.MongoBackend {
	return backends.NewMongoBackend(host, port, user, password, db, collection)
}

type iLogger struct {
}

func (i iLogger) NewYTaskLogger() log.LoggerInterface {
	return log.NewYTaskLogger(log.YTaskLog)
}
