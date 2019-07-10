# YTask
YTask is an asynchronous task queue for handling distributed jobs in golang

# install
```bash
go get github.com/gojuukaze/YTask
```

# todo
[] save result  
[] task retry  
[] support amqp

#example

##server

```go
package main

import (
	"github.com/gojuukaze/YTask/v1/brokers/redisBroker"
	"github.com/gojuukaze/YTask/v1/config"
	"github.com/gojuukaze/YTask/v1/ymsg"
	"github.com/gojuukaze/YTask/v1/ytask"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type NumArgs struct {
	A int `json:"a"`
	B int `json:"b"`
}

type AddWorker struct {
}

func (a AddWorker) Name() string {
	return "add"
}

func (a AddWorker) Run(msg ymsg.Message) error {
	var args NumArgs
	_ = json.Unmarshal([]byte(msg.JsonArgs), &args)

	fmt.Println(args.A + args.B)
	return nil
}

func main() {

	var numWorkers = 3
	t := ytask.NewYTask(config.Config{
		Broker: redisBroker.NewRedisBroker("127.0.0.1", "6379", "", 0, numWorkers),
		Debug:  true,
	})

	t.Add("ytask", AddWorker{})

	t.Run("ytask", numWorkers)

	quit := make(chan os.Signal, 1)

	ctx := context.Background()

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	t.Shutdown(ctx)

}

```

##worker

```go
package main

import (
	"github.com/gojuukaze/YTask/v1/brokers/redisBroker"
	"github.com/gojuukaze/YTask/v1/config"
	"github.com/gojuukaze/YTask/v1/ymsg"
	"github.com/gojuukaze/YTask/v1/ytask"
)

type NumArgs struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {

	t := ytask.NewYTask(config.Config{
		Broker: redisBroker.NewRedisBroker("127.0.0.1", "6379", "", 0, 3),
		Debug:  true,
	})

	t.Send("ytask", ymsg.Message{
		WorkerName: "add",
		JsonArgs:   `{"a":1,"b":2}`,
	})

	t.Send("ytask", "add", `{"a":1,"b":2}`)

	t.Send("ytask", "add", NumArgs{3, 1})

}

```

## other example
Also take a look at [example](https://github.com/gojuukaze/YTask/tree/master/example) directory
```bash
go run example/server/main.go -g ytask1

go run example/send/main.go -g ytask1
```

```bash
go run example/server/main.go -g ytask2

go run example/send/main.go -g ytask2
```