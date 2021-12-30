package config

import "github.com/gojuukaze/YTask/v1/brokers"

type Config struct {
	Broker     brokers.BrokerInterface
	Debug      bool
}
