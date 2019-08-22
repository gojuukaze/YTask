package config

import "github.com/gojuukaze/YTask/v1.1/brokers"

type Config struct {
	Broker brokers.BrokerInterface
	Debug  bool
}

type Opt struct {
}
type SetConfigFunc func(*Config)

func Broker(b brokers.BrokerInterface) SetConfigFunc {
	return func(config *Config) {
		config.Broker = b
	}
}
func Debug(debug bool) SetConfigFunc {
	return func(config *Config) {
		config.Debug = debug
	}
}
func NewConfig(setConfigFunc ...SetConfigFunc) Config {
	var config = Config{}
	for _, f := range setConfigFunc {
		f(&config)
	}
	return config
}
