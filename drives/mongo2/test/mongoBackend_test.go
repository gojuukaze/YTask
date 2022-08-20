package test

import (
	"fmt"
	"github.com/gojuukaze/YTask/v3/drives/mongo2"
	"github.com/gojuukaze/YTask/v3/message"
	"github.com/gojuukaze/YTask/v3/yerrors"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestMongoBackend(t *testing.T) {
	b := mongo2.NewMongoBackend("127.0.0.1", "27017", "", "", "test", "test", 2)
	result := message.NewResult("xx123")
	result.FuncReturn = nil
	b.Activate()
	err := b.SetResult(result, 0)
	if err != nil {
		t.Fatal(err)
	}

	r2, err := b.GetResult(result.GetBackendKey())
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%v", r2) != fmt.Sprintf("%v", result) {
		t.Fatalf("%v != %v", r2, result)
	}

	time.Sleep(3 * time.Second)

	_, err = b.GetResult(result.GetBackendKey())
	if !yerrors.IsEqual(err, yerrors.ErrTypeNilResult) {
		t.Fatal("err != ErrNilResult")
	}

}

func TestMongoBackend2(t *testing.T) {
	// 测试断线重连
	b := mongo2.NewMongoBackend("127.0.0.1", "27017", "", "", "test", "test", 0)
	result := message.NewResult("xx123")
	result.FuncReturn = nil
	b.Activate()
	err := b.SetResult(result, 0)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("重启mongo后通过 \" kill -CONT %d \" 继续运行测试\n", os.Getpid())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGCONT)
	<-quit

	r2, err := b.GetResult(result.GetBackendKey())
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%v", r2) != fmt.Sprintf("%v", result) {
		t.Fatalf("%v != %v", r2, result)
	}

}
