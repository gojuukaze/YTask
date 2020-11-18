package test

import (
	"fmt"
	"github.com/gojuukaze/YTask/v2/brokers"
	"github.com/gojuukaze/YTask/v2/controller"
	"github.com/gojuukaze/YTask/v2/message"
	"testing"
)

func TestRocketMqBroker(t *testing.T) {

	broker := brokers.NewRocketMqBroker("127.0.0.1", "9876")
	broker.Activate()
	//broker.Shutdown()主要是为了关闭consumer,同步offset到broker
	//BUG：会出现同步失败
	defer broker.Shutdown()
	msg := message.NewMessage(controller.NewTaskCtl())
	msg2 := message.NewMessage(controller.NewTaskCtl())

	err := broker.Send("test_rock", msg)
	if err != nil {
		t.Fatal(err)
	}
	err = broker.Send("test_rock", msg2)
	if err != nil {
		t.Fatal(err)
	}
	m, err := broker.Next("test_rock")
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%v", m) != fmt.Sprintf("%v", msg) {
		t.Fatalf("%v != %v", m, msg)
	}

	m2, err := broker.Next("test_rock")
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%v", m2) != fmt.Sprintf("%v", msg2) {
		t.Fatalf("%v != %v", m2, msg2)

	}


}

func TestRocketMqBrokerLSend(t *testing.T) {
	broker := brokers.NewRocketMqBroker("127.0.0.1", "9876")
	broker.Activate()
	defer broker.Shutdown()
	msg := message.NewMessage(controller.NewTaskCtl())
	msg.Id = "1"
	msg2 := message.NewMessage(controller.NewTaskCtl())
	msg2.Id = "2"
	err := broker.Send("test_rock", msg)
	if err != nil {
		t.Fatal(err)
	}
	err = broker.LSend("test_rock", msg2)
	if err != nil {
		t.Fatal(err)
	}

	m, err := broker.Next("test_rock")
	if err != nil {
		t.Fatal(err)
	}
	if m.Id != msg2.Id {
		t.Fatalf("%v != %v", m, msg2)
	}

	m2, err := broker.Next("test_rock")
	if err != nil {
		t.Fatal(err)
	}
	if m2.Id != msg.Id {
		t.Fatalf("%v != %v", m2, msg)

	}
}
