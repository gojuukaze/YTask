package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"path/filepath"
	"runtime"
	"strings"
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

// 记录行号的hook
type lineNumHook struct {
}

func (hook lineNumHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook lineNumHook) Fire(entry *logrus.Entry) error {
	pc := make([]uintptr, 6, 6)
	cnt := runtime.Callers(6, pc)
	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		if !strings.Contains(name, "github.com/sirupsen/logrus") {
			file, line := fu.FileLine(pc[i] - 1)
			p, f := filepath.Split(file)
			l := strings.Split(p, "/")
			i := 0
			for ; i < len(l); i++ {
				if l[i] == "YTask" {
					break
				}
			}
			file = strings.Join(l[i:], "/") + f
			entry.Data["file"] = file + ":" + fmt.Sprintf("%d", line)
			entry.Data["func"] = path.Base(name)
			break
		}
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
