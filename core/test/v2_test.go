package test

import (
	ytask "github.com/gojuukaze/YTask/core"
	"github.com/gojuukaze/YTask/v3/core/consts"
	"github.com/gojuukaze/YTask/v3/core/message"
	"github.com/gojuukaze/YTask/v3/core/server"
	"github.com/gojuukaze/YTask/v3/core/util/yjson"
	jsoniter "github.com/json-iterator/go"
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
	// 测试完后改回去，否则运行多个测试时会影响后面测试
	defer func() {
		consts.UserV2Name = false
		yjson.YJson = jsoniter.Config{
			EscapeHTML:             true,
			ValidateJsonRawMessage: true,
			TagKey:                 "yjson",
		}.Froze()
	}()
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
