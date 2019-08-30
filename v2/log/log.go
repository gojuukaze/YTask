package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

var YTaskLog *logrus.Logger

// 记录行号的hook
type YTaskHook struct {
}

func (hook YTaskHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook YTaskHook) Fire(entry *logrus.Entry) error {
	goroutineName,ok := entry.Data["goroutine"]
	delete(entry.Data, "goroutine")
	if ok {
		entry.Message = fmt.Sprintf("goroutine[%v]: %s", goroutineName, entry.Message)

	}
	return nil
}

func init() {

	YTaskLog = logrus.New()
	YTaskLog.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	YTaskLog.SetLevel(logrus.InfoLevel)
	YTaskLog.AddHook(&YTaskHook{})

}
