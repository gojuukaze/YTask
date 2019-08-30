package test

import (
	"github.com/gojuukaze/YTask/v2/backends"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/yerrors"
	"testing"
	"time"
)

func TestRedisBackend(t *testing.T) {
	b := backends.NewRedisBackend("127.0.0.1", "6379", "", 0, 1)
	result := message.NewResult("xx123")
	result.JsonResult = "[]"
	b.Activate()
	err := b.SetResult(result, 2)
	if err != nil {
		t.Fatal(err)
	}

	r2, err := b.GetResult(result.GetBackendKey())
	if err != nil {
		t.Fatal(err)
	}
	if r2 != result {
		t.Fatalf("%v != %v", r2, result)
	}

	time.Sleep(2 * time.Second)

	_, err = b.GetResult(result.GetBackendKey())
	if !yerrors.Compare(err, yerrors.ErrTypeNilResult) {
		t.Fatal("err != ErrNilResult")

	}

}
