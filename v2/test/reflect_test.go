package test

import (
	"fmt"
	"github.com/gojuukaze/YTask/v2/controller"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/util"
	"github.com/gojuukaze/YTask/v2/worker"
	"reflect"
	"testing"
)

type s1 struct {
	A int
	B int64
	C uint64
	D float64
	F []int64
	G []float64
}
type s2 struct {
	A int
	B string
}
type s3 struct {
	S2 s2
	C  bool
	D  float64
}

var (
	a int       = 1
	b int64     = 259933429192721385
	c uint      = 3
	d uint64    = 4
	e float32   = 5.5
	f float64   = 133.7976931348623
	g bool      = true
	h string    = "TestJsonArgs"
	i []int     = []int{123, 4456, 56756, 234, 123, 4, 5, 6, 7, 812, 123, 345, 756, 678, 7686, 7, 2, 23, 4}
	j []int64   = []int64{259933429192721385, 219933429192721385, 4}
	k []uint    = []uint{44, 56546, 2311, 567,}
	l []uint64  = []uint64{18446744073709551615, 18446744073709551600, 184467440737095516}
	m []float64 = []float64{445535.3321, 133.7976931348623}
	n []string  = []string{"", "YTask", "is", "good", "!!", " "}
	o           = s1{A: 12344, B: 4444444444444, C: 123789, D: 123.444456, G: []float64{677.4, 345.78221}}
	p           = s3{S2: s2{345, "ggggggg"}, C: true, D: 344,}
)

func TestGoVarToYJson(t *testing.T) {

	jsonSlice, _ := util.GoVarsToYJsonSlice(a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p)
	var base = []string{fmt.Sprintf("%+v", a),
		fmt.Sprintf("%+v", b),
		fmt.Sprintf("%+v", c),
		fmt.Sprintf("%+v", d),
		fmt.Sprintf("%+v", e),
		fmt.Sprintf("%+v", f),
		fmt.Sprintf("%+v", g),
		fmt.Sprintf(`"%+v"`, h),
		"[123,4456,56756,234,123,4,5,6,7,812,123,345,756,678,7686,7,2,23,4]",
		"[259933429192721385,219933429192721385,4]",
		"[44,56546,2311,567]",
		"[18446744073709551615,18446744073709551600,184467440737095516]",
		"[445535.3321,133.7976931348623]",
		`["","YTask","is","good","!!"," "]`,
		`{"A":12344,"B":4444444444444,"C":123789,"D":123.444456,"Func":null,"G":[677.4,345.78221]}`,
		`{"S2":{"A":345,"B":"ggggggg"},"C":true,"D":344}`,}
	for i, v := range base {
		if v != jsonSlice[i] {
			t.Fatalf("%+v != %+v", v, jsonSlice[i])
		}
	}

}

func TestGetCallInArgs(t *testing.T) {
	testFunc := func(int, int64, uint, uint64, float32, float64, bool, string, []int, []int64, []uint64, []string, s1, s3) {}

	jsonSlice, _ := util.GoVarsToYJsonSlice(a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p)
	inValues, _ := util.GetCallInArgs(reflect.ValueOf(testFunc), jsonSlice, 0)
	var base = []reflect.Value{
		reflect.ValueOf(a),
		reflect.ValueOf(b),
		reflect.ValueOf(c),
		reflect.ValueOf(d),
		reflect.ValueOf(e),
		reflect.ValueOf(f),
		reflect.ValueOf(g),
		reflect.ValueOf(h),
		reflect.ValueOf(i),
		reflect.ValueOf(j),
		reflect.ValueOf(k),
		reflect.ValueOf(l),
		reflect.ValueOf(m),
		reflect.ValueOf(n),
		reflect.ValueOf(o),
		reflect.ValueOf(p),
	}
	for i, v := range inValues {
		if reflect.DeepEqual(b, base[i]) {
			t.Fatalf("%v!=%v", v, base[i])
		}
	}

}

//
func TestGetCallInArgs2(t *testing.T) {
	testFunc := func() {}

	inValues, err := util.GetCallInArgs(reflect.ValueOf(testFunc), nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(inValues) != 0 {
		t.Fatal("len(inValues)!=0")
	}

}

func TestRunFunc(t *testing.T) {
	fun := func(aa, bb int) (int, int, int64, uint, uint64, float32, float64, bool, string, []int64, []uint64, []float64, []string, s1, s3) {
		return aa + bb, a, b, c, d, e, f, g, h, j, l, m, n, o, p
	}

	w := &worker.FuncWorker{
		Func: fun,
	}
	s, _ := util.GoVarsToYJsonSlice(12, 33)
	msg := message.NewMessage(controller.NewTaskCtl())
	msg.FuncArgs = s
	result := message.Result{}
	err := w.Run(&msg.TaskCtl, msg.FuncArgs, &result)
	if err != nil {
		t.Fatal(err)
	}
	var base = []interface{}{int(45), a, b, c, d, e, f, g, h, j, l, m, n, o, p,}
	var (
		raa int
		ra  int
		rb  int64
		rc  uint
		rd  uint64
		re  float32
		rf  float64
		rg  bool
		rh  string
		rj  []int64
		rl  []uint64
		rm  []float64
		rn  []string
		ro  s1
		rp  s3
	)

	var returnV = []interface{}{&raa, &ra, &rb, &rc, &rd, &re, &rf, &rg, &rh, &rj, &rl, &rm, &rn, &ro, &rp,}
	for i, v := range base {
		temp := returnV[i]
		err := result.Get(i, temp)
		if err != nil {
			t.Fatal(err)
		}
		var realV = reflect.ValueOf(temp).Elem()
		if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", realV) {
			t.Fatalf("%v != %v", v, realV)
		}
	}

	err = result.Gets(&raa, &ra, &rb, &rc, &rd, &re, &rf, &rg, &rh, &rj, &rl, &rm, &rn, &ro, &rp)
	if err != nil {
		t.Fatal(err)
	}
	returnV = []interface{}{raa, ra, rb, rc, rd, re, rf, rg, rh, rj, rl, rm, rn, ro, rp,}
	for i, v := range base {
		if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", returnV[i]) {
			t.Fatalf("%v != %v", v, returnV[i])
		}
	}

}
