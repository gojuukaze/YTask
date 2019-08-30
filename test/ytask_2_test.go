package test

import (
	"context"
	"github.com/gojuukaze/YTask/v2/backends"
	"github.com/gojuukaze/YTask/v2/brokers"
	"github.com/gojuukaze/YTask/v2/config"
	"github.com/gojuukaze/YTask/v2/log"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/server"
	"io/ioutil"
	"testing"
	"time"
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
	ser.Run("test_g2",1)
	client := ser.GetClient()
	id, err := client.Send("test_g2", "workerTestStatusExpires")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1*time.Second)
	_, err = ser.GetResult(id)
	if err == nil {
		t.Fatal("err==nill")
	}

	result, _ := client.GetResult(id, 3*time.Second, 300*time.Millisecond)
	if !result.IsSuccess(){
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
	ser.Run("test_g2",1)

	client := ser.GetClient()
	id, err := client.Send("test_g2", "workerTestResultExpires")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1*time.Second)
	r1, err := ser.GetResult(id)
	if err != nil {
		t.Fatal(err)
	}
	if r1.Status!=message.ResultStatus.FirstRunning{
		t.Fatal("r1.Status!=message.ResultStatus.FirstRunning")

	}

	_, err = client.GetResult(id, 2*time.Second, 300*time.Millisecond)
	if err==nil{
		t.Fatal("err==nil")

	}
	ser.Shutdown(context.TODO())
}
