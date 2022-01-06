快速开始
==========

这里给出一个简单的server端，client端代码，更多的样例可以参考 `example/v2 <https://github.com/gojuukaze/YTask/tree/master/example/v2>`_


server
----------

.. code:: go

   package main

   import (
       "context"
       "github.com/gojuukaze/YTask/v2"
       "os"
       "os/signal"
       "syscall"
   )

   // 定义两个任务，任务参数、返回值支持所有能被序列化为json的类型

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
       // RedisBroker最后一个参数是连接池大小， Server端无需设置
       broker := ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 0)

       // RedisBackend最后一个参数是连接池大小，对于server端 如果<=0 会使用默认值，
       // 默认值是 min(10, numWorkers)
       backend := ytask.Backend.NewRedisBackend("127.0.0.1", "6379", "", 0, 0)

       logger := ytask.Logger.NewYTaskLogger()  // v2.5+支持

       ser := ytask.Server.NewServer(
           ytask.Config.Broker(&broker),
           ytask.Config.Backend(&backend),  // 可不设置
           ytask.Config.Logger(logger),  // 可不设置
           ytask.Config.Debug(true),
           ytask.Config.StatusExpires(60*5),
           ytask.Config.ResultExpires(60*5),
       )

       // 注册任务
       ser.Add("group1", "add", add)
       ser.Add("group1", "append_user", appendUser)

       // 运行server，并发数3
       ser.Run("group1", 3)

       quit := make(chan os.Signal, 1)

       signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
       <-quit
       ser.Shutdown(context.Background())

   }

client
----------

.. code:: go

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
       // 对于client端你需要设置连接池大小
       broker := ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 5)

       // 对于client端，如果连接池<=0，poolSize会默认为10
       backend := ytask.Backend.NewRedisBackend("127.0.0.1", "6379", "", 0, 5)

       logger := ytask.Logger.NewYTaskLogger()  // v2.5+支持

       ser := ytask.Server.NewServer(
           ytask.Config.Broker(&broker),
           ytask.Config.Backend(&backend),
           ytask.Config.Logger(logger),
           ytask.Config.Debug(true),
           ytask.Config.StatusExpires(60*5),
           ytask.Config.ResultExpires(60*5),
       )

       client = ser.GetClient()

       // 提交任务
       taskId, _ := client.Send("group1", "add", 123, 44)
       // 获取结果
       result, _ := client.GetResult(taskId, 2*time.Second, 300*time.Millisecond)

       if result.IsSuccess() {
           // 有多种方法获取返回值，具体可以参考： https://doc.ikaze.cn/YTask/client.html#id4
           sum, _ := result.GetInt64(0)
           // or
           var sum2 int
           result.Get(0, &sum2)

           fmt.Println("add(123,44) =", int(sum))
       }

       // 提交结构体，slice等
       taskId, _ = client.Send("group1", "append_user", User{1, "aa"}, []int{322, 11}, []string{"bb", "cc"})
       result, _ = client.GetResult(taskId, 2*time.Second, 300*time.Millisecond)
       var users []User
       result.Get(0, &users)
       fmt.Println(users)

   }
