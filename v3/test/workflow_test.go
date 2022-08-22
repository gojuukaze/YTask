package test

import (
	"context"
	"fmt"
	"github.com/gojuukaze/YTask/v3/backends"
	"github.com/gojuukaze/YTask/v3/brokers"
	"github.com/gojuukaze/YTask/v3/config"
	"github.com/gojuukaze/YTask/v3/server"
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
	testWorkflow1(ser, t)
	ser.Shutdown(context.TODO())

}

func testWorkflow1(ser server.Server, t *testing.T) {
	client := ser.GetClient()
	id, _ := client.Workflow().
		Send("test_g", "workflow1", 1, 2).
		Send("test_g", "workflow2").
		Done()

	time.Sleep(time.Second * 1)
	result, _ := client.GetResult(id, time.Second*2, time.Millisecond*300)
	fmt.Println(result)
	a, _ := result.GetInt64(0)
	t.Logf("a=%d", a)
	if a != 9 {
		t.Fatalf("a is %d , !=3", a)
	}
}
