package backends

import "github.com/gojuukaze/YTask/v2/message"

type BackendInterface interface {
	SetResult(result message.Result, exTime int) error
	GetResult(key string) (message.Result, error)

	Activate()
	SetPoolSize(int)
	GetPoolSize() int
}
