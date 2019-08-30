package main

import (
	"fmt"
	"github.com/gojuukaze/YTask/v2"
	"github.com/gojuukaze/YTask/v2/server"
	"time"
)

var client server.Client

func sendAndGet() {
	// task add
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

	// task add_sub
	taskId, err = client.Send("group1", "add_sub", 123, 44)
	_ = err
	result, err = client.GetResult(taskId, 2*time.Second, 300*time.Millisecond)
	_ = err

	if result.IsSuccess() {
		sum, _ := result.GetInt64(0)
		sub, _ := result.GetInt64(1)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("add_sub(123,44) =", int(sum), int(sub))
	} else {
		fmt.Println("result failure")
	}
}

func retry() {

	// set retry count, default is 3
	tId, _ := client.SetTaskCtl(client.RetryCount, 5).Send("group1", "retry", 123, 44)
	result, _ := client.GetResult(tId, 3*time.Second, 300*time.Millisecond)
	fmt.Println("retry times =", result.RetryCount)

	// do not retry
	client.SetTaskCtl(client.RetryCount, 5).Send("group1", "retry", 123, 44)


}
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

	sendAndGet()

	retry()

}
