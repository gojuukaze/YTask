package test

import (
	"fmt"
	"github.com/gojuukaze/YTask/v3/drives/redis"
	"github.com/gojuukaze/YTask/v3/message"
	"github.com/gojuukaze/YTask/v3/yerrors"
	"testing"
	"time"
)

func TestRedisBackend(t *testing.T) {
	b := redis.NewRedisBackend("127.0.0.1", "6379", "", 0, 1)
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
