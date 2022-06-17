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

var workerTestCallbackChan chan workerTestCallbackResult

type workerTestCallbackResult struct {
	a      int
	s      string
	result *message.Result
}

func workerTestCallback(a int, s string) int {

	a = a * 6
	return 233
}

func workerTestCallback2(a int, s string) int {
	panic("err")
	return 233
}

func callbackTestCallback(a int, s string, result *message.Result) {

	workerTestCallbackChan <- workerTestCallbackResult{a, s, result}

}

func callbackTestCallback2(a int, s string, result *message.Result) {
	panic("err")
}

func TestCallback(t *testing.T) {
	workerTestCallbackChan = make(chan workerTestCallbackResult, 2)
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

	ser.Add("test_callback", "workerTestCallback", workerTestCallback, callbackTestCallback)
	ser.Add("test_callback", "workerTestCallback2", workerTestCallback2, callbackTestCallback)
	ser.Add("test_callback", "workerTestCallback3", workerTestCallback, callbackTestCallback2)

	ser.Run("test_callback", 1)
	client := ser.GetClient()
	testCallback_1(t, client)
	testCallback_2(t, client)
	testCallback_3(t, client)

	ser.Shutdown(context.TODO())

}

func testCallback_1(t *testing.T, client server.Client) {
	_, _ = client.Send("test_callback", "workerTestCallback", 1, "s")

	r := <-workerTestCallbackChan
	// worker函数中修改参数不会改变callback中的参数
	if r.a != 1 || r.s != "s" {
		t.Fatal("r.a!=1 || r.s!=\"s\"")
	}
	// 验证callback中的result
	if r.result.IsSuccess() != true {
		t.Fatal("IsSuccess!=true")

	}
	r0, _ := r.result.GetInt64(0)
	if r0 != 233 {
		t.Fatal("r0!=233")
	}
}

func testCallback_2(t *testing.T, client server.Client) {
	// worker错误，callback中能获取错误
	_, _ = client.Send("test_callback", "workerTestCallback2", 1, "s")

	r := <-workerTestCallbackChan

	if r.result.IsSuccess() != false {
		t.Fatal("IsSuccess!=false")
	}

}

func testCallback_3(t *testing.T, client server.Client) {
	// callback错误不会影响结果获取
	id, _ := client.Send("test_callback", "workerTestCallback3", 1, "s")
	r, _ := client.GetResult(id, 2*time.Second, 300*time.Millisecond)

	if r.IsSuccess() != true {
		t.Fatal("IsSuccess!=true")

	}
	r0, _ := r.GetInt64(0)
	if r0 != 233 {
		t.Fatal("r0!=233")
	}

}
