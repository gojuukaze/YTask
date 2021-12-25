package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/vua/YTask/v2/backends"
	"github.com/vua/YTask/v2/message"
	"github.com/vua/YTask/v2/yerrors"
)

func TestMemcacheBackend(t *testing.T) {
	b := backends.NewMemCacheBackend("127.0.0.1", "11211", 1)
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
