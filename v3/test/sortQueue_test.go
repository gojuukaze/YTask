package test

import (
	"fmt"
	"github.com/gojuukaze/YTask/v3/message"
	"github.com/gojuukaze/YTask/v3/server"
	"math/rand"
	"testing"
	"time"
)

func newMsg(t time.Time) message.Message {
	c := message.NewMsgArgs()
	c.RunTime = t
	return message.NewMessage(c)
}

func TestSortQ(t *testing.T) {
	testSortQ(t)

}
func testSortQ(t *testing.T) {
	q := server.NewSortQueue(20)

	for i := 0; i < 20; i++ {
		t := time.Now().Add(time.Duration(rand.Intn(600)) * time.Second)
		q.Insert(newMsg(t))
	}
	// 排序是否正确
	for i := 0; i < 20; i += 2 {
		if !q.Queue[i].RunTimeBefore(q.Queue[i+1].MsgArgs.GetRunTime()) && !q.Queue[i].RunTimeEqual(q.Queue[i+1].MsgArgs.GetRunTime()) {
			fmt.Println(q.Queue)
			t.Fatal("排序错误")
		}
	}

	// 插入一个较前的任务时，最后一个任务会出队
	temp := q.Queue[19]
	pop := q.Insert(newMsg(q.Queue[3].MsgArgs.GetRunTime()))
	if !temp.RunTimeEqual(pop.MsgArgs.GetRunTime()) {
		t.Fatal("temp!=pop", temp, pop)

	}

	// 插入的任务比最后一个运行时间还旧时，任务不会插入
	temp = newMsg(q.Queue[19].MsgArgs.GetRunTime().Add(2000))
	pop = q.Insert(temp)
	if !temp.RunTimeEqual(pop.MsgArgs.GetRunTime()) {
		t.Fatal("temp!=pop", temp, pop)

	}

}
