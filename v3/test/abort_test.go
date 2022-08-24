package test

import (
	"context"
	"github.com/gojuukaze/YTask/v3/backends"
	"github.com/gojuukaze/YTask/v3/brokers"
	"github.com/gojuukaze/YTask/v3/config"
	"github.com/gojuukaze/YTask/v3/log"
	"github.com/gojuukaze/YTask/v3/server"
	"io/ioutil"
	"testing"
	"time"
)

func abortW1(a int, b int) int {
	return a + b
}

func abortW2(a int) int {

	return a * a
}

func TestAbort(t *testing.T) {
	b := brokers.NewLocalBroker()
	b2 := backends.NewLocalBackend()
	log.YTaskLog.Out = ioutil.Discard

	ser := server.NewServer(
		config.NewConfig(
			config.Broker(&b),
			config.Backend(&b2),
			config.ResultExpires(100),
			config.Debug(true),
		),
	)

	ser.Add("test_g", "abortW1", abortW1)
	ser.Add("test_g", "abortW2", abortW2)

	client := ser.GetClient()
	ser.Run("test_g", 2, true)

	testAbort1(client, t)
	//testWorkflow2(client, t)
	//testWorkflow3(client, t)
	ser.Shutdown(context.TODO())

}

func testAbort1(client server.Client, t *testing.T) {
	id, _ := client.SetTaskCtl(client.RunAfter, 1*time.Second).
		Send("test_g", "abortW1", 1, 2)
	client.AbortTask(id, 10)
	result, _ := client.GetResult(id, time.Second*2, time.Millisecond*300)

	if !result.IsFailure() {
		t.Fatalf("!result.IsFailure()")
	}
}
