package config

import "YTask/v1/brokers"

type Config struct {
	Broker     brokers.BrokerInterface
	Debug      bool
}
