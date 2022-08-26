package backends

import (
	"github.com/gojuukaze/YTask/v3/drive"
	"github.com/gojuukaze/YTask/v3/message"
	"github.com/gojuukaze/YTask/v3/util/yjson"
	"github.com/gojuukaze/YTask/v3/yerrors"
)

type LocalBackend struct {
	client drive.LocalDrive
}

func NewLocalBackend() LocalBackend {
	return LocalBackend{}
}

func (l *LocalBackend) Activate() {
	l.client = drive.NewLocalDrive(false)
}

func (l *LocalBackend) SetResult(result message.Result, exTime int) error {
	b, err := yjson.YJson.Marshal(result)

	if err != nil {
		return err
	}
	err = l.client.Set(result.GetBackendKey(), b, exTime)
	return err
}

func (l *LocalBackend) GetResult(key string) (message.Result, error) {
	var result message.Result

	b, err := l.client.Get(key)
	if err != nil {
		if err == drive.NilResultError {
			return result, yerrors.ErrNilResult{}
		}
		return result, err
	}

	err = yjson.YJson.Unmarshal(b, &result)
	return result, err
}

func (l *LocalBackend) SetPoolSize(i int) {

}

func (l LocalBackend) GetPoolSize() int {
	return 0
}

func (l *LocalBackend) Clone() BackendInterface {
	return &LocalBackend{}
}
