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

func workflow1(a int, b int) int {
	return a + b
}

func workflow2(a int) int {

	return a * a
}

func TestWorkflow(t *testing.T) {
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

	ser.Add("test_g", "workflow1", workflow1)
	ser.Add("test_g", "workflow2", workflow2)

	client := ser.GetClient()
	ser.Run("test_g", 2, true)

	testWorkflow1(client, t)
	testWorkflow2(client, t)
	testWorkflow3(client, t)
	ser.Shutdown(context.TODO())

}

func testWorkflow1(client server.Client, t *testing.T) {
	id, _ := client.Workflow().
		Send("test_g", "workflow1", 1, 2).
		Send("test_g", "workflow2").
		Done()

	result, err := client.GetResult(id, time.Second*2, time.Millisecond*300)

	a, _ := result.GetInt64(0)
	if a != 9 {
		t.Fatalf("a is %d , !=3 ; err=%s", a, err)
	}
}

// 测试延时任务执行
func testWorkflow2(client server.Client, t *testing.T) {
	id, _ := client.Workflow().
		SetTaskCtl(client.RunAfter, 3*time.Second).
		Send("test_g", "workflow1", 1, 2).
		SetTaskCtl(client.RunAfter, 3*time.Second).
		Send("test_g", "workflow2").
		Done()

	result, _ := client.GetResult(id, time.Second*2, time.Millisecond*300)
	if result.IsFinish() {
		t.Fatalf("result.IsFinish()")
	}
	time.Sleep(1 * time.Second)
	result, _ = client.GetResult2(id, time.Second*2, time.Millisecond*300)
	if result.Workflow[0][1] != "success" || result.Workflow[1][1] != "waiting" {
		t.Fatalf("WorkflowStatus error %v", result.Workflow)
	}
	result, _ = client.GetResult(id, time.Second*5, time.Millisecond*300)

	a, _ := result.GetInt64(0)
	if a != 9 {
		t.Fatalf("a is %d , !=3", a)
	}
}

// 测试任务过期
func testWorkflow3(client server.Client, t *testing.T) {
	id, _ := client.Workflow().
		SetTaskCtl(client.RunAfter, 2*time.Second).
		SetTaskCtl(client.ExpireTime, time.Now()).
		Send("test_g", "workflow1", 1, 2).
		Send("test_g", "workflow2").
		Done()

	result, _ := client.GetResult(id, time.Second*4, time.Millisecond*300)
	if result.Status != message.ResultStatus.Expired ||
		result.Workflow[0][1] != message.WorkflowStatus.Expired {
		t.Fatalf("result.Status error %v", result)
	}

	id, _ = client.Workflow().
		Send("test_g", "workflow1", 1, 2).
		SetTaskCtl(client.RunAfter, 2*time.Second).
		SetTaskCtl(client.ExpireTime, time.Now()).
		Send("test_g", "workflow2").
		Done()

	result, _ = client.GetResult(id, time.Second*5, time.Millisecond*300)
	if result.Status != message.ResultStatus.Expired ||
		result.Workflow[1][1] != message.WorkflowStatus.Expired {
		t.Fatalf("result.Status error %v", result)
	}
}
