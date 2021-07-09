快速开始
==========

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
       taskId, _ := client.Send("group1", "add", 123, 44)
       result, _ := client.GetResult(taskId, 2*time.Second, 300*time.Millisecond)

       if result.IsSuccess() {
           sum, _ := result.GetInt64(0)
           // or
           var sum2 int
           result.Get(0, &sum2)

           fmt.Println("add(123,44) =", int(sum))
       }

       // task append user
       taskId, _ = client.Send("group1", "append_user", User{1, "aa"}, []int{322, 11}, []string{"bb", "cc"})
       result, _ = client.GetResult(taskId, 2*time.Second, 300*time.Millisecond)
       var users []User
       result.Get(0, &users)
       fmt.Println(users)

   }
