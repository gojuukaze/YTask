package util

import (
	"github.com/gojuukaze/YTask/v3/core/util/yjson"
	"reflect"
)

func GoVarToYJson(v interface{}) (string, error) {
	b, err := yjson.YJson.Marshal(v)
	return string(b), err
}

func GoVarsToYJsonSlice(args ...interface{}) ([]string, error) {
	var r = make([]string, len(args))
	for i, v := range args {
		yJsonStr, err := GoVarToYJson(v)
		if err != nil {
			return r, err
		}
		r[i] = yJsonStr
	}
	return r, nil
}

func GoValuesToYJsonSlice(values []reflect.Value) ([]string, error) {
	var r = make([]string, len(values))
	for i, v := range values {
		s, err := GoVarToYJson(v.Interface())
		if err != nil {
			return r, err
		}
		r[i] = s

	}
	return r, nil
}

func GetCallInArgs(funcValue reflect.Value, funcArgs []string, inStart int) ([]reflect.Value, error) {

	var inArgs = make([]reflect.Value, funcValue.Type().NumIn()-inStart)
	for i := inStart; i < funcValue.Type().NumIn(); i++ {
		if i-inStart >= len(funcArgs) {
			break
		}
		inType := funcValue.Type().In(i)
		inValue := reflect.New(inType)
		// yjson to go value
		err := yjson.YJson.Unmarshal([]byte(funcArgs[i-inStart]), inValue.Interface())

		if err != nil {
			return inArgs, err
		}
		inArgs[i-inStart] = inValue.Elem()
	}
	return inArgs, nil
}
