package test

import (
	"fmt"
	"github.com/gojuukaze/YTask/drives/rabbitmq/v3"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func putMsg() {

}
func TestRabbitmqPool(t *testing.T) {
	c := rabbitmq.NewRabbitMqClient("127.0.0.1", "5672", "guest", "guest", "", 2)
	channel0, _ := c.GetChannel()
	c.QueueDeclare("test_amqp", channel0)
	c.PutChannel(channel0, false)

	channel1, _ := c.GetChannel()

	// channel 0 与 1 应该是同一个
	if channel1 != channel0 {
		t.Fatal("channel1 != channel")
	}
	channel2, err := c.GetChannel()
	if err != nil {
		t.Fatal(err)
	}

	// channel 2 与 1 不是同一个
	if channel1 == channel2 {
		t.Fatal("channel1 == channel")
	}
	// 临时修改超时时间
	rabbitmq.GetChanTimeout = 5 * time.Second
	defer func() {
		rabbitmq.GetChanTimeout = 60 * time.Second
	}()
	// 这里应该是超时
	_, err = c.GetChannel()
	if err != rabbitmq.ErrNoIdleChannel {
		t.Fatal(err)
	}
}

func TestRabbitmqPoolBadChan(t *testing.T) {
	/*
		测试Rabbitmq重启后能否创建链接
	*/
	c := rabbitmq.NewRabbitMqClient("127.0.0.1", "5672", "guest", "guest", "", 2)
	channel0, _ := c.GetChannel()

	// test中无法通过Scanf()获取输入, 因此通过信号传递信息
	fmt.Printf("重启Rabbitmq后通过 \" kill -CONT %d \" 继续运行测试\n", os.Getpid())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGCONT)
	<-quit

	c.PutChannel(channel0, false)

	channel1, _ := c.GetChannel()

	// 因为链接断了，所以channel0和channel1应该是不同的
	if channel0 == channel1 {
		t.Fatal("channel0 == channel1")
	}
	// 这里应该是成功的
	err := c.QueueDeclare("test_amqp22", channel1)
	if err != nil {
		t.Fatal(err)
	}

}
