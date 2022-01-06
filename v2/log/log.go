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
	serverName, ok := entry.Data["server"]
	s := ""
	if ok {
		s = fmt.Sprintf("server[%s", serverName)

	}

	goroutineName, ok := entry.Data["goroutine"]
	if ok {
		goroutineName = fmt.Sprintf("|%s", goroutineName)

	}
	if !ok{
		goroutineName=""
	}
	delete(entry.Data, "goroutine")
	delete(entry.Data, "server")

	entry.Message = fmt.Sprintf("%s%s]: %s", s,goroutineName, entry.Message)

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

type LoggerInterface interface {
	Debug(string)
	DebugWithField(string, string, interface{})
	Info(string)
	InfoWithField(string, string, interface{})
	Warn(string)
	WarnWithField(string, string, interface{})
	Error(string)
	ErrorWithField(string, string, interface{})
	Fatal(string)
	FatalWithField(string, string, interface{})
	Panic(string)
	PanicWithField(string, string, interface{})
	SetLevel(string)
	Clone() LoggerInterface
}

type YTaskLogger struct {
	logger *logrus.Logger
}

func NewYTaskLogger(logger *logrus.Logger) *YTaskLogger {
	return &YTaskLogger{
		logger: logger,
	}
}

func (yl *YTaskLogger) Debug(msg string) {
	yl.logger.Debug(msg)
}

func (yl *YTaskLogger) DebugWithField(msg string, key string, val interface{}) {
	yl.logger.WithField(key, val).Debug(msg)
}

func (yl *YTaskLogger) Info(msg string) {
	yl.logger.Info(msg)
}

func (yl *YTaskLogger) InfoWithField(msg string, key string, val interface{}) {
	yl.logger.WithField(key, val).Info(msg)
}

func (yl *YTaskLogger) Warn(msg string) {
	yl.logger.Warn(msg)
}

func (yl *YTaskLogger) WarnWithField(msg string, key string, val interface{}) {
	yl.logger.WithField(key, val).Warn(msg)
}

func (yl *YTaskLogger) Error(msg string) {
	yl.logger.Error(msg)
}

func (yl *YTaskLogger) ErrorWithField(msg string, key string, val interface{}) {
	yl.logger.WithField(key, val).Error(msg)
}

func (yl *YTaskLogger) Fatal(msg string) {
	yl.logger.Fatal(msg)
}

func (yl *YTaskLogger) FatalWithField(msg string, key string, val interface{}) {
	yl.logger.WithField(key, val).Fatal(msg)
}

func (yl *YTaskLogger) Panic(msg string) {
	yl.logger.Panic(msg)
}

func (yl *YTaskLogger) PanicWithField(msg string, key string, val interface{}) {
	yl.logger.WithField(key, val).Panic(msg)
}

func (yl *YTaskLogger) SetLevel(level string)  {
	switch level {
	case "debug":
		yl.logger.SetLevel(logrus.DebugLevel)
	case "info":
		yl.logger.SetLevel(logrus.InfoLevel)
	case "warn":
		yl.logger.SetLevel(logrus.WarnLevel)
	case "error":
		yl.logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		yl.logger.SetLevel(logrus.FatalLevel)
	case "panic":
		yl.logger.SetLevel(logrus.PanicLevel)
	default:
		yl.logger.SetLevel(logrus.InfoLevel)
	}
}

func (yl *YTaskLogger) Clone() LoggerInterface {
	return yl
}
