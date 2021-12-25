package test

import (
	"fmt"
	"testing"

	"github.com/vua/YTask/v2/backends"
	"github.com/vua/YTask/v2/message"
)

func TestMongoBackend(t *testing.T) {
	b := backends.NewMongoBackend("127.0.0.1", "27017", "", "", "task", "task")
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

	// mongo不支持过期时间

	//time.Sleep(2 * time.Second)
	//
	//_, err = b.GetResult(result.GetBackendKey())
	//if !yerrors.IsEqual(err, yerrors.ErrTypeNilResult) {
	//	t.Fatal("err != ErrNilResult")
	//
	//}

}
