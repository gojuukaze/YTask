package util

import (
	"fmt"
	"github.com/gojuukaze/YTask/v2/consts"
	"github.com/gojuukaze/YTask/v2/errors"
	"github.com/tidwall/sjson"
	"reflect"
)

// return : [ {"type":"int", "value":123} ]
//
func GoArgsToJson(args ...interface{}) (string, error) {
	var s string
	for i, v := range args {
		_type := reflect.TypeOf(v).String()
		if !consts.SupportedTypes[_type] {
			return s, errors.ErrUnsupportedType{_type}
		}
		s, _ = sjson.Set(s, fmt.Sprintf("%d.type", i), _type)
		s, _ = sjson.Set(s, fmt.Sprintf("%d.value", i), v)

	}
	return s,nil
}