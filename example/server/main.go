package main

import (
	"context"
	"github.com/gojuukaze/YTask/v3/core"
	"github.com/gojuukaze/YTask/v3/drives/redis"
	"github.com/gojuukaze/YTask/v3/example/server/workers"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//
	// PoolSize : server端, 如果brokerPoolSize<=0时默认为3;
	//              如果需要频繁使用工作流，则可适当调大此项，最大不要超过 并发任务数+1
	broker := redis.NewRedisBroker("127.0.0.1", "6379", "", 0, 0)
	// poolSize: Maximum number of idle connections in the pool. If poolSize<=0 use default value
	//           default value is min(10, numWorkers) at server
	//           -------------
	//           如果poolSize<=0 会使用默认值，对于server端backendPoolSize的默认值是 min(10, numWorkers)
	backend := redis.NewRedisBackend("127.0.0.1", "6379", "", 0, 0)

	ser := ytask.Server.NewServer(
		ytask.Config.Broker(&broker),
		ytask.Config.Backend(&backend),
		ytask.Config.Debug(true),
		ytask.Config.StatusExpires(60*5),
		ytask.Config.ResultExpires(60*5),
	)

	ser.Add("group1", "add", workers.Add)
	ser.Add("group1", "retry", workers.Retry)
	ser.Add("group1", "add_user", workers.AppendUser)

	ser.Add("group2", "add_sub", workers.AddSub)

	ser.Run("group1", 3)
	ser.Run("group2", 3, true)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ser.Shutdown(context.Background())

}
