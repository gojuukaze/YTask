package main

import (
	"YTask/example/service/workers"
	"YTask/v1/brokers/redisBroker"
	"YTask/v1/config"
	"YTask/v1/ytask"
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func addWorkers(t *ytask.YTask) {
	t.Add("ytask1", workers.AddWorker{})
	t.Add("ytask1", workers.SubWorker{})

	t.Add("ytask2", workers.MulWorker{})
}

func main() {

	groupName := flag.String("g", "ytask1", "start groupName")
	flag.Parse()


	var numWorkers = 3
	t := ytask.NewYTask(config.Config{
		Broker: redisBroker.NewRedisBroker("127.0.0.1", "6379", "", 0, numWorkers),
		Debug:  true,
	})

	addWorkers(&t)

	t.Run(*groupName, numWorkers)

	quit := make(chan os.Signal, 1)

	ctx := context.Background()

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	t.Shutdown(ctx)

}
