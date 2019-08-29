package test

import (
	"fmt"
	"github.com/gojuukaze/YTask/v2/config"
	"github.com/gojuukaze/YTask/v2/controller"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/server"
	"github.com/gojuukaze/YTask/v2/util"
	"github.com/gojuukaze/YTask/v2/worker"

	"github.com/tidwall/gjson"
	"reflect"
	"testing"
)

func TestGoArgsToJson(t *testing.T) {

	var (
		a int     = 1
		b int64   = 259933429192721385
		c uint    = 3
		d uint64  = 4
		e float32 = 5.5
		f float64 = 133.7976931348623
		g bool    = true
		h string  = "TestJsonArgs"
	)
	jsonStr, _ := util.GoArgsToJson(a, b, c, d, e, f, g, h)
	var base = [][]string{
		{"int", "1"},
		{"int64", "259933429192721385"},
		{"uint", "3"},
		{"uint64", "4"},
		{"float32", "5.5"},
		{"float64", "133.7976931348623"},
		{"bool", "true"},
		{"string", "TestJsonArgs"},
	}

	json := gjson.Parse(jsonStr)
	for i, j := range json.Array() {
		if j.Get("type").String() != base[i][0] {
			t.Fatalf("'%v'.type != %s", j, base[i][0])
		}
		if j.Get("value").String() != base[i][1] {
			t.Fatalf("'%v'.value != %s", j, base[i][1])
		}
	}

}

func TestGetCallInArgs(t *testing.T) {
	testFunc := func(int, int64, uint, uint64, float32, float64, bool, string) {}

	var (
		a int     = 1
		b int64   = 259933429192721385
		c uint    = 3
		d uint64  = 259933429192721385
		e float32 = 5.5
		f float64 = 133.7976931348623
		g bool    = true
		h string  = "TestJsonArgs"
	)
	jsonStr, _ := util.GoArgsToJson(a, b, c, d, e, f, g, h)
	inValues, _ := util.GetCallInArgs(reflect.ValueOf(testFunc), jsonStr, 0)
	var base = []reflect.Value{
		reflect.ValueOf(a),
		reflect.ValueOf(b),
		reflect.ValueOf(c),
		reflect.ValueOf(d),
		reflect.ValueOf(e),
		reflect.ValueOf(f),
		reflect.ValueOf(g),
		reflect.ValueOf(h),
	}
	for i, v := range inValues {
		if v.Interface() != base[i].Interface() {
			t.Fatalf("%v!=%v", v, base[i])
		}
	}

}

func TestGetCallInArgs2(t *testing.T) {
	testFunc := func() {}

	inValues, err := util.GetCallInArgs(reflect.ValueOf(testFunc), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(inValues) != 0 {
		t.Fatal("len(inValues)!=0")
	}

}

func TestRunFunc(t *testing.T) {
	var (
		ra int     = 1
		rb int64   = 259933429192721385
		rc uint    = 3
		rd uint64  = 4
		re float32 = 5.5
		rf float64 = 133.7976931348623
		rg bool    = true
		rh string  = "TestJsonArgs"
	)
	f := func(a, b int) (int, int, int64, uint, uint64, float32, float64, bool, string) {
		return a + b, ra, rb, rc, rd, re, rf, rg, rh
	}

	w := worker.FuncWorker{
		F: f,
	}
	s, _ := util.GoArgsToJson(12, 33)
	msg := message.NewMessage(controller.NewTaskCtl())
	msg.JsonArgs = s
	result := message.Result{}
	err := w.Run(&msg.TaskCtl, msg.JsonArgs, &result)
	if err != nil {
		t.Fatal(err)
	}

	var base = []reflect.Value{
		reflect.ValueOf(int(45)),
		reflect.ValueOf(ra),
		reflect.ValueOf(rb),
		reflect.ValueOf(rc),
		reflect.ValueOf(rd),
		reflect.ValueOf(re),
		reflect.ValueOf(rf),
		reflect.ValueOf(rg),
		reflect.ValueOf(rh),
	}
	for i, v := range base {
		inter, _ := result.GetInterface(i)
		if inter != v.Interface() {
			t.Fatalf("%v != %v", inter, v)
		}
	}
}

func cc1(ctl controller.TaskCtl) {
	ctl.RetryCount = 777
}

func cc2(ctl *controller.TaskCtl) {
	ctl.RetryCount = 777
}
func TestX(t *testing.T) {
	s:=server.NewServer(config.Config{})
	client := server.NewClient(&s)
	fmt.Printf("%+v\n", client)
	c2 := client.SetCtl(client.RetryCount, 123)
	fmt.Printf("%+v\n", c2)

	fmt.Printf("%+v\n", client)

}
