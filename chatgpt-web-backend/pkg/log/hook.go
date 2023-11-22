package log

import (
	nativeLog "log"

	"github.com/sirupsen/logrus"
)

type errorHook struct {
}

func (*errorHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}
}
func (*errorHook) Fire(entry *logrus.Entry) error {
	nativeLog.Println(entry.Message, entry.Data)
	return nil
}
