package worker

import (
	"fmt"
	"github.com/gojuukaze/YTask/v2/consts"
	"github.com/gojuukaze/YTask/v2/errors"
	"github.com/tidwall/gjson"
	"strings"

	"reflect"
)

func GetValue(inType reflect.Type, jsonType string, jsonValue gjson.Result) (reflect.Value, error) {
	var v = reflect.New(inType)
	if !consts.SupportedTypes[jsonType] {
		return v, errors.ErrUnsupportedType{jsonType}
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
		return v.Elem(), errors.ErrUnsupportedType{jsonType}
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
