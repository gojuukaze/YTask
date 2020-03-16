package test

import (
	"fmt"
	"github.com/gojuukaze/YTask/v2/brokers"
	"github.com/gojuukaze/YTask/v2/controller"
	"github.com/gojuukaze/YTask/v2/message"
	"testing"
)

func TestRabbitmqBroker(t *testing.T) {
	b := brokers.NewRabbitMqBroker("127.0.0.1", "5672", "guest", "guest")
	var broker brokers.BrokerInterface = &b
	broker.Activate()
	msg := message.NewMessage(controller.NewTaskCtl())
	msg2 := message.NewMessage(controller.NewTaskCtl())

	err := broker.Send("test_redis", msg)
	if err != nil {
		t.Fatal(err)
	}
	err = broker.Send("test_redis", msg2)
	if err != nil {
		t.Fatal(err)
	}

	m, err := broker.Next("test_redis")
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%v", m) != fmt.Sprintf("%v", msg) {
		t.Fatalf("%v != %v", m, msg)
	}

	m2, err := broker.Next("test_redis")
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%v", m2) != fmt.Sprintf("%v", msg2) {
		t.Fatalf("%v != %v", m2, msg2)

	}
}
