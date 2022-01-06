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
	Run(ctl *controller.TaskCtl, funcArgs []string, result *message.Result) error
	WorkerName() string
	After(ctl *controller.TaskCtl, funcArgs []string, result *message.Result) error
}

type FuncWorker struct {
	Func         interface{} // 执行的函数
	CallbackFunc interface{} // 回调函数
	Name         string
	Logger       log.LoggerInterface
}

func (f *FuncWorker) Run(ctl *controller.TaskCtl, funcArgs []string, result *message.Result) error {
	return runFunc(f.Func, ctl, funcArgs, result, false, f.Logger)
}

func (f *FuncWorker) After(ctl *controller.TaskCtl, funcArgs []string, result *message.Result) error {
	if f.CallbackFunc != nil {
		return runFunc(f.CallbackFunc, ctl, funcArgs, result, true, f.Logger)

	}
	return nil
}

func (f *FuncWorker) WorkerName() string {
	return f.Name
}

// isCallBack: 是否是回调函数
func runFunc(f interface{}, ctl *controller.TaskCtl, funcArgs []string, result *message.Result, isCallBack bool, logger log.LoggerInterface) (err error) {
	defer func() {
		e := recover()
		if e != nil {
			t, ok := e.(error)
			if ok {
				err = t
			} else {
				err = errors.New(fmt.Sprintf("%v", e))
			}
			if !isCallBack {
				result.Status = message.ResultStatus.Failure
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

	inValue, err = util.GetCallInArgs(funcValue, funcArgs, inStart)
	if err != nil {
		return
	}
	if inStart == 1 {
		inValue = append(inValue, reflect.Value{})
		copy(inValue[1:], inValue)
		inValue[0] = reflect.ValueOf(ctl)
	}

	if isCallBack {
		inValue[len(inValue)-1]=reflect.ValueOf(result)

	}


	funcOut := funcValue.Call(inValue)

	if isCallBack {
		return
	}
	err = ctl.GetError()
	if err == nil {
		result.Status = message.ResultStatus.Success
		if len(funcOut) > 0 {
			re, err2 := util.GoValuesToYJsonSlice(funcOut)
			if err2 != nil {
				//log.YTaskLog.Error(err2)
				logger.Error(err2.Error())
			} else {
				result.FuncReturn = re
			}
		}
	}

	return
}
