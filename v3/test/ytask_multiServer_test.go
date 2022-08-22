package test

//
// cd v3
// go test -v -count=1 test/*
//

import (
	"context"
	"github.com/gojuukaze/YTask/v3/backends"
	"github.com/gojuukaze/YTask/v3/brokers"
	"github.com/gojuukaze/YTask/v3/config"
	"github.com/gojuukaze/YTask/v3/log"
	"github.com/gojuukaze/YTask/v3/server"
	"github.com/gojuukaze/YTask/v3/yerrors"
	"io/ioutil"
	"testing"
	"time"
)

func multiWorker1() int {
	time.Sleep(2 * time.Second)
	return 123
}

func multiWorker2() int {
	return 123
}

func TestMulti(t *testing.T) {
	b := brokers.NewLocalBroker()
	b2 := backends.NewLocalBackend()

	ser := server.NewServer(
		config.NewConfig(
			config.Broker(&b),
			config.Backend(&b2),
			config.Debug(true),
		),
	)
	log.YTaskLog.Out = ioutil.Discard

	ser.Add("test_group1", "multiWorker1", multiWorker1)
	ser.Add("test_group2", "multiWorker2", multiWorker2)

	ser.Run("test_group1", 1)
	ser.Run("test_group2", 1)

	client := ser.GetClient()

	testMulti1(t, client)
	testMulti2(t, client)

	ser.Shutdown(context.TODO())

}

func testMulti1(t *testing.T, client server.Client) {

	client.Send("test_group1", "multiWorker1")
	id, _ := client.Send("test_group1", "multiWorker1")

	// 连续发两次，因为并发是1，此时第二次应该还没运行
	_, err := client.GetStatus(id, 1*time.Second, 300*time.Millisecond)
	if !yerrors.IsEqual(err, yerrors.ErrTypeTimeOut) {
		t.Fatal("err!=yerrors.ErrTypeTimeOut")
	}

	_, err = client.GetResult(id, 4*time.Second, 300*time.Millisecond)
	if err != nil {
		t.Fatal("err!=nil")
	}

}

func testMulti2(t *testing.T, client server.Client) {

	client.Send("test_group1", "multiWorker1")
	id, _ := client.Send("test_group2", "multiWorker2")

	result, err := client.GetResult(id, 1*time.Second, 300*time.Millisecond)
	if err != nil {
		t.Fatal("err!=nil")
	}
	r, _ := result.GetInt64(0)
	if r != 123 {
		t.Fatal("r!=123")

	}

}
