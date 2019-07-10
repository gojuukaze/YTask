package main

import (
	"YTask/v1/brokers/redisBroker"
	"YTask/v1/config"
	"YTask/v1/ymsg"
	"YTask/v1/ytask"
	"flag"
	"time"
)

type NumArgs struct {
	A int `json:"a"`
	B int `json:"b"`
}

func send1(t ytask.YTask, groupName string) {
	go func() {
		t.Send(groupName, ymsg.Message{
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

func send2(t ytask.YTask, groupName string) {
	go func() {
		t.Send(groupName, ymsg.Message{
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
	groupName := flag.String("g", "ytask1", "start groupName")
	flag.Parse()

	t := ytask.NewYTask(config.Config{
		Broker: redisBroker.NewRedisBroker("127.0.0.1", "6379", "", 0, 3),
		Debug:  true,
	})

	if *groupName=="ytask1" {
		send1(t, *groupName)
	}else if *groupName=="ytask2" {
		send2(t, *groupName)
	}

	time.Sleep(2*time.Second)




}
