package worker

import (
	"errors"
	"fmt"
	"github.com/gojuukaze/YTask/v2/log"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/util"
	"reflect"
)

type WorkerInterface interface {
	Run(msg message.Message, result *message.Result) error
	WorkerName() string
}

type FuncWorker struct {
	F    interface{}
	Name string
}

func (f FuncWorker) Run(msg message.Message, result *message.Result) error {
	return runFunc(f.F, msg, result)
}
func (f FuncWorker) WorkerName() string {
	return f.Name
}

func runFunc(f interface{}, msg message.Message, result *message.Result) (err error) {
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
	inValue, err := util.GetCallInArgs(funcValue, msg.JsonArgs)
	if err != nil {
		return err
	}

	r := funcValue.Call(inValue)
	if len(r) > 0 {
		result.Status = message.ResultStatus.Success
		s, err2 := util.GoValuesToJson(r)
		if err2 != nil {
			log.YTaskLog.Error(err2)
		} else {
			result.JsonResult = s
		}
	}
	return
}
