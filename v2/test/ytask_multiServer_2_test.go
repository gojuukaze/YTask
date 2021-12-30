package test

//
// cd v2
// go test -v -count=1 test/*
//

import (
	"context"
	"github.com/gojuukaze/YTask/v2/backends"
	"github.com/gojuukaze/YTask/v2/brokers"
	"github.com/gojuukaze/YTask/v2/config"
	"github.com/gojuukaze/YTask/v2/log"
	"github.com/gojuukaze/YTask/v2/server"
	"github.com/gojuukaze/YTask/v2/yerrors"
	"io/ioutil"
	"testing"
	"time"
)

func delayWorker1() int {
	return 123
}

func TestMultit2(t *testing.T) {
	b := brokers.NewRedisBroker("127.0.0.1", "6379", "", 0, 0)
	b2 := backends.NewRedisBackend("127.0.0.1", "6379", "", 0, 0)

	ser := server.NewServer(
		config.NewConfig(
			config.Broker(&b),
			config.Backend(&b2),
			config.Debug(true),
			config.EnableDelayServer(true),
		),
	)
	log.YTaskLog.Out = ioutil.Discard

	ser.Add("TestMulti2Group", "delayWorker1", delayWorker1)

	ser.Run("TestMulti2Group", 2)

	client := ser.GetClient()

	testMulti2_1(t, client)
	testMulti2_2(t, client)
	testMulti2_3(t, client)

	ser.Shutdown(context.TODO())

}

func testMulti2_1(t *testing.T, client server.Client) {
	// 测试两种任务能正常执行

	id, _ := client.Send("TestMulti2Group", "delayWorker1")
	id2, _ := client.SetTaskCtl(client.RunAfter, 100*time.Millisecond).Send("TestMulti2Group", "delayWorker1")

	_, err := client.GetStatus(id, 1*time.Second, 300*time.Millisecond)
	if err != nil {
		t.Fatal("err=", err)
	}

	_, err = client.GetResult(id2, 1*time.Second, 300*time.Millisecond)
	if err != nil {
		t.Fatal("err=", err)
	}

}

func testMulti2_2(t *testing.T, client server.Client) {
	// 测试延时任务是否延时执行

	client.Send("TestMulti2Group", "delayWorker1")

	id2, _ := client.SetTaskCtl(client.RunAfter, 2*time.Second).Send("TestMulti2Group", "delayWorker1")

	_, err := client.GetResult(id2, 1*time.Second, 300*time.Millisecond)
	if !yerrors.IsEqual(err, yerrors.ErrTypeTimeOut) {
		t.Fatal("err!=yerrors.ErrTypeTimeOut")
	}
	time.Sleep(1 * time.Second)

	_, err = client.GetResult(id2, 1*time.Second, 300*time.Millisecond)
	if err != nil {
		t.Fatal("err=", err)
	}

}

func testMulti2_3(t *testing.T, client server.Client) {
	// 测试多个延时任务的执行顺序

	id2, _ := client.SetTaskCtl(client.RunAfter, 2*time.Second).Send("TestMulti2Group", "delayWorker1")

	id, _ := client.SetTaskCtl(client.RunAfter, 100*time.Millisecond).Send("TestMulti2Group", "delayWorker1")

	_, err := client.GetResult(id, 500*time.Millisecond, 100*time.Millisecond)
	if err != nil {
		t.Fatal("err=", err)
	}
	_, err = client.GetResult(id2, 1*time.Second, 300*time.Millisecond)
	if !yerrors.IsEqual(err, yerrors.ErrTypeTimeOut) {
		t.Fatal("err!=yerrors.ErrTypeTimeOut ",err)
	}

	time.Sleep(1 * time.Second)

	_, err = client.GetResult(id2, 1*time.Second, 300*time.Millisecond)
	if err != nil {
		t.Fatal("err=", err)
	}

}
