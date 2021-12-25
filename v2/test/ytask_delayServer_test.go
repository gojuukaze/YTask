package test

//
// cd v2
// go test -v -count=1 test/*
//

import (
	"context"
	"io/ioutil"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/vua/YTask/v2/brokers"
	"github.com/vua/YTask/v2/config"
	"github.com/vua/YTask/v2/log"
	"github.com/vua/YTask/v2/message"
	"github.com/vua/YTask/v2/server"
)

func TestDelayServerSend(t *testing.T) {
	// 测试client send
	b := brokers.NewRedisBroker("127.0.0.1", "6379", "", 0, 0)
	client := server.NewClient(config.NewConfig(
		config.Broker(&b),
		config.Debug(true),
	))
	runTime, _ := time.Parse("2006-01-02 15:04:05", "2018-04-23 12:24:51")

	client.SetTaskCtl(client.RunAt, runTime).Send("TestDelayServer_Send", "1")
	runTime2 := time.Now().Add(5 * time.Second)
	client.SetTaskCtl(client.RunAfter, 5*time.Second).Send("TestDelayServer_Send", "2")

	b2 := b.Clone()
	b2.Activate()
	msg, _ := b2.Next("YTask:Query:Delay:TestDelayServer_Send")
	if msg.TaskCtl.GetRunTime() != runTime {
		t.Fatal(msg.TaskCtl.GetRunTime(), "!=", runTime)
	}

	msg, _ = b2.Next("YTask:Query:Delay:TestDelayServer_Send")
	if msg.TaskCtl.GetRunTime().Sub(runTime2).Milliseconds() > 300 {
		t.Fatal(msg.TaskCtl.GetRunTime().Sub(runTime2), ">300 ms")

	}

}

func TestDelayServer(t *testing.T) {
	// 测试处理顺序是否正确

	b := brokers.NewRedisBroker("127.0.0.1", "6379", "", 0, 0)
	ch := make(chan message.Message, 5)
	ds := server.NewDelayServer("testDelay", config.NewConfig(
		config.Broker(&b),
		config.Debug(true),
	), ch)
	client := server.NewClient(config.NewConfig(
		config.Broker(&b),
		config.Debug(true),
	))
	log.YTaskLog.Out = ioutil.Discard
	ds.Run()

	wg := sync.WaitGroup{}
	client.SetTaskCtl(client.RunAfter, 6*time.Second).Send("testDelay", "3")
	client.SetTaskCtl(client.RunAfter, 6*time.Second).Send("testDelay", "4")

	client.SetTaskCtl(client.RunAfter, 3*time.Second).Send("testDelay", "2")
	client.SetTaskCtl(client.RunAfter, 2*time.Second).Send("testDelay", "1")
	wg.Add(4)
	go func() {
		wg.Wait()
		close(ch)
	}()
	i := 1
	for msg := range ch {
		if msg.WorkerName != strconv.Itoa(i) {
			t.Fatal(msg.WorkerName, "!=", i)
		}
		wg.Done()
		i++
	}

	ds.Shutdown(context.TODO())

	//清空测试用的队列
	b2 := b.Clone()
	b2.Activate()
	for true {
		_, err := b2.Next("YTask:Query:Delay:testDelay")
		if err != nil {
			break
		}
	}

}

func TestDelayServer2(t *testing.T) {
	// 测试服务关闭后，在本地队列中的任务是否能插入到broker中
	b := brokers.NewRedisBroker("127.0.0.1", "6379", "", 0, 0)
	ch := make(chan message.Message, 5)
	ds := server.NewDelayServer("testDelay2", config.NewConfig(
		config.Broker(&b),
		config.Debug(true),
	), ch)
	client := server.NewClient(config.NewConfig(
		config.Broker(&b),
		config.Debug(true),
	))
	log.YTaskLog.Out = ioutil.Discard
	//log.YTaskLog.SetLevel(logrus.DebugLevel)
	ds.Run()
	/*
		放22个任务后，在放3个任务。
		2秒后：1，2依次出队；3在本地队列的头，若使用RedisBroker，在server关闭后，3应该在Redis队列的第一个

	*/
	for i := 0; i < 22; i++ {
		client.SetTaskCtl(client.RunAfter, 15*time.Second).Send("testDelay2", "d"+strconv.Itoa(i))
	}

	client.SetTaskCtl(client.RunAfter, 3*time.Second).Send("testDelay2", "3")
	client.SetTaskCtl(client.RunAfter, 2*time.Second).Send("testDelay2", "1")
	client.SetTaskCtl(client.RunAfter, 2*time.Second).Send("testDelay2", "2")

	time.Sleep(2 * time.Second)

	msg := <-ch

	if msg.WorkerName != "1" {
		t.Fatal(msg.WorkerName, "!=1")
	}
	msg = <-ch
	if msg.WorkerName != "2" {
		t.Fatal(msg.WorkerName, "!=2")
	}
	close(ch)
	ds.Shutdown(context.TODO())

	b2 := b.Clone()
	b2.Activate()
	msg, _ = b2.Next("YTask:Query:Delay:testDelay2")
	if msg.WorkerName != "3" {
		t.Fatal(msg.WorkerName, "!=3")
	}
	//清空测试用的队列
	for true {
		_, err := b2.Next("YTask:Query:Delay:testDelay2")
		if err != nil {
			break
		}
	}

}

func TestDelayServer3(t *testing.T) {
	// 测试服务关闭后，在readyChan中的任务是否能插入到broker中
	groupName := "testDelay3"

	b := brokers.NewRedisBroker("127.0.0.1", "6379", "", 0, 0)
	// 因为要模拟inlineServer处理任务的情况，这个chan不能有缓存
	ch := make(chan message.Message)

	ds := server.NewDelayServer(groupName, config.NewConfig(
		config.Broker(&b),
		config.Debug(true),
	), ch)
	client := server.NewClient(config.NewConfig(
		config.Broker(&b),
		config.Debug(true),
	))
	log.YTaskLog.Out = ioutil.Discard
	//log.YTaskLog.SetLevel(logrus.DebugLevel)
	ds.Run()
	/*
		放22个任务后，在放3个任务。
		2秒后：1，2在inlineServerMsgChan中，关闭服务后，1，2应该插入到broker

	*/
	for i := 0; i < 22; i++ {
		client.SetTaskCtl(client.RunAfter, 15*time.Second).Send(groupName, "d"+strconv.Itoa(i))
	}

	client.SetTaskCtl(client.RunAfter, 5*time.Second).Send(groupName, "3")
	client.SetTaskCtl(client.RunAfter, 2*time.Second).Send(groupName, "1")
	client.SetTaskCtl(client.RunAfter, 2*time.Second).Send(groupName, "2")

	time.Sleep(3 * time.Second)

	close(ch)

	ds.Shutdown(context.TODO())

	b2 := b.Clone()
	b2.Activate()
	var names = map[string]string{"1": "", "2": ""}

	for true {
		msg, err := b2.Next("YTask:Query:Delay:" + groupName)
		if err != nil {
			break
		}
		_, ok := names[msg.WorkerName]
		if ok {
			delete(names, msg.WorkerName)
		}
	}
	if len(names) != 0 {
		t.Fatal("not found ", names)

	}

}
