package test

import (
	"fmt"
	"github.com/gojuukaze/YTask/v3/core/message"
	"github.com/gojuukaze/YTask/v3/core/yerrors"
	"github.com/gojuukaze/YTask/v3/drives/rabbitmq"
	"github.com/gojuukaze/YTask/v3/taskMessage"
	"testing"
)

func TestRabbitmqBroker(t *testing.T) {
	broker := rabbitmq.NewRabbitMqBroker("127.0.0.1", "5672", "guest", "guest", "", 2)
	broker.Activate()
	msg := message.NewMessage(taskMessage.NewTaskCtl())
	msg2 := message.NewMessage(taskMessage.NewTaskCtl())

	_, err := broker.Next("test_amqp")
	if !yerrors.IsEqual(err, yerrors.ErrTypeEmptyQueue) {
		t.Fatal(err)
	}

	err = broker.Send("test_amqp", msg)
	if err != nil {
		t.Fatal(err)
	}
	err = broker.Send("test_amqp", msg2)
	if err != nil {
		t.Fatal(err)
	}

	m, err := broker.Next("test_amqp")
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%v", m) != fmt.Sprintf("%v", msg) {
		t.Fatalf("%v != %v", m, msg)
	}

	m2, err := broker.Next("test_amqp")
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%v", m2) != fmt.Sprintf("%v", msg2) {
		t.Fatalf("%v != %v", m2, msg2)

	}
}

func TestRabbitmqBrokerLSend(t *testing.T) {
	broker := rabbitmq.NewRabbitMqBroker("127.0.0.1", "5672", "guest", "guest", "", 2)
	broker.Activate()
	msg := message.NewMessage(taskMessage.NewTaskCtl())
	msg.Id = "1"
	msg2 := message.NewMessage(taskMessage.NewTaskCtl())
	msg2.Id = "2"
	err := broker.Send("test_amqp", msg)
	if err != nil {
		t.Fatal(err)
	}
	err = broker.LSend("test_amqp", msg2)
	if err != nil {
		t.Fatal(err)
	}

	m, err := broker.Next("test_amqp")
	if err != nil {
		t.Fatal(err)
	}
	if m.Id != msg2.Id {
		t.Fatalf("%v != %v", m, msg2)
	}

	m2, err := broker.Next("test_amqp")
	if err != nil {
		t.Fatal(err)
	}
	if m2.Id != msg.Id {
		t.Fatalf("%v != %v", m2, msg)

	}
}
