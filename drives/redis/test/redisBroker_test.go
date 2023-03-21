package test

import (
	"fmt"
	"github.com/gojuukaze/YTask/drives/redis/v3"
	"github.com/gojuukaze/YTask/v3/brokers"
	"github.com/gojuukaze/YTask/v3/message"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestRedisBroker(t *testing.T) {
	b := redis.NewRedisBroker([]string{"127.0.0.1:6379"}, "", 0, 1, 0)
	var broker brokers.BrokerInterface = &b
	broker.Activate()
	msg := message.NewMessage(message.NewMsgArgs())
	msg2 := message.NewMessage(message.NewMsgArgs())

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
	broker := redis.NewRedisBroker([]string{"127.0.0.1:6379"}, "", 0, 1, 0)
	broker.Activate()
	msg := message.NewMessage(message.NewMsgArgs())
	msg.Id = "1"
	msg2 := message.NewMessage(message.NewMsgArgs())
	msg2.Id = "2"
	err := broker.Send("test_redis", msg)
	if err != nil {
		t.Fatal(err)
	}
	err = broker.LSend("test_redis", msg2)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("重启Redis后通过 \" kill -CONT %d \" 继续运行测试\n", os.Getpid())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGCONT)
	<-quit

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

func TestRedisBroker2(t *testing.T) {
	broker := redis.NewRedisBroker([]string{"127.0.0.1:6379"}, "", 0, 1, 0)
	broker.Activate()
	msg := message.NewMessage(message.NewMsgArgs())
	msg.Id = "1"
	msg2 := message.NewMessage(message.NewMsgArgs())
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
