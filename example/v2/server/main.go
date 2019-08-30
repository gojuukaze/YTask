package main

import (
	"context"
	"fmt"
	"github.com/gojuukaze/YTask/example/v2/server/workers"
	"github.com/gojuukaze/YTask/v2"
)

func main() {
	b := ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 0)
	b2 := ytask.Backend.NewRedisBackend("127.0.0.1", "6379", "", 0, 0)

	ser := ytask.Server.NewServer(
		ytask.Config.Broker(&b),
		ytask.Config.Backend(&b2),
		ytask.Config.Debug(true),
	)

	ser.Add("g1", "add", workers.Add)

	ser.Run("g1", 5)
	fmt.Scanln()
	ser.Shutdown(context.Background())

}
