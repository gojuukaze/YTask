package main

import (
	"flag"
	"fmt"
	"github.com/gojuukaze/YTask/v2"
	"time"
)

func main() {
	flag.Parse()

	b := ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 5)
	b2 := ytask.Backend.NewRedisBackend("127.0.0.1", "6379", "", 0, 5)

	ser := ytask.Server.NewServer(
		ytask.Config.Broker(&b),
		ytask.Config.Backend(&b2),
		ytask.Config.Debug(true),
	)

	client := ser.GetClient()
	id, err := client.Send("g1", "add", 123, 44)
	fmt.Println(err)
	result, err := client.GetResult(id, 5*time.Second, 300*time.Millisecond)
	fmt.Println(err)

	if err == nil {
		if result.IsSuccess() {
			sum, _ := result.GetInt64(0)

			fmt.Println("result=", sum)
		} else {
			fmt.Println("result failure")
		}
	} else {
		fmt.Println(err)
	}

}
