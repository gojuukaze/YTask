YTask
-----------

YTask is an asynchronous task queue for handling distributed jobs in golang  
golang异步任务/队列 框架  

* [中文文档](https://doc.ikaze.cn/YTask) (Chinese document has more detailed instructions. If you know Chinese, read Chinese document)
* [En Doc](https://github.com/gojuukaze/YTask/wiki)
* [Github](https://github.com/gojuukaze/YTask)
* [Brokers And Backends](https://github.com/gojuukaze/YTask/drives)
* [V2 Doc](https://doc.ikaze.cn/YTaskV2)


# install
```bash
# install core
go get -u github.com/gojuukaze/YTask/v3

#install broker and backend
go get -u github.com/gojuukaze/YTask/drives/redis/v3
go get -u github.com/gojuukaze/YTask/drives/rabbitmq/v3
go get -u github.com/gojuukaze/YTask/drives/mongo2/v3
go get -u github.com/gojuukaze/YTask/drives/memcache/v3

```

# architecture diagram
<img src="./architecture_diagram.png" alt="architecture_diagram" width="80%">



# Quick Start

## server demo

```go
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"github.com/gojuukaze/YTask/drives/redis/v3"
	"github.com/gojuukaze/YTask/v3"
)

// Define two tasks.
// Task parameters and return values ​​support all types that can be serialized to json

func add(a int, b int) int {
	return a + b
}

type User struct {
	Id   int
	Name string
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
	// The last parameter of RedisBroker is the connection pool size. If you don't need workflow, you can set it to 0 (the default value is used when it is 0, and the default value on the server side is 3)
	// Otherwise, set it as needed, the maximum should not exceed the number of concurrent tasks
	broker := redis.NewRedisBroker("127.0.0.1", "6379", "", 0, 3)

	// The last parameter of RedisBackend is the connection pool size. For the server side, if <=0, the default value will be used.
	// the default value is min(10, numWorkers)
	backend := redis.NewRedisBackend("127.0.0.1", "6379", "", 0, 0)

	ser := ytask.Server.NewServer(
		ytask.Config.Broker(&broker),
		ytask.Config.Backend(&backend), // 可不设置
		ytask.Config.Debug(true),
		ytask.Config.StatusExpires(60*5),
		ytask.Config.ResultExpires(60*5),
	)

	// register task
	ser.Add("group1", "add", add)
	ser.Add("group1", "append_user", appendUser)

	// Run the server, the number of concurrency is 3
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
	"github.com/gojuukaze/YTask/drives/redis/v3"
	"github.com/gojuukaze/YTask/v3"
	"github.com/gojuukaze/YTask/v3/server"
	"time"
)

type User struct {
	Id   int
	Name string
}

var client server.Client

func main() {
	// For the client side you need to set the connection pool size
	broker := redis.NewRedisBroker("127.0.0.1", "6379", "", 0, 5)
	backend := redis.NewRedisBackend("127.0.0.1", "6379", "", 0, 5)

	ser := ytask.Server.NewServer(
		ytask.Config.Broker(&broker),
		ytask.Config.Backend(&backend),
		ytask.Config.Debug(true),
		ytask.Config.StatusExpires(60*5),
		ytask.Config.ResultExpires(60*5),
	)

	client = ser.GetClient()

	// Submit a task
	taskId, _ := client.Send("group1", "add", 123, 44)
	// get results
	result, _ := client.GetResult(taskId, 2*time.Second, 300*time.Millisecond)

	if result.IsSuccess() {
		// There are multiple ways to get the return value, please refer to the documentation for details
		sum, _ := result.GetInt64(0)
		// or
		var sum2 int
		result.Get(0, &sum2)

		fmt.Println("add(123,44) =", int(sum))
	}

	// Submit structs, slices, etc.
	taskId, _ = client.Send("group1", "append_user", User{1, "aa"}, []int{322, 11}, []string{"bb", "cc"})
	result, _ = client.GetResult(taskId, 2*time.Second, 300*time.Millisecond)
	var users []User
	result.Get(0, &users)
	fmt.Println(users)

}


```

# Example

Also take a look at [example](https://github.com/gojuukaze/YTask/tree/master/example) directory.

```bash
cd example/server
go run main.go 

cd example/send
go run send/main.go
```

捐赠 / Sponsor
================

开源不易，如果你觉得对你有帮助，求打赏个一块两块的

![](https://gitee.com/gojuukaze/liteAuth/raw/master/shang.jpg)





