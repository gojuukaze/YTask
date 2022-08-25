package test

import (
	ytask "github.com/gojuukaze/YTask/v3"
	"github.com/gojuukaze/YTask/v3/message"
	"github.com/gojuukaze/YTask/v3/server"
	"github.com/gojuukaze/YTask/v3/util/yjson"
	"strings"
	"testing"
)

// 测试v2兼容
func TestV2(t *testing.T) {

	s, _ := yjson.YJson.MarshalToString(message.NewMessage(message.NewMsgArgs()))
	// v3版本，msg序列化后不应包含TaskCtl
	i := strings.Index(s, "TaskCtl")
	if i != -1 {
		t.Fatalf("%s", s)
	}
	// 测试v3版本，队列名
	su := server.ServerUtils{}
	s = su.GetQueueName("a")
	i = strings.Index(s, "Query")
	if i != -1 {
		t.Fatalf("%s", s)
	}

	// 改为v2版本名字
	ytask.UseV2Name()
	s, _ = yjson.YJson.MarshalToString(message.NewMessage(message.NewMsgArgs()))
	i = strings.Index(s, "TaskCtl")
	if i == -1 {
		t.Fatalf("%s", s)
	}

	s = su.GetQueueName("a")
	i = strings.Index(s, "Query")
	if i == -1 {
		t.Fatalf("%s", s)
	}

}
