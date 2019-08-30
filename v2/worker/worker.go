package worker

import (
	"errors"
	"fmt"
	"github.com/gojuukaze/YTask/v2/controller"
	"github.com/gojuukaze/YTask/v2/log"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/util"
	"reflect"
)

type WorkerInterface interface {
	Run(ctl *controller.TaskCtl, jsonArgs string, result *message.Result) error
	WorkerName() string
}

type FuncWorker struct {
	F    interface{}
	Name string
}

func (f FuncWorker) Run(ctl *controller.TaskCtl, jsonArgs string, result *message.Result) error {
	return runFunc(f.F, ctl, jsonArgs, result)
}
func (f FuncWorker) WorkerName() string {
	return f.Name
}

func runFunc(f interface{}, ctl *controller.TaskCtl, jsonArgs string, result *message.Result) (err error) {
	defer func() {
		e := recover()
		if e != nil {
			result.Status = message.ResultStatus.Failure
			t, ok := e.(error)
			if ok {
				err = t
			} else {
				err = errors.New(fmt.Sprintf("%v", e))
			}
		}
	}()
	funcValue := reflect.ValueOf(f)
	funcType := reflect.TypeOf(f)
	var inStart = 0
	var inValue []reflect.Value
	if funcType.NumIn() > 0 && funcType.In(0) == reflect.TypeOf(&controller.TaskCtl{}) {
		inStart = 1
	}

	inValue, err = util.GetCallInArgs(funcValue, jsonArgs, inStart)
	if err != nil {
		return
	}
	if inStart == 1 {
		inValue = append(inValue, reflect.Value{})
		copy(inValue[1:], inValue)
		inValue[0] = reflect.ValueOf(ctl)

	}

	funcOut := funcValue.Call(inValue)

	if ctl.GetError() != nil {
		err = ctl.GetError()
	} else {
		result.Status = message.ResultStatus.Success
		if len(funcOut) > 0 {
			s, err2 := util.GoValuesToJson(funcOut)
			if err2 != nil {
				log.YTaskLog.Error(err2)
			} else {
				result.JsonResult = s
			}
		}
	}

	return
}
