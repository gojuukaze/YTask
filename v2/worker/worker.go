package worker

import (
	"github.com/gojuukaze/YTask/v2/message"
	"reflect"
)

type RunnerInterface interface {
	Run(msg message.Message) error
}
type WorkerInterface interface {
	RunnerInterface
	WorkerName() string
}

type FuncWorker struct {
	F    interface{}
	Name string
}

func (f FuncWorker) Run(msg message.Message) error {
	funcValue := reflect.ValueOf(f.F)
	r := funcValue.Call([]reflect.Value{reflect.ValueOf(msg)})
	if r[0].Interface() == nil {
		return nil
	} else {
		return r[0].Interface().(error)
	}
}
func (f FuncWorker) WorkerName() string {
	return f.Name
}

type StructWorker struct {
	S    interface{}
	Name string
}

func (s StructWorker) Run(msg message.Message) error {
	return s.S.(RunnerInterface).Run(msg)
}

func (s StructWorker) WorkerName() string {
	return s.Name
}
