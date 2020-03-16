# YTask
YTask is an asynchronous task queue for handling distributed jobs in golang  
golang异步任务/队列 框架  

* [中文文档](README.ZH.md) (Chinese document has more detailed instructions. If you know Chinese, read Chinese document)
* [En Doc](README.md)
* [Github](https://github.com/gojuukaze/YTask)

# install
```bash
go get github.com/gojuukaze/YTask/v2
```
# architecture diagram
<img src="./architecture_diagram.png" alt="architecture_diagram" width="75%">

# todo
- [x] save result  
- [x] task retry  
- [ ] run multi group
- [ ] more option in TaskCtl
- [x] support more type

# doc

* [Quick Start](#quick-start)
  * [server](#server-demo)
  * [client](#client-demo)
  * [other example](#other-example)
* [Usage](#usage)
  * [server](#server)
    * [server config](#server-config)
    * [add worker](#add-worker-func)
    * [run and shutdown](#run-and-shutdown)
  * [client](#client)
    * [get client](#get-client)
    * [send message](#send-msg)
    * [get result](#get-result)
  * [retry](#retry)
    * [set retry count](#set-retry-count)
    * [disable retry](#disable-retry)
  * [broker](#broker)
    * [redis broker](#redisbroker)
    * [rabbitmq broker](#rabbitmqbroker)
    * [custom broker](#custom-broker)
  * [backend](#backend)
    * [redis backend](#redisbackend)
    * [memcache backend](#memcachebackend)
    * [mongo backend](#mongobackend)
    * [custom backend](#custom-backend)
  * [support type](#support-type)
  * [log](#log)
  * [error](#error)



# Quick Start

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

func add(a int,b int)int {
    return a+b
}

func appendUser(user User, ids []int, names []string) []User {
	var r = make([]User, 0)
	r = append(r, user)
	for i := range ids {
		r = append(r, User{ids[i],names[i],})
	}
	return r
}

func main() {
	// For the server, you do not need to set up the poolSize
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
	// For the client, you need to set up the poolSize
	// 对于client你需要设置poolSize
	broker := ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 5)
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

## other example
Also take a look at [example](https://github.com/gojuukaze/YTask/tree/master/example/v2) directory
```bash
cd example/v2
go run server/main.go 

go run send/main.go
```


# usage

## server
* init
```go
import "github.com/gojuukaze/YTask/v2"

ser := ytask.Server.NewServer(
		ytask.Config.Broker(&broker),
		ytask.Config.Backend(&backend),
		...
)
```
### server config
| Config        | require | default | code                         | other                                            |
|---------------|---------|---------|------------------------------|--------------------------------------------------|
| Broker        | \*      |         | ytask\.Config\.Broker        |                                                  |
| Backend       |         | nil     | ytask\.Config\.Backend       |                                                  |
| Debug         |         | FALSE   | ytask\.Config\.Debug         |                                                  |
| StatusExpires |         | 1day    | ytask\.Config\.StatusExpires | "task status expires in ex seconds, \-1:forever" |
| ResultExpires |         | 1day    | ytask\.Config\.ResultExpires | "task result expires in ex seconds, \-1:forever" |

* StatusExpires, ResultExpires is not valid for Mongo backend, 0 means no storage, > 0 means permanent storage
  

### add worker func

```go
// group1 : group name is also the query name
// add : worker name 
// addFunc : worker func
ser.Add("group1","add",addFunc)
```

### run and shutdown
```go
// group1 : run group name
// 3 : number of worker goroutine
ser.Run("group1",3)

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
ser.Shutdown(context.Background())
```
> You cannot run multiple groups with the same server.
> ```go
> ser:=ytask.Server.NewServer(...)
> ser.Run("g1",1)
> // panic
> ser.Run("g2",1)
> ``` 
> This feature is already under development

## client

### get client
```go
import "github.com/gojuukaze/YTask/v2"

ser := ytask.Server.NewServer(
		ytask.Config.Broker(&broker),
		ytask.Config.Backend(&backend),
		...
)

client = ser.GetClient()
```

### send msg
```go
// group1 : group name
// add : worker name
// 12,33 ... : func args
// return :
//   - taskId : taskId
//   - err : error
taskId,err:=client.Send("group1","add",12,33)

// set retry count
taskId,err=client.SetTaskCtl(client.RetryCount, 5).Send("group1","add",12,33)

```

### get result
```go
// taskId :
// 3*time.Second : timeout
// 300*time.Millisecond : sleep time
result, _ := client.GetResult(taskId, 3*time.Second, 300*time.Millisecond)

// get worker func return
if result.IsSuccess(){
    // get worker func return
    a,err:=result.GetInt64(0)
    b,err:=result.GetBool(1)
    
    // or
    var a int
    var b bool
    err:=result.Get(0, &a)
    err:=result.Get(1, &b)

    // or
    var a int
    var b bool
    err:=result.Gets(&a, &b)
}
```
> **Warning!!!**  
> Although YTask provides the ability to get results, don't rely on transitions.  
> If the backend error causes the result to not be saved, YTask will not retry again. 
> Keep retrying will cause the task to fail to start or end.  
> If you need task results in particular, it is recommended that you save them yourself in the task function.

## retry
**default retry count is 3**  

there are 2 way to trigger retry
* use panic
```go

func add(a, b int){
    panic("xx")
}
```

* use TaskCtl
```go

func add(ctl *controller.TaskCtl,a, b int){
    ctl.Retry(errors.New("xx"))
    return
}
```

### set retry count

* in client
```go
client.SetTaskCtl(client.RetryCount, 5).Send("group1", "retry", 123, 44)
```

### disable retry
* in server
```go
func add(ctl *controller.TaskCtl,a, b int){
    ctl.SetRetryCount(0)
    return
}
```
* in client
```go
client.SetTaskCtl(client.RetryCount, 0).Send("group1", "retry", 123, 44)
```

## broker

### redisBroker

```go
import "github.com/gojuukaze/YTask/v2"

// 127.0.0.1 : host
// 6379 : port
// "" : password
// 0 : db
// 10 : connection pool size. 
//      For server, if poolSize is 0, the pool size will be set automatically.
//      For client, you need to set up the poolSize by yourself
ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 10)
```

### rabbitMqBroker

```go
import "github.com/gojuukaze/YTask/v2"
// 127.0.0.1 : host
// 5672 : port
// guest : username
// guest : password

ytask.Broker.NewRabbitMqBroker("127.0.0.1", "5672", "guest", "guest")
```

### custom broker
```go
type BrokerInterface interface {
    // get task
	Next(queryName string) (message.Message, error)
    // send task
	Send(queryName string, msg message.Message) error
	// Activate connection
	Activate()
	SetPoolSize(int)
	GetPoolSize()int
}
```

## backend

### redisBackend

```go
import "github.com/gojuukaze/YTask/v2"

// 127.0.0.1 : host
// 6379 : port
// "" : password
// 0 : db
// 10 : connection pool size. 
//      For server, if poolSize is 0, the pool size will be set automatically.
//      For client, you need to set up the poolSize by yourself

ytask.Backend.NewRedisBackend("127.0.0.1", "6379", "", 0, 10)
```

### memCacheBackend

```go
import "github.com/gojuukaze/YTask/v2"

// 127.0.0.1 : host
// 11211 : port
// 10 : connection pool size. 

ytask.Backend.NewMemCacheBackend("127.0.0.1", "11211", 10)
```

### mongoBackend

```go
import "github.com/gojuukaze/YTask/v2"

// 127.0.0.1 : host
// 27017 : port
// "" : username
// "" : password
// "task": db
// "taks": collection

ytask.Backend.NewMongoBackend("127.0.0.1", "27017", "", "", "task", "task")
```

### custom backend

```go
type BackendInterface interface {
	SetResult(result message.Result, exTime int) error
	GetResult(key string) (message.Result, error)
	// Activate connection
	Activate()
	SetPoolSize(int)
	GetPoolSize() int
}
```

## support type
Support all types what can be serialized to JSON


## log


```go
import (
"github.com/gojuukaze/YTask/v2/log"
"github.com/gojuukaze/go-watch-file")

// write to file
file,err:=watchFile.OpenWatchFile("xx.log")
if err != nil {
	panic(err)
}
log.YTaskLog.SetOutput(file)

// set level
log.YTaskLog.SetLevel(logrus.InfoLevel)
```

## error
error type
```go
const (
	ErrTypeEmptyQuery      = 1
	ErrTypeUnsupportedType = 2
	ErrTypeOutOfRange      = 3
	ErrTypeNilResult       = 4
	ErrTypeTimeOut         = 5
)
```

compare err
```go
import 	"github.com/gojuukaze/YTask/v2/yerrors"
yerrors.IsEqual(err, yerrors.ErrTypeNilResult)

```
