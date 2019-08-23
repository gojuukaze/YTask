package main

import (
	"github.com/gojuukaze/YTask/v2"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/server"

	"flag"
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

	//ser.Send("g1", "func-add", NumArgs{
	//	A: 123,
	//	B: 111,
	//})


	ser.Send("g1", "struct-add", NumArgs{
		A: 1,
		B: 1,
	})

	ser.Send("g1", "struct-add", NumArgs{
		A: 2,
		B: 2,
	})

	ser.Send("g1", "struct-add", NumArgs{
		A: 3,
		B: 3,
	})
	ser.Send("g1", "struct-add", NumArgs{
		A: 4,
		B: 4,
	})
	ser.Send("g1", "struct-add", NumArgs{
		A: 5,
		B: 5,
	})


}
