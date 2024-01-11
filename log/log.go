package log

import (
	"fmt"
	"io"
	"log"
)

const (
	LogLevelNull    = 0
	LogLevelTrace   = 1
	LogLevelDebug   = 2
	LogLevelInfo    = 3
	LogLevelWarning = 4
	LogLevelError   = 5
	LogLevelFatal   = 6
)

// 默认log级别
var defaultLogLevel uint8 = LogLevelDebug

type Logger struct {
	*log.Logger
}

func SetFlags(flag int) {
	log.SetFlags(flag)
}

func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

func SetLevel(level uint8) {
	defaultLogLevel = level
}

func Trace(format string, v ...interface{}) {
	if defaultLogLevel > LogLevelTrace {
		return
	}
	log.Output(2, string("[TRACE] ")+fmt.Sprintf(format, v...))
}

func Debug(format string, v ...interface{}) {
	if defaultLogLevel > LogLevelDebug {
		return
	}
	log.Output(2, string("[DEBUG] ")+fmt.Sprintf(format, v...))
}

func Info(format string, v ...interface{}) {
	if defaultLogLevel > LogLevelInfo {
		return
	}
	log.Output(2, string("[INFO] ")+fmt.Sprintf(format, v...))
}

func Warning(format string, v ...interface{}) {
	if defaultLogLevel > LogLevelWarning {
		return
	}
	log.Output(2, string("[WARNING] ")+fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	if defaultLogLevel > LogLevelError {
		return
	}
	log.Output(2, string("[ERROR] ")+fmt.Sprintf(format, v...))
}
