package log

import (
	"chatgpt-proxy/pkg/config"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"runtime"
)

type logger struct {
	entry *logrus.Entry
	// panic,fatal,error,warn,warning,info,debug,trace
	level       string
	printCaller bool
	caller      func() (file string, line int, funcName string, err error)
}

// 设置日志打印级别
func (l *logger) SetLevel(lvl string) {
	if lvl == "" {
		return
	}
	level, err := logrus.ParseLevel(lvl)
	if err == nil {
		l.level = lvl
		l.entry.Logger.Level = level
	}
}

// 设置日志输出位置
func (l *logger) SetOutput(writer io.Writer) {
	l.entry.Logger.SetOutput(writer)
}

// 设置是否打印调用信息
func (l *logger) SetPrintCaller(printCaller bool) {
	l.printCaller = printCaller
}
func (l *logger) SetCaller(caller func() (file string, line int, funcName string, err error)) {
	l.caller = caller
}

// 获取caller信息
func (l *logger) getCallerInfo(level logrus.Level) map[string]interface{} {
	mp := make(map[string]interface{})
	if l.printCaller == true || level != logrus.InfoLevel {
		file, line, funcName, err := l.caller()
		if err == nil {
			mp["file"] = fmt.Sprintf("%s:%d", file, line)
			mp["func"] = funcName
		}
	}
	return mp
}

func (l *logger) log(level logrus.Level, args ...interface{}) {
	l.entry.WithFields(l.getCallerInfo(level)).Log(level, args...)
}
func (l *logger) logf(level logrus.Level, format string, args ...interface{}) {
	l.entry.WithFields(l.getCallerInfo(level)).Logf(level, format, args...)
}
func (l *logger) Trace(args ...interface{}) {
	l.log(logrus.TraceLevel, args...)
}
func (l *logger) Debug(args ...interface{}) {
	l.log(logrus.DebugLevel, args...)
}
func (l *logger) Info(args ...interface{}) {
	l.log(logrus.InfoLevel, args...)
}
func (l *logger) Warning(args ...interface{}) {
	l.log(logrus.WarnLevel, args...)
}
func (l *logger) Error(args ...interface{}) {
	l.log(logrus.ErrorLevel, args...)
}
func (l *logger) Fatal(args ...interface{}) {
	l.log(logrus.FatalLevel, args...)
}
func (l *logger) Panic(args ...interface{}) {
	l.log(logrus.PanicLevel, args...)
}
func (l *logger) TraceF(format string, args ...interface{}) {
	l.logf(logrus.TraceLevel, format, args...)
}
func (l *logger) DebugF(format string, args ...interface{}) {
	l.logf(logrus.DebugLevel, format, args...)
}
func (l *logger) InfoF(format string, args ...interface{}) {
	l.logf(logrus.InfoLevel, format, args...)
}
func (l *logger) WarningF(format string, args ...interface{}) {
	l.logf(logrus.WarnLevel, format, args...)
}
func (l *logger) ErrorF(format string, args ...interface{}) {
	l.logf(logrus.ErrorLevel, format, args...)
}
func (l *logger) FatalF(format string, args ...interface{}) {
	l.logf(logrus.FatalLevel, format, args...)
}
func (l *logger) PanicF(format string, args ...interface{}) {
	l.logf(logrus.PanicLevel, format, args...)
}
func (l *logger) WithFields(fields map[string]interface{}) *logger {
	entry := l.entry.WithFields(fields)
	return &logger{entry: entry, level: l.level, printCaller: l.printCaller, caller: l.caller}
}

var log *logger

func NewLogger() *logger {
	log := logrus.New()
	log.SetLevel(logrus.WarnLevel)
	log.AddHook(&errorHook{})
	log.SetOutput(getRotateWriter())
	logger := &logger{
		entry:  logrus.NewEntry(log),
		caller: defaultCaller,
	}
	return logger
}

func init() {
	cnf := config.GetConf()
	log = NewLogger()
	log.SetLevel(cnf.Log.Level)
	log.SetOutput(getRotateWriter())
}

// 设置是否打印调用信息
func SetPrintCaller(printCaller bool) {
	log.printCaller = printCaller
}

func SetCaller(caller func() (file string, line int, funcName string, err error)) {
	log.caller = caller
}

func defaultCaller() (file string, line int, funcName string, err error) {
	pc, f, l, ok := runtime.Caller(4)
	if !ok {
		err = errors.New("caller failure")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	file, line = f, l
	return
}

func Trace(args ...interface{}) {
	log.log(logrus.TraceLevel, args...)
}
func Debug(args ...interface{}) {
	log.log(logrus.DebugLevel, args...)
}
func Info(args ...interface{}) {
	log.log(logrus.InfoLevel, args...)
}
func Warning(args ...interface{}) {
	log.log(logrus.WarnLevel, args...)
}
func Error(args ...interface{}) {
	log.log(logrus.ErrorLevel, args...)
}
func Fatal(args ...interface{}) {
	log.log(logrus.FatalLevel, args...)
}
func Panic(args ...interface{}) {
	log.log(logrus.PanicLevel, args...)
}
func TraceF(format string, args ...interface{}) {
	log.logf(logrus.TraceLevel, format, args...)
}
func DebugF(format string, args ...interface{}) {
	log.logf(logrus.DebugLevel, format, args...)
}
func InfoF(format string, args ...interface{}) {
	log.logf(logrus.InfoLevel, format, args...)
}
func WarningF(format string, args ...interface{}) {
	log.logf(logrus.WarnLevel, format, args...)
}
func ErrorF(format string, args ...interface{}) {
	log.logf(logrus.ErrorLevel, format, args...)
}
func FatalF(format string, args ...interface{}) {
	log.logf(logrus.FatalLevel, format, args...)
}
func PanicF(format string, args ...interface{}) {
	log.logf(logrus.PanicLevel, format, args...)
}
func WithFields(fields map[string]interface{}) *logger {
	entry := log.entry.WithFields(fields)
	return &logger{entry: entry, level: log.level, printCaller: log.printCaller, caller: log.caller}
}
