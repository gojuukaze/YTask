package test

import (
	"fmt"
	"testing"

	"github.com/vua/YTask/v2/brokers"
	"github.com/vua/YTask/v2/controller"
	"github.com/vua/YTask/v2/message"
)

func TestRedisBroker(t *testing.T) {
	b := brokers.NewRedisBroker("127.0.0.1", "6379", "", 0, 1)
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

func TestRedisBrokerLSend(t *testing.T) {
	broker := brokers.NewRedisBroker("127.0.0.1", "6379", "", 0, 1)
	broker.Activate()
	msg := message.NewMessage(controller.NewTaskCtl())
	msg.Id = "1"
	msg2 := message.NewMessage(controller.NewTaskCtl())
	msg2.Id = "2"
	err := broker.Send("test_redis", msg)
	if err != nil {
		t.Fatal(err)
	}
	err = broker.LSend("test_redis", msg2)
	if err != nil {
		t.Fatal(err)
	}

	m, err := broker.Next("test_redis")
	if err != nil {
		t.Fatal(err)
	}
	if m.Id != msg2.Id {
		t.Fatalf("%v != %v", m, msg2)
	}

	m2, err := broker.Next("test_redis")
	if err != nil {
		t.Fatal(err)
	}
	if m2.Id != msg.Id {
		t.Fatalf("%v != %v", m2, msg)

	}
}
