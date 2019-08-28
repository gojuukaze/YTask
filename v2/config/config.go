package config

import (
	"github.com/gojuukaze/YTask/v2/backends"
	"github.com/gojuukaze/YTask/v2/brokers"
)

type Config struct {
	Broker        brokers.BrokerInterface
	Backend       backends.BackendInterface
	Debug         bool
	StatusExpires int // second, -1:forever , default: 1 day
	ResultExpires int // second, -1:forever , default: 1 day
}

type Opt struct {
}
type SetConfigFunc func(*Config)

func NewConfig(setConfigFunc ...SetConfigFunc) Config {
	var config = Config{
		StatusExpires: 60 * 60 * 24,
		ResultExpires: 60 * 60 * 24,
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
