package main

import (
	"flag"
	"fmt"
	"github.com/gojuukaze/YTask/v2"
)

func main() {
	flag.Parse()

	b := ytask.Broker.NewRedisBroker("127.0.0.1", "6379", "", 0, 3)
	ser := ytask.Server.NewServer(
		ytask.Config.Broker(b),
		ytask.Config.Debug(true),
	)

	err:=ser.Send("g1", "add", 123, 33)
	fmt.Println(err)

}
