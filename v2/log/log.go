package log

import (
	"github.com/sirupsen/logrus"
)

var YTaskLog *logrus.Logger


func init() {

	YTaskLog = logrus.New()
	YTaskLog.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	YTaskLog.SetLevel(logrus.InfoLevel)

}
