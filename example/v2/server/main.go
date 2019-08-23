package main

import (
	"context"
	"fmt"
	"github.com/gojuukaze/YTask/example/v2/server/workers"
	"github.com/gojuukaze/YTask/v2"
)


func main() {
	b := ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 3)
	ser := ytask.Server.NewServer(
		ytask.Config.Broker(b),
		ytask.Config.Debug(true),
	)

	ser.Add("g1", "func-add", workers.AddFunc)
	ser.Add("g1", "struct-add", workers.AddStruct{})

	ser.Run("g1",2)
	fmt.Scanln()
	ser.Shutdown(context.Background())

}
