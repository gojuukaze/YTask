package main

import (
	"context"
	"github.com/gojuukaze/YTask/example/v1.1/server/workers"
	"github.com/gojuukaze/YTask/v1.1"
	"time"
)

type addArgs struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	b := ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 3)
	ser := ytask.Server.NewServer(
		ytask.Config.Broker(b),
		ytask.Config.Debug(true),
	)

	ser.Add("g1", "func-add", workers.AddFunc)
	ser.Add("g1", "struct-add", workers.AddStruct{})

	ser.Run("g1",2)
	time.Sleep(5*time.Second)
	ser.Shutdown(context.Background())

}
