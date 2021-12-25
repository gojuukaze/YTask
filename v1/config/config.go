package config

import "github.com/vua/YTask/v1/brokers"

type Config struct {
	Broker brokers.BrokerInterface
	Debug  bool
}
