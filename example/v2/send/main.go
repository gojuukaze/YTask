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

func taskAdd() {
	taskId, err := client.Send("group1", "add", 123, 44)
	_ = err
	result, err := client.GetResult(taskId, 2*time.Second, 300*time.Millisecond)
	_ = err

	if result.IsSuccess() {
		sum, err := result.GetInt64(0)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("add(123,44) =", int(sum))
	} else {
		fmt.Println("result failure")
	}
}

func taskAddSub() {
	taskId, _ := client.Send("group2", "add_sub", 123, 44)
	result, _ := client.GetResult(taskId, 2*time.Second, 300*time.Millisecond)

	if result.IsSuccess() {
		sum, _ := result.GetInt64(0)
		sub, _ := result.GetInt64(1)

		// or
		var sum2, sub2 int
		result.Get(0, &sum2)
		result.Get(1, &sub2)

		// or
		var sum3, sub3 int
		result.Gets(&sum3, &sub3)

		fmt.Println("add_sub(123,44) =", sum, sub)
	}
}

func taskAppendUser() {
	taskId, err := client.Send("group1", "add_user", User{1, "aa"}, []int{322, 11}, []string{"bb", "cc"})
	_ = err
	result, err := client.GetResult(taskId, 2*time.Second, 300*time.Millisecond)
	_ = err

	if result.IsSuccess() {
		var users []User
		err := result.Get(0, &users)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(`add_user({1,"aa"}, [322,11], ["bb","cc"]) =`, users)
	}
}
func sendAndGet() {
	taskAdd()
	taskAddSub()
	taskAppendUser()
}

func retry() {

	// set retry count, default is 3
	tId, _ := client.SetTaskCtl(client.RetryCount, 5).Send("group1", "retry", 123, 44)
	result, _ := client.GetResult(tId, 3*time.Second, 300*time.Millisecond)
	fmt.Println("retry times =", result.RetryCount)

	// do not retry
	tId,_=client.SetTaskCtl(client.RetryCount, 0).Send("group1", "retry", 123, 44)
	result, _ = client.GetResult(tId, 3*time.Second, 300*time.Millisecond)
	fmt.Println("retry times =", result.RetryCount)

}

func delay() {
	// RunAfter
	tId, _ := client.SetTaskCtl(client.RunAfter, 1*time.Second).Send("group2", "add_sub", 123, 44)
	result, _ := client.GetResult(tId, 3*time.Second, 300*time.Millisecond)
	var sum, sub int
	result.Gets(&sum, &sub)
	fmt.Println("add_sub(123,44) =", sum, sub)

	// RunAt
	runTime := time.Now().Add(1 * time.Second)
	tId, _ = client.SetTaskCtl(client.RunAt, runTime).Send("group2", "add_sub", 123, 44)
	result, _ = client.GetResult(tId, 3*time.Second, 300*time.Millisecond)
	result.Gets(&sum, &sub)
	fmt.Println("add_sub(123,44) =", sum, sub)
}
func main() {
	// clientPoolSize: Maximum number of idle connections in client pool.
	//                 If clientPoolSize<=0, clientPoolSize=10
	//
	broker := ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 5)
	// poolSize: Maximum number of idle connections in the pool. If poolSize<=0 use default value
	//           default value is 10 at client
	//           ---------------
	//           对于client端，如果poolSize<=0，poolSize会设为10
	backend := ytask.Backend.NewRedisBackend("127.0.0.1", "6379", "", 0, 5)

	ser := ytask.Server.NewServer(
		ytask.Config.Broker(&broker),
		ytask.Config.Backend(&backend),
		ytask.Config.Debug(true),
		ytask.Config.StatusExpires(60*5),
		ytask.Config.ResultExpires(60*5),
	)

	client = ser.GetClient()

	fmt.Println("Send and get result\n---")
	sendAndGet()
	fmt.Println("\nRetry\n---")
	retry()
	fmt.Println("\nDelay\n---")
	delay()

}
