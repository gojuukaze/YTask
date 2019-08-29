package ytask

import (
	"github.com/gojuukaze/YTask/v2/backends"
	"github.com/gojuukaze/YTask/v2/brokers"
	"github.com/gojuukaze/YTask/v2/config"
	"github.com/gojuukaze/YTask/v2/server"
)

var (
	Server  = iServer{}
	Broker  = iBroker{}
	Config  = iConfig{}
	Backend = iBackend{}
)

type iServer struct {
}

func (is iServer) NewServer(setConfigFunc ...config.SetConfigFunc) server.Server {
	c := config.NewConfig(setConfigFunc...)
	return server.NewServer(c)
}

type iBroker struct {
}

// poolSize: ( default: 1(for server) 10(for client) ) Maximum number of idle connections in the pool. if poolSize<=0 use default
func (i iBroker) NewRedisBroker(host string, port string, password string, db int, poolSize int) brokers.RedisBroker {
	return brokers.NewRedisBroker(host, port, password, db, poolSize)
}

type iConfig struct {
}

func (i iConfig) Broker(b brokers.BrokerInterface) config.SetConfigFunc {
	return config.Broker(b)
}

func (i iConfig) Backend(b backends.BackendInterface) config.SetConfigFunc {
	return config.Backend(b)
}
func (i iConfig) Debug(debug bool) config.SetConfigFunc {
	return config.Debug(debug)
}

func (i iConfig) StatusExpires(ex int) config.SetConfigFunc {
	return config.StatusExpires(ex)
}

func (i iConfig) ResultExpires(ex int) config.SetConfigFunc {
	return config.ResultExpires(ex)
}

type iBackend struct {
}

// poolSize: ( default: numWorkers(for server) 10(for client) ) Maximum number of idle connections in the pool. if poolSize<=0 use default
func (i iBackend) NewRedisBackend(host string, port string, password string, db int, poolSize int) backends.RedisBackend {
	return backends.NewRedisBackend(host, port, password, db, poolSize)
}
