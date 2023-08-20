package wits

import (
	"github.com/shura1014/logger"
)

var applog *logger.Logger

const (
	DebugLevel = logger.DebugLevel

	InfoLevel  = logger.InfoLevel
	WarnLevel  = logger.WarnLevel
	ErrorLevel = logger.ErrorLevel
	TEXT       = logger.TEXT
)

func SetLog(l *logger.Logger) {
	applog = l
}

func init() {
	applog = logger.Default(GetLogPrefix())
}

func Info(msg any, a ...any) {
	applog.DoPrint(appCtx, InfoLevel, msg, logger.GetFileNameAndLine(0), a...)
}

func Debug(msg any, a ...any) {
	applog.DoPrint(appCtx, DebugLevel, msg, logger.GetFileNameAndLine(0), a...)
}

func Error(msg any, a ...any) {
	applog.DoPrint(appCtx, ErrorLevel, msg, logger.GetFileNameAndLine(0), a...)
}

func Warn(msg any, a ...any) {
	applog.DoPrint(appCtx, WarnLevel, msg, logger.GetFileNameAndLine(0), a...)
}

func Text(msg any, a ...any) {
	applog.DoPrint(appCtx, TEXT, msg, logger.GetFileNameAndLine(0), a...)
}

func Fatal(msg any, a ...any) {
	applog.DoPrint(appCtx, ErrorLevel, msg, logger.GetFileNameAndLine(0), a...)
}
