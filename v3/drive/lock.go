package drive

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type UnsafeFileLock struct {
	filePath string
}

func NewFileLock(s string) UnsafeFileLock {
	return UnsafeFileLock{filepath.Join(os.TempDir(), s)}
}
func (l UnsafeFileLock) Init() {
	os.Remove(l.filePath)
}
func (l UnsafeFileLock) Lock() error {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*400)
	for {
		f, err := os.OpenFile(l.filePath, os.O_CREATE|os.O_EXCL, os.FileMode(0600))
		if err == nil {
			f.Close()
			return nil
		}
		select {
		case <-ctx.Done():
			return errors.New("lock timeout")
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func (l UnsafeFileLock) Unlock() {
	os.Remove(l.filePath)
}
