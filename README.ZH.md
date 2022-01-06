# YTask
YTask is an asynchronous task queue for handling distributed jobs in golang  
golang异步任务/队列 框架  

* [中文文档](https://doc.ikaze.cn/YTask) (中文文档更加全面，优先阅读中文文档)
* [En Doc](https://github.com/gojuukaze/YTask/wiki)
* [Github](https://github.com/gojuukaze/YTask)

# install
```bash
go get github.com/gojuukaze/YTask/v2
```
# 架构图
<img src="./architecture_diagram.png" alt="architecture_diagram" width="80%">

# 特点
- 简单无侵入  
- 方便扩展broker，backend，logger
- 支持所有能被序列化为json的类型
- 支持任务重试，延时任务

# 快速开始

## server demo

```go
package main

import (
	"context"
	"github.com/gojuukaze/YTask/v2"
	"os"
	"os/signal"
	"syscall"
)

type User struct {
	Id   int
	Name string
}

func add(a int, b int) int {
	return a + b
}

func appendUser(user User, ids []int, names []string) []User {
	var r = make([]User, 0)
	r = append(r, user)
	for i := range ids {
		r = append(r, User{ids[i], names[i]})
	}
	return r
}

func main() {
	// clientPoolSize: Server端无需设置broker clientPoolSize
	broker := ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 0)

	// poolSize: 如果backend poolSize<=0 会使用默认值，
	//           对于server端backendPoolSize的默认值是 min(10, numWorkers)
	backend := ytask.Backend.NewRedisBackend("127.0.0.1", "6379", "", 0, 0)

	logger := ytask.Logger.NewYTaskLogger()  // v2.5+支持
	
	ser := ytask.Server.NewServer(
		ytask.Config.Broker(&broker),
		ytask.Config.Backend(&backend),
		ytask.Config.Logger(logger),		// 可以不设置 v2.5+支持
		ytask.Config.Debug(true),
		ytask.Config.StatusExpires(60*5),
		ytask.Config.ResultExpires(60*5),
	)

	ser.Add("group1", "add", add)
	ser.Add("group1", "append_user", appendUser)

	ser.Run("group1", 3)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ser.Shutdown(context.Background())

}
```

## client demo

```go
package main

import (
	"fmt"
	"github.com/gojuukaze/YTask/v2"
	"github.com/gojuukaze/YTask/v2/server"
	"time"
)

type User struct {
	Id   int
	Name string
}

var client server.Client

func main() {
	// 对于client你需要设置broker clientPoolSize
	broker := ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 5)
	// 对于client端，如果backend poolSize<=0，poolSize会设为10
	backend := ytask.Backend.NewRedisBackend("127.0.0.1", "6379", "", 0, 5)

	ser := ytask.Server.NewServer(
		ytask.Config.Broker(&broker),
		ytask.Config.Backend(&backend),
		ytask.Config.Debug(true),
		ytask.Config.StatusExpires(60*5),
		ytask.Config.ResultExpires(60*5),
	)

	client = ser.GetClient()

	// task add
	taskId, err := client.Send("group1", "add", 123, 44)
	_ = err
	result, err := client.GetResult(taskId, 2*time.Second, 300*time.Millisecond)
	_ = err

	if result.IsSuccess() {
		sum, err := result.GetInt64(0)
		// or
		var sum2 int
		err = result.Get(0, &sum2)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("add(123,44) =", int(sum))
	} else {
		fmt.Println("result failure")
	}
	// task append user
	taskId, _ = client.Send("group1", "append_user", User{1, "aa"}, []int{322, 11}, []string{"bb", "cc"})
	_ = err
	result, _ = client.GetResult(taskId, 2*time.Second, 300*time.Millisecond)
	var users []User
	result.Get(0, &users)
	fmt.Println(users)

}

```

# Example
[example](https://github.com/gojuukaze/YTask/tree/master/example/v2) 目录下有更多的样例可供参考
```bash
cd example/v2
go run server/main.go 

go run send/main.go
```

