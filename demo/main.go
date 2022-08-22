package main

import (
	"context"
	"fmt"
	"github.com/gojuukaze/YTask/v3/backends"
	"github.com/gojuukaze/YTask/v3/brokers"
	"github.com/gojuukaze/YTask/v3/config"
	"github.com/gojuukaze/YTask/v3/log"
	"github.com/gojuukaze/YTask/v3/server"
	"time"
)

func workflow1(a int, b int) int {
	return a + b
}

func workflow2(a int) int {
	return a * a
}
func main() {
	fmt.Println("hello world")
	b := brokers.NewLocalBroker()
	b2 := backends.NewLocalBackend()
	l := log.NewYTaskLogger(log.YTaskLog)
	l.Info("hello world")

	ser := server.NewServer(
		config.NewConfig(
			config.Broker(&b),
			config.Backend(&b2),
			config.Debug(true),
		),
	)
	//log.YTaskLog.Out = ioutil.Discard

	ser.Add("test_g", "workflow1", workflow1)
	ser.Add("test_g", "workflow2", workflow2)
	ser.Run("test_g", 2, true)
	//testWorkflow1(ser, t)
	time.Sleep(time.Second * 5)
	ser.Shutdown(context.TODO())
}
