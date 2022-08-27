package test

import (
	"context"
	"github.com/gojuukaze/YTask/v3/backends"
	"github.com/gojuukaze/YTask/v3/brokers"
	"github.com/gojuukaze/YTask/v3/config"
	"github.com/gojuukaze/YTask/v3/log"
	"github.com/gojuukaze/YTask/v3/message"
	"github.com/gojuukaze/YTask/v3/server"
	"io/ioutil"
	"testing"
	"time"
)

func abortW1(a int, b int) int {
	return a + b
}

func abortW2(ctl *server.TaskCtl, a int) int {
	time.Sleep(3 * time.Second)

	f, _ := ctl.IsAbort()
	if f {
		ctl.Abort("手动中止")
		return 0
	}
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
	testAbort2(client, t)
	ser.Shutdown(context.TODO())

}

func testAbort1(client server.Client, t *testing.T) {
	id, _ := client.SetTaskCtl(client.RunAfter, 1*time.Second).
		Send("test_g", "abortW1", 1, 2)
	client.AbortTask(id, 10)
	result, _ := client.GetResult(id, time.Second*2, time.Millisecond*300)

	if !result.IsFailure() || result.Status != message.ResultStatus.Abort {
		t.Fatalf("!result.IsFailure()")
	}
}

func testAbort2(client server.Client, t *testing.T) {
	id, _ := client.Send("test_g", "abortW2", 1, 2)
	time.Sleep(1 * time.Second)
	client.AbortTask(id, 10)

	time.Sleep(1 * time.Second)
	result, _ := client.GetResult(id, time.Second*3, time.Millisecond*100)

	if !result.IsFailure() || result.Status != message.ResultStatus.Abort {
		t.Fatalf("!result.IsFailure()")
	}
}
