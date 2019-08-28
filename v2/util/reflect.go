package util

import (
	"fmt"
	"github.com/gojuukaze/YTask/v2/consts"
	"github.com/gojuukaze/YTask/v2/yerrors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"reflect"
	"strings"
)

var typeMap = map[string]reflect.Type{
	"int":     reflect.TypeOf(int(0)),
	"int64":   reflect.TypeOf(int64(0)),
	"uint":    reflect.TypeOf(uint(0)),
	"uint64":  reflect.TypeOf(uint64(0)),
	"float32": reflect.TypeOf(float32(0)),
	"float64": reflect.TypeOf(float64(0)),
	"bool":    reflect.TypeOf(false),
	"string":  reflect.TypeOf(""),
}

// return : [ {"type":"int", "value":123} ]
//
func GoArgsToJson(args ...interface{}) (string, error) {
	var s string
	for i, v := range args {
		_type := reflect.TypeOf(v).String()
		if !consts.SupportedTypes[_type] {
			return s, yerrors.ErrUnsupportedType{_type}
		}
		s, _ = sjson.Set(s, fmt.Sprintf("%d.type", i), _type)
		s, _ = sjson.Set(s, fmt.Sprintf("%d.value", i), v)

	}
	return s, nil
}

func GoValuesToJson(values []reflect.Value) (string, error) {
	var s string
	for i, v := range values {
		_type := v.Type().String()
		if !consts.SupportedTypes[_type] {
			return s, yerrors.ErrUnsupportedType{_type}
		}
		s, _ = sjson.Set(s, fmt.Sprintf("%d.type", i), _type)
		s, _ = sjson.Set(s, fmt.Sprintf("%d.value", i), v.Interface())

	}
	return s, nil
}

func GetValueFromJson(j string) (reflect.Value, error) {
	jsonType := gjson.Get(j, "type").String()
	v := gjson.Get(j, "value")
	return GetValue(typeMap[jsonType], jsonType, v)
}
func GetValue(inType reflect.Type, jsonType string, jsonValue gjson.Result) (reflect.Value, error) {
	var v = reflect.New(inType)
	if !consts.SupportedTypes[jsonType] {
		return v, yerrors.ErrUnsupportedType{jsonType}
	}
	if strings.HasPrefix(jsonType, "int") {
		v.Elem().SetInt(jsonValue.Int())
	} else if strings.HasPrefix(jsonType, "bool") {
		v.Elem().SetBool(jsonValue.Bool())
	} else if strings.HasPrefix(jsonType, "uint") {
		v.Elem().SetUint(jsonValue.Uint())
	} else if strings.HasPrefix(jsonType, "float") {
		v.Elem().SetFloat(jsonValue.Float())
	} else if strings.HasPrefix(jsonType, "string") {
		v.Elem().SetString(jsonValue.String())
	} else {
		return v.Elem(), yerrors.ErrUnsupportedType{jsonType}
	}
	return v.Elem(), nil
}

func GetCallInArgs(funcValue reflect.Value, jsonArgs string) ([]reflect.Value, error) {
	var inArgs = make([]reflect.Value, funcValue.Type().NumIn())
	j := gjson.Parse(jsonArgs)
	for i := 0; i < funcValue.Type().NumIn(); i++ {
		inType := funcValue.Type().In(i)
		jsonValue := j.Get(fmt.Sprintf("%d.value", i))
		jsonType := j.Get(fmt.Sprintf("%d.type", i)).String()
		inValue, err := GetValue(inType, jsonType, jsonValue)
		if err != nil {
			return inArgs, err
		}
		inArgs[i] = inValue
	}
	return inArgs, nil
}
