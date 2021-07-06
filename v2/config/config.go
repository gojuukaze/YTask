package config

import (
	"github.com/gojuukaze/YTask/v2/backends"
	"github.com/gojuukaze/YTask/v2/brokers"
)

type Config struct {
	// require: true
	Broker brokers.BrokerInterface

	// require: false
	Backend backends.BackendInterface

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

	EnableDelayServer           bool
	DelayServerReadyMsgChanSize int
}

func (c Config) Clone() Config {
	newC := Config{
		Broker:                      c.Broker.Clone(),
		Backend:                     nil,
		Debug:                       c.Debug,
		StatusExpires:               c.StatusExpires,
		ResultExpires:               c.ResultExpires,
		EnableDelayServer:           c.EnableDelayServer,
		DelayServerReadyMsgChanSize: c.DelayServerReadyMsgChanSize,
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
		StatusExpires:               60 * 60 * 24,
		ResultExpires:               60 * 60 * 24,
		DelayServerReadyMsgChanSize: 20,
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

func DelayServerReadyMsgChanSize(size int) SetConfigFunc {
	return func(config *Config) {
		config.DelayServerReadyMsgChanSize = size
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
