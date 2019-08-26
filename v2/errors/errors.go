package errors

import (
	"fmt"
)

const (
	ErrTypeEmptyQuery      = 1
	ErrTypeUnsupportedType = 2
)

type YtaskError interface {
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
	return fmt.Sprintf("UnsupportedType: %s", e.T)
}

func (e ErrUnsupportedType) Type() int {
	return ErrTypeUnsupportedType
}
