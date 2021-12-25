package test

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/vua/YTask/v2/backends"
	"github.com/vua/YTask/v2/brokers"
	"github.com/vua/YTask/v2/config"
	"github.com/vua/YTask/v2/log"
	"github.com/vua/YTask/v2/message"
	"github.com/vua/YTask/v2/server"
)

func workerTestStatusExpires() int {
	time.Sleep(2 * time.Second)
	return 233
}

func workerTestResultExpires() int {
	time.Sleep(2 * time.Second)

	return 233
}

func TestStatusExpires(t *testing.T) {
	b := brokers.NewRedisBroker("127.0.0.1", "6379", "", 0, 0)
	b2 := backends.NewRedisBackend("127.0.0.1", "6379", "", 0, 0)

	ser := server.NewServer(
		config.NewConfig(
			config.Broker(&b),
			config.Backend(&b2),
			config.Debug(false),
			config.StatusExpires(0),
		),
	)
	log.YTaskLog.Out = ioutil.Discard

	ser.Add("test_g2", "workerTestStatusExpires", workerTestStatusExpires)
	ser.Run("test_g2", 1)
	client := ser.GetClient()
	id, err := client.Send("test_g2", "workerTestStatusExpires")
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.GetStatus(id, 1*time.Second, 300*time.Millisecond)
	if err == nil {
		t.Fatal("err==nill")
	}

	result, _ := client.GetResult(id, 3*time.Second, 300*time.Millisecond)
	if !result.IsSuccess() {
		t.Fatal("!result.IsSuccess()")

	}
	a, _ := result.GetInt64(0)

	if int(a) != 233 {
		t.Fatal("int(a)!=233")
	}
	ser.Shutdown(context.TODO())

}

func TestResultExpires(t *testing.T) {
	b := brokers.NewRedisBroker("127.0.0.1", "6379", "", 0, 0)
	b2 := backends.NewRedisBackend("127.0.0.1", "6379", "", 0, 0)

	ser := server.NewServer(
		config.NewConfig(
			config.Broker(&b),
			config.Backend(&b2),
			config.Debug(false),
			config.ResultExpires(0),
		),
	)
	log.YTaskLog.Out = ioutil.Discard

	ser.Add("test_g2", "workerTestResultExpires", workerTestResultExpires)
	ser.Run("test_g2", 1)

	client := ser.GetClient()
	id, err := client.Send("test_g2", "workerTestResultExpires")
	if err != nil {
		t.Fatal(err)
	}
	status, err := client.GetStatus(id, 1*time.Second, 100*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if status != message.ResultStatus.FirstRunning {
		t.Fatal("r1.Status!=message.ResultStatus.FirstRunning", status)

	}

	_, err = client.GetResult(id, 3*time.Second, 300*time.Millisecond)
	if err == nil {
		t.Fatal("err==nil")

	}
	ser.Shutdown(context.TODO())
}
