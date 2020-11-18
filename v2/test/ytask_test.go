package test

//
// cd v2
// go test -v -count=1 test/*
//

import (
	"context"
	"errors"
	"fmt"
	"github.com/gojuukaze/YTask/v2/backends"
	"github.com/gojuukaze/YTask/v2/brokers"
	"github.com/gojuukaze/YTask/v2/config"
	"github.com/gojuukaze/YTask/v2/controller"
	"github.com/gojuukaze/YTask/v2/log"
	"github.com/gojuukaze/YTask/v2/server"
	"io/ioutil"
	"os"
	"testing"
	"time"
)
func init(){
	fmt.Println("ROCKETMQ_GO_LOG_LEVEL setting")
	os.Setenv("ROCKETMQ_GO_LOG_LEVEL","info")
}
type User struct {
	Id   int
	Name string
}

func worker1() {
}

func worker2(a int, b float32, c uint64, d bool) (float32, uint64, bool) {
	return float32(a) + b, c, d
}

func worker3(user User, ids []int, names []string) []User {
	var r = make([]User, 0)
	r = append(r, user)
	for i := range ids {
		r = append(r, User{
			Id:   ids[i],
			Name: names[i],
		})
	}
	return r
}
func workerTestRetry1() {
	panic("test retry")
}

func workerTestRetry2(ctl *controller.TaskCtl, a int) int {
	if ctl.RetryCount == 3 {
		panic("test retry")
	} else if ctl.RetryCount == 2 {
		ctl.Retry(errors.New("test retry 2"))
		return 0
	}
	return a + ctl.RetryCount
}
func TestYTask1(t *testing.T) {
	b := brokers.NewRocketMqBroker("127.0.0.1", "9876")
	//b := brokers.NewRedisBroker("127.0.0.1", "6379", "", 0, 0)

	b2 := backends.NewRedisBackend("127.0.0.1", "6379", "", 0, 0)

	ser := server.NewServer(
		config.NewConfig(
			config.Broker(&b),
			config.Backend(&b2),
			config.Debug(true),
		),
	)
	log.YTaskLog.Out = ioutil.Discard

	ser2 := ser
	ser.Add("test_g", "worker1", worker1)
	ser.Add("test_g", "worker2", worker2)
	ser.Add("test_g", "worker3", worker3)

	ser.Add("test_g", "workerTestRetry1", workerTestRetry1)
	ser.Add("test_g", "workerTestRetry2", workerTestRetry2)

	ser.Run("test_g", 3)
	testWorker1(ser2, t)
	testWorker2(ser2, t)
	testWorker3(ser2, t)

	testRetry1(ser2, t)
	testRetry2(ser2, t)

	ser.Shutdown(context.TODO())

}

func testWorker1(ser server.Server, t *testing.T) {
	client := ser.GetClient()

	id, err := client.Send("test_g", "worker1")

	if err != nil {
		t.Fatal(err)
	}
	result, _ := client.GetResult(id, 200*time.Second, 300*time.Millisecond)

	if !result.IsSuccess() {
		t.Fatal("result is not success")
	}
}

func testWorker2(ser server.Server, t *testing.T) {
	client := ser.GetClient()
	var (
		a int     = 12
		b float32 = 22.1
		c uint64  = 18446744073709551610
		d         = true
	)
	id, err := client.Send("test_g", "worker2", a, b, c, d)
	if err != nil {
		t.Fatal(err)
	}

	result, _ := client.GetResult(id, 2*time.Second, 300*time.Millisecond)

	if !result.IsSuccess() {
		t.Fatal("result is not success")
	}

	r1, _ := result.GetFloat64(0)
	if float32(a)+b != float32(r1) {
		t.Fatalf("%v != %v", float32(a)+b, float32(r1))
	}

	r2, _ := result.GetUint64(1)
	if c != r2 {
		t.Fatalf("%v != %v", c, r2)
	}

	r3, _ := result.GetBool(2)
	if d != r3 {
		t.Fatalf("%v != %v", c, r3)
	}
}

func testWorker3(ser server.Server, t *testing.T) {
	client := ser.GetClient()

	id, err := client.Send("test_g", "worker3",
		User{
			Id:   1,
			Name: "a",
		},
		[]int{233, 44},
		[]string{"bb", "cc"})
	if err != nil {
		t.Fatal(err)
	}
	result, _ := client.GetResult(id, 20*time.Second, 300*time.Millisecond)

	if !result.IsSuccess() {
		t.Fatal("result is not success")
	}
	var base = []User{{1, "a"}, {233, "bb"}, {44, "cc"},}
	var r []User
	err = result.Get(0, &r)

	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(base) != fmt.Sprint(r) {
		t.Fatalf("%v !=%v", base, r)
	}
}

func testRetry1(ser server.Server, t *testing.T) {
	client := ser.GetClient()

	id, err := client.SetTaskCtl(client.RetryCount, 5).Send("test_g", "workerTestRetry1")
	if err != nil {
		t.Fatal(err)
	}
	result, _ := client.GetResult(id, 2*time.Second, 300*time.Millisecond)

	if result.IsSuccess() {
		t.Fatal("result is success")
	}

	if result.RetryCount != 5 {
		t.Fatal("result.RetryCount!=5")

	}

}

func testRetry2(ser server.Server, t *testing.T) {
	client := ser.GetClient()

	id, err := client.Send("test_g", "workerTestRetry2", 6)
	if err != nil {
		t.Fatal(err)
	}
	result, _ := client.GetResult(id, 2*time.Second, 300*time.Millisecond)

	if !result.IsSuccess() {
		t.Fatal("result is not success")
	}
	r1, _ := result.GetInt64(0)
	if int(r1) != 6+1 {
		t.Fatal("int(r1)!=6+1")

	}

}
