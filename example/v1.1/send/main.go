package main

import (
	"github.com/gojuukaze/YTask/v1.1"
	"github.com/gojuukaze/YTask/v1.1/message"
	"github.com/gojuukaze/YTask/v1.1/server"

	"flag"
	"time"
)

type NumArgs struct {
	A int `json:"a"`
	B int `json:"b"`
}

func send1(t server.Server, groupName string) {
	go func() {
		t.Send(groupName, message.Message{
			WorkerName: "add",
			JsonArgs:   `{"a":1,"b":2}`,
		})
	}()

	go func() {
		t.Send(groupName, "sub", `{"a":1,"b":2}`)
	}()

	go func() {
		t.Send(groupName, "sub", NumArgs{3, 1})
	}()

}

func send2(t server.Server, groupName string) {
	go func() {
		t.Send(groupName, message.Message{
			WorkerName: "mul",
			JsonArgs:   `{"a":1,"b":2}`,
		})
	}()

	go func() {
		t.Send(groupName, "mul", `{"a":1,"b":2}`)
	}()

	go func() {
		t.Send(groupName, "mul", NumArgs{3, 1})
	}()

}

func main() {
	flag.Parse()

	b := ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 3)
	ser := ytask.Server.NewServer(
		ytask.Config.Broker(b),
		ytask.Config.Debug(true),
	)

	ser.Send("g1", "func-add", NumArgs{
		A: 123,
		B: 111,
	})
	ser.Send("g1", "struct-add", NumArgs{
		A: 0,
		B: 111,
	})

	time.Sleep(2 * time.Second)

}
