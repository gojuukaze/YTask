package test

import (
	"github.com/gojuukaze/YTask/v2/brokers"
	"github.com/gojuukaze/YTask/v2/controller"
	"github.com/gojuukaze/YTask/v2/message"
	"testing"
)

func TestRedisBroker(t *testing.T) {
	b := brokers.NewRedisBroker("127.0.0.1", "6379", "", 0, 1)
	var broker brokers.BrokerInterface = &b
	broker.Activate()
	msg:=message.NewMessage(controller.NewTaskCtl())
	msg2:=message.NewMessage(controller.NewTaskCtl())

	err:=broker.Send("test_redis",msg)
	if err!=nil{
		t.Fatal(err)
	}
	err=broker.Send("test_redis",msg2)
	if err!=nil{
		t.Fatal(err)
	}



	m,err:=broker.Next("test_redis")
	if err!=nil{
		t.Fatal(err)
	}
	if m!=msg{
		t.Fatalf("%v != %v",m,msg)
	}

	m2,err:=broker.Next("test_redis")
	if err!=nil{
		t.Fatal(err)
	}
	if m2!=msg2{
		t.Fatalf("%v != %v",m2,msg2)

	}
}
