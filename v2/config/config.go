package config

import (
	"github.com/gojuukaze/YTask/v2/backends"
	"github.com/gojuukaze/YTask/v2/brokers"
	"github.com/gojuukaze/YTask/v2/log"
)

type Config struct {
	// require: true
	Broker brokers.BrokerInterface

	// require: false
	Backend backends.BackendInterface

	// require: false
	Logger log.LoggerInterface

	// require: false
	// default:false
	Debug bool

	// require: false
	// default: 1 day
	// task status expires in ex seconds, -1:forever
	StatusExpires int

	// require: false
	// default: 1 day
	// task result expires in ex seconds, -1:forever
	ResultExpires int

	EnableDelayServer    bool
	DelayServerQueueSize int
}

func (c Config) Clone() Config {
	newC := Config{
		Broker:               c.Broker.Clone(),
		Backend:              nil,
		Logger:               c.Logger.Clone(),
		Debug:                c.Debug,
		StatusExpires:        c.StatusExpires,
		ResultExpires:        c.ResultExpires,
		EnableDelayServer:    c.EnableDelayServer,
		DelayServerQueueSize: c.DelayServerQueueSize,
	}
	if c.Backend != nil {
		newC.Backend = c.Backend.Clone()
	}
	return newC

}

type Opt struct {
}
type SetConfigFunc func(*Config)

func NewConfig(setConfigFunc ...SetConfigFunc) Config {
	var config = Config{
		StatusExpires:        60 * 60 * 24,
		ResultExpires:        60 * 60 * 24,
		DelayServerQueueSize: 20,
		Logger: log.NewYTaskLogger(log.YTaskLog),
	}
	for _, f := range setConfigFunc {
		f(&config)
	}
	return config
}
func Broker(b brokers.BrokerInterface) SetConfigFunc {
	return func(config *Config) {
		config.Broker = b
	}
}

func Backend(b backends.BackendInterface) SetConfigFunc {
	return func(config *Config) {
		config.Backend = b
	}
}

func Logger(l log.LoggerInterface) SetConfigFunc {
	return func(config *Config) {
		config.Logger = l
	}
}

func DelayServerQueueSize(size int) SetConfigFunc {
	return func(config *Config) {
		config.DelayServerQueueSize = size
	}
}

func EnableDelayServer(enable bool) SetConfigFunc {
	return func(config *Config) {
		config.EnableDelayServer = enable
	}
}

func Debug(debug bool) SetConfigFunc {
	return func(config *Config) {
		config.Debug = debug
	}
}

func StatusExpires(ex int) SetConfigFunc {
	return func(config *Config) {
		config.StatusExpires = ex
	}
}

func ResultExpires(ex int) SetConfigFunc {
	return func(config *Config) {
		config.ResultExpires = ex
	}
}
