package worker

import (
	"github.com/gojuukaze/YTask/v2/message"
	"reflect"
)

type WorkerInterface interface {
	Run(msg message.Message) error
	WorkerName() string
}

type FuncWorker struct {
	F    interface{}
	Name string
}

func (f FuncWorker) Run(msg message.Message) error {
	return runFunc(f.F, msg)
}
func (f FuncWorker) WorkerName() string {
	return f.Name
}

func runFunc(f interface{}, msg message.Message) error {
	funcValue := reflect.ValueOf(f)
	inValue, err := GetCallInArgs(funcValue, msg.JsonArgs)
	if err != nil {
		return err
	}

	r := funcValue.Call(inValue)
	if r[0].IsNil() {
		return nil
	}
	return r[0].Interface().(error)

}
