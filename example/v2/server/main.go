package main

import (
	"context"
	"github.com/gojuukaze/YTask/example/v2/server/workers"
	"github.com/gojuukaze/YTask/v2"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// For the client, you need to set up the poolSize
	// Server端无需设置poolSize，
	broker := ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 0)
	backend := ytask.Backend.NewRedisBackend("127.0.0.1", "6379", "", 0, 0)

	ser := ytask.Server.NewServer(
		ytask.Config.Broker(&broker),
		ytask.Config.Backend(&backend),
		ytask.Config.Debug(true),
		ytask.Config.StatusExpires(60*5),
		ytask.Config.ResultExpires(60*5),
	)

	ser.Add("group1", "add", workers.Add)
	ser.Add("group1", "add_sub", workers.AddSub)
	ser.Add("group1", "retry", workers.Retry)
	ser.Add("group1", "add_user", workers.AppendUser)

	ser.Run("group1", 3)
	ser.Run("group1", 3)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ser.Shutdown(context.Background())

}
