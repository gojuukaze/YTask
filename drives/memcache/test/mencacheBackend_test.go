package test

import (
	"fmt"
	"github.com/gojuukaze/YTask/v3/core/message"
	"github.com/gojuukaze/YTask/v3/core/yerrors"
	"github.com/gojuukaze/YTask/v3/drives/memcache"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestMemcachedBackend(t *testing.T) {
	b := memcache.NewMemCacheBackend([]string{"127.0.0.1:11211"}, 1)
	result := message.NewResult("xx123")
	result.FuncReturn = nil
	b.Activate()
	err := b.SetResult(result, 2)
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

	time.Sleep(2 * time.Second)

	_, err = b.GetResult(result.GetBackendKey())
	if !yerrors.IsEqual(err, yerrors.ErrTypeNilResult) {
		t.Fatal("err != ErrNilResult")

	}

}

func TestMemcachedBackend2(t *testing.T) {
	b := memcache.NewMemCacheBackend([]string{"127.0.0.1:11211"}, 1)
	result := message.NewResult("xx123")
	result.FuncReturn = nil
	b.Activate()
	err := b.SetResult(result, 100)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("重启memcache后通过 \" kill -CONT %d \" 继续运行测试\n", os.Getpid())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGCONT)
	<-quit

	err = b.SetResult(result, 100)
	if err != nil {
		t.Fatal(err)
	}

}
