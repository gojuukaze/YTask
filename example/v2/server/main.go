package main

import (
	"context"
	workers2 "github.com/gojuukaze/YTask/example/v2/server/workers"
	"github.com/gojuukaze/YTask/v2"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// clientPoolSize: brokerPoolSize need not be set at server
	//                 -------------
	//                 server端不需要设置brokerPoolSize
	broker := ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 0)
	// poolSize: Maximum number of idle connections in the pool. If poolSize<=0 use default value
	//           default value is min(10, numWorkers) at server
	//           -------------
	//           如果poolSize<=0 会使用默认值，对于server端backendPoolSize的默认值是 min(10, numWorkers)
	backend := ytask.Backend.NewRedisBackend("127.0.0.1", "6379", "", 0, 0)

	ser := ytask.Server.NewServer(
		ytask.Config.Broker(&broker),
		ytask.Config.Backend(&backend),
		ytask.Config.Debug(true),
		ytask.Config.StatusExpires(60*5),
		ytask.Config.ResultExpires(60*5),
	)

	ser.Add("group1", "add", workers2.Add)
	ser.Add("group1", "retry", workers2.Retry)
	ser.Add("group1", "add_user", workers2.AppendUser)

	ser.Add("group2", "add_sub", workers2.AddSub)

	// v2.2开始支持运行多个group
	ser.Run("group1", 3)

	// If you want to use delayServer, set enableDelayServer
	// -------
	// 如果你要使用延时任务，把enableDelayServer设为true
	ser.Run("group2", 3,true)


	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ser.Shutdown(context.Background())

}
