package test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/vua/YTask/v2/controller"
	"github.com/vua/YTask/v2/message"
	"github.com/vua/YTask/v2/server"
)

func newMsg(t time.Time) message.Message {
	c := controller.NewTaskCtl()
	c.SetRunTime(t)
	return message.NewMessage(c)
}

func TestSortQ(t *testing.T) {

	for i := 0; i < 20; i++ {
		testSortQ(t)
	}

}
func testSortQ(t *testing.T) {
	q := server.SortQueue{}

	for i := 0; i < 20; i++ {
		t := time.Now().Add(time.Duration(rand.Intn(600)) * time.Second)
		q.Insert(newMsg(t))
	}

	for i := 0; i < 20; i += 2 {
		if !q.Queue[i].RunTimeBefore(q.Queue[i+1].TaskCtl.GetRunTime()) && !q.Queue[i].RunTimeEqual(q.Queue[i+1].TaskCtl.GetRunTime()) {
			fmt.Println(q.Queue)
			t.Fatal("排序错误")
		}
	}

	temp := q.Queue[19]
	pop := q.Insert(newMsg(q.Queue[3].TaskCtl.GetRunTime()))
	if !temp.RunTimeEqual(pop.TaskCtl.GetRunTime()) {
		t.Fatal("temp!=pop", temp, pop)

	}

	temp = newMsg(q.Queue[19].TaskCtl.GetRunTime().Add(2000))
	pop = q.Insert(temp)
	if !temp.RunTimeEqual(pop.TaskCtl.GetRunTime()) {
		t.Fatal("temp!=pop", temp, pop)

	}

}
