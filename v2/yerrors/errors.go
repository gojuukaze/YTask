package yerrors

import (
	"fmt"
)

const (
	ErrTypeEmptyQuery      = 1 // 队列为空， broker获取任务时用到
	ErrTypeUnsupportedType = 2 // 不支持此参数类型
	ErrTypeOutOfRange      = 3 // 暂时没用
	ErrTypeNilResult       = 4 // 任务结果为空
	ErrTypeTimeOut         = 5 // broker，backend超时
	ErrTypeServerStop      = 6 // 服务已停止
)

func IsEqual(err error, errType int) bool {
	yerr, ok := err.(YTaskError)
	if !ok {
		return ok
	}
	if yerr.Type() == errType {
		return true
	}
	return false
}

type YTaskError interface {
	Error() string
	Type() int
}

type ErrEmptyQuery struct {
}

func (e ErrEmptyQuery) Error() string {
	return "YTask: empty query"
}

func (e ErrEmptyQuery) Type() int {
	return ErrTypeEmptyQuery
}

type ErrUnsupportedType struct {
	T string
}

func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("YTask: UnsupportedType: %s", e.T)
}

func (e ErrUnsupportedType) Type() int {
	return ErrTypeUnsupportedType
}

type ErrOutOfRange struct {
}

func (e ErrOutOfRange) Error() string {
	return "YTask: index out of range"
}

func (e ErrOutOfRange) Type() int {
	return ErrTypeOutOfRange
}

type ErrNilResult struct {
}

func (e ErrNilResult) Error() string {
	return "YTask: nil result"
}

func (e ErrNilResult) Type() int {
	return ErrTypeNilResult
}

type ErrTimeOut struct {
}

func (e ErrTimeOut) Error() string {
	return "YTask: timeout"
}

func (e ErrTimeOut) Type() int {
	return ErrTypeTimeOut
}

type ErrServerStop struct {
}

func (e ErrServerStop) Error() string {
	return "YTask: server stop"
}

func (e ErrServerStop) Type() int {
	return ErrTypeServerStop
}
