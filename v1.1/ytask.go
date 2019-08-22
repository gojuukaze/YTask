package ytask

import (
	"github.com/gojuukaze/YTask/v1.1/brokers"
	"github.com/gojuukaze/YTask/v1.1/brokers/redisBroker"
	"github.com/gojuukaze/YTask/v1.1/config"
	"github.com/gojuukaze/YTask/v1.1/server"
)

var (
	Server = iServer{}
	Broker = iBroker{}
	Config = iConfig{}
)

type iServer struct {
}

func (is iServer) NewServer(setConfigFunc ...config.SetConfigFunc) server.Server {
	c := config.NewConfig(setConfigFunc...)
	return server.NewServer(c)
}

type iBroker struct {
}

func (i iBroker) NewRedisBroker(host string, port string, password string, db int, numConns int) redisBroker.RedisBroker {
	return redisBroker.NewRedisBroker(host, port, password, db, numConns)
}

type iConfig struct {
}

func (i iConfig) Broker(b brokers.BrokerInterface) config.SetConfigFunc {
	return config.Broker(b)
}

func (i iConfig) Debug(debug bool) config.SetConfigFunc {
	return config.Debug(debug)
}



