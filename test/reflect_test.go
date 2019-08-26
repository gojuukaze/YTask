package test

import (
	"github.com/gojuukaze/YTask/v2/message"
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
	inValues, _ := worker.GetCallInArgs(reflect.ValueOf(testFunc), jsonStr)
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

	inValues, err := worker.GetCallInArgs(reflect.ValueOf(testFunc), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(inValues) != 0 {
		t.Fatal("len(inValues)!=0")
	}

}

func TestRunFunc(t *testing.T) {
	var result int
	f := func(a, b int) error {
		result = a + b
		return nil
	}

	w := worker.FuncWorker{
		F: f,
	}
	s, _ := util.GoArgsToJson(12, 33)
	err := w.Run(message.Message{
		JsonArgs: s,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result != 45 {
		t.Fatal("result!=45")

	}
}
