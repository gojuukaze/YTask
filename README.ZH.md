# YTask
YTask is an asynchronous task queue for handling distributed jobs in golang  
golang异步任务/队列 框架  

* [中文文档](README.ZH.md) (中文文档更加全面，优先阅读中文文档)
* [En Doc](README.md)
* [Github](https://github.com/gojuukaze/YTask)

# install
```bash
go get github.com/gojuukaze/YTask/v2
```
# 架构图
<img src="./architecture_diagram.png" alt="architecture_diagram" width="75%">

# todo
- [x] save result  
- [x] task retry  
- [x] 支持 RabbitMQ
- [ ] 一次运行多了group
- [ ] 扩展TaskCtl参数
- [x] 支持更多类型

# 文档
* [快速开始](#快速开始)
  * [server样例](#server-demo)
  * [client样例](#client-demo)
  * [其他样例](#other-example)
* [使用指南](#使用指南)
  * [服务端](#服务端)
    * [服务端配置](#服务端配置)
    * [注册任务](#注册任务)
    * [运行与停止](#运行与停止)
  * [客户端](#客户端)
    * [获取连接](#获取连接)
    * [发送信息](#发送信息)
    * [获取结果](#获取结果)
  * [重试](#重试)
    * [设置重试次数](#设置重试次数)
    * [禁用重试](#禁用重试)
  * [broker](#broker)
    * [redis broker](#redisbroker)
    * [rabbitmq broker](#rabbitmqbroker)
    * [自定义broker](#自定义broker)
  * [backend](#backend)
    * [redis backend](#redisbackend)
    * [memcache backend](#memcachebackend)
    * [mongo backend](#mongobackend)
    * [自定义backend](#自定义backend)
  * [支持的类型](#支持的类型)
  * [log](#log)
  * [error](#error)



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
[example](https://github.com/gojuukaze/YTask/tree/master/example/v2) 目录下有更多的样例可供参考
```bash
cd example/v2
go run server/main.go 

go run send/main.go
```


# 使用指南

## 服务端
* 使用`NewServer()`初始化服务，其参数是server的配置，所有配置在下面
```go
import "github.com/gojuukaze/YTask/v2"

ser := ytask.Server.NewServer(
		ytask.Config.Broker(&broker),
		ytask.Config.Backend(&backend),
		...
)
```
### 服务端配置
| 配置        | 是否必须 | 默认值 | code                         | 说明                                            |
|---------------|---------|---------|------------------------------|--------------------------------------------------|
| Broker        | \*      |         | ytask\.Config\.Broker        |                                                  |
| Backend       |         | nil     | ytask\.Config\.Backend       |                                                  |
| Debug         |         | FALSE   | ytask\.Config\.Debug         |   是否开启debug                      |
| StatusExpires |         | 1day    | ytask\.Config\.StatusExpires | 单位：秒，任务状态的过期时间, -1:永久保存 |
| ResultExpires |         | 1day    | ytask\.Config\.ResultExpires | 单位：秒，任务结果的过期时间, -1:永久保存 |

* 任务状态、结果有什么不同？
  * 状态： 任务的开始、运行、成功、失败状态
  * 结果： 函数的返回值

### 注册任务
使用`Add`注册任务
```go
// group1 : 任务所属组，也是队列的名字
// add : 任务名 
// addFunc : 任务函数
ser.Add("group1","add",addFunc)
```
任务函数的参数，返回值支持int,float等类型，如果需要在函数中控制任务的重试等东西，则函数的第一个参数为`TaskCtl`，如：

```go
func Add(ctl *controller.TaskCtl, a int, b int) (int, int) {
	if ... {
        // retry
		ctl.Retry(errors.New("ctl.Retry"))
        return 0, 0
	}

	return a + b, a - b
}
```

带`TaskCtl`函数和其他函数使用上没区别

### 运行与停止
```go
// group1 : 运行的组名
// 3 : 并发任务数
ser.Run("group1",3)

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
ser.Shutdown(context.Background())
```
> 你不能用同一个server运行不同的组，比如：
> ```go
> ser:=ytask.Server.NewServer(...)
> ser.Run("g1",1)
> // 这样会报错
> ser.Run("g2",1)
> ``` 
> 这个功能会在接下来的版本中加入

## 客户端

### 获取连接
获取连接前一样需要初始化Server，然后调用`GetClient()`。`NewServer`的参数可以和服务端不同，但建议使用相同的参数配置
```go
import "github.com/gojuukaze/YTask/v2"

ser := ytask.Server.NewServer(
		ytask.Config.Broker(&broker),
		ytask.Config.Backend(&backend),
		...
)

client = ser.GetClient()
```

### 发送信息
使用`Send`发送任务信息，函数前两个参数为组名、任务，后面的参数是任务函数的参数。函数第一个返回值为任务id，可以用来获取任务结果。   
发送消息时可以使用`SetTaskCtl()`配置该次任务的重试次数等
```go
// group1 : 组名
// add : 任务名
// 12,33 ... : 任务参数
// return :
//   - taskId : taskId
//   - err : error
taskId,err:=client.Send("group1","add",12,33)

// set retry count
taskId,err=client.SetTaskCtl(client.RetryCount, 5).Send("group1","add",12,33)

```

### 获取结果
调用`GetResult()`获取任务结果，第2个参数为超时时间，第3个参数为重新获取时间。  
获取结果后可调用`GetXX()`，`Get()`，`Gets()`获取任务函数的返回结果。
```go
// taskId :
// 3*time.Second : timeout
// 300*time.Millisecond : sleep time
result, _ := client.GetResult(taskId, 3*time.Second, 300*time.Millisecond)

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

> **重要！！**  
> YTask虽然提供获取结果功能，但不要过渡依赖。  
> 如果backend出错导致无法保存结果，YTask不会再次重试。因为对任务状态、结果的保存与运行任务的goroutine是同一个，不断重试会导致任务无法开始或无法结束。
> YTask优先保障任务运行，而不是结果保存。  
> 如果你特别需要任务结果，推荐你在任务函数中自行保存。

## 重试
**默认的重试次数是3次**  

有两种方法可以触发2重试
* 使用 panic
```go

func add(a, b int){
    panic("xx")
}
```

* 使用 TaskCtl
```go

func add(ctl *controller.TaskCtl,a, b int){
    ctl.Retry(errors.New("xx"))
    return
}
```

### 设置重试次数

* 目前只支持在client端设置
```go
client.SetTaskCtl(client.RetryCount, 5).Send("group1", "retry", 123, 44)
```

### 禁用重试
* 在server端针对某个任务禁用
```go
func add(ctl *controller.TaskCtl,a, b int){
    ctl.SetRetryCount(0)
    return
}
```
* 在client端对此次任务禁用
```go
client.SetTaskCtl(client.RetryCount, 0).Send("group1", "retry", 123, 44)
```

## broker
YTask使用broker与任务队列通信，发送或接收任务。  
支持的broker有：
### redisBroker

```go
import "github.com/gojuukaze/YTask/v2"

// 127.0.0.1 : host
// 6379 : port
// "" : password
// 0 : db
// 10 : 连接池大小. 
//      对于server端，你无需自定义连接池，如果连接池为0，系统会自动设置合适的连接池
//      对于client端, 你需要根据情况自行设置连接池
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

### 自定义broker
你可以自行定义broker。需要注意，因为系统中会调用`SetPoolSize`设置连接池，所以初始化broker时不要建立连接，调用`Activate()`时再建立。
如果你的broker不支持连接池，那可以不用管Activate,SetPoolSize,GetPoolSize三个方法，直接返回空就行。
```go
type BrokerInterface interface {
    // 获取任务
	Next(queryName string) (message.Message, error)
    // 发送任务
	Send(queryName string, msg message.Message) error
	// 建立连接
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
// 10 : 连接池大小. 
//      对于server端，你无需自定义连接池，如果连接池为0，系统会自动设置合适的连接池
//      对于client端, 你需要根据情况自行设置连接池
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

### 自定义backend

你可以自行定义backend。同broker一样，调用`Activate()`时再建立连接。

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

## 支持的类型
支持所有能序列化为json格式的类型

## log
YTask使用logrus打印日志，下面给出了一个输出日志到文件的样例

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

* [go-watch-file](https://github.com/gojuukaze/go-watch-file) ：一个专为日志系统编写的读写文件库，会自动监听文件的变化，文件被删除时自动创建新文件。

## error
内置的错误类型
```go
const (
	ErrTypeEmptyQuery      = 1
	ErrTypeUnsupportedType = 2
	ErrTypeOutOfRange      = 3
	ErrTypeNilResult       = 4
	ErrTypeTimeOut         = 5
)
```

比较错误
```go
import 	"github.com/gojuukaze/YTask/v2/yerrors"
yerrors.IsEqual(err, yerrors.ErrTypeNilResult)

```
