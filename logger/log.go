package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-kit/log/term"
	"github.com/ranggadablues/gosok/common"
)

type ILogLevel interface {
	LogInfoLevel(keyvals ...interface{})
	LogWarnLevel(keyvals ...interface{})
	LogErrorLevel(keyvals ...interface{})
	LogDebugLevel(keyvals ...interface{})
	LogDebugLevelWithCaller(msg string)
	UTC() *LogLevel
}

type LogLevel struct {
	logger log.Logger
	isUTC  bool
}

func NewLogger() ILogLevel {
	logger := setNewLogger(false)
	return &LogLevel{logger: logger, isUTC: false}
}

func (l *LogLevel) UTC() *LogLevel {
	l.isUTC = true
	return l
}

func (l *LogLevel) defaultLogTime() *LogLevel {
	if l.isUTC {
		l.logger = setNewLogger(l.isUTC)
	}
	return l
}

func setNewLogger(isUTC bool) log.Logger {
	logTime := log.DefaultTimestamp
	if isUTC {
		logTime = log.DefaultTimestampUTC
	}
	logger := term.NewLogger(os.Stdout, log.NewLogfmtLogger, ColorInit)
	logger = log.With(logger, "ts", logTime, "caller", log.Caller(4))
	return logger
}

func (l *LogLevel) LogInfoLevel(keyvals ...interface{}) {
	l.defaultLogTime()
	level.Info(l.logger).Log(keyvals...)
}

func (l *LogLevel) LogWarnLevel(keyvals ...interface{}) {
	l.defaultLogTime()
	level.Warn(l.logger).Log(keyvals...)
}

func (l *LogLevel) LogErrorLevel(keyvals ...interface{}) {
	l.defaultLogTime()
	level.Error(l.logger).Log(keyvals...)
}

func (l *LogLevel) LogDebugLevel(keyvals ...interface{}) {
	l.defaultLogTime()
	level.Debug(l.logger).Log(keyvals...)
}

func (l *LogLevel) LogDebugLevelWithCaller(msg string) {
	l.defaultLogTime()
	file, line, fn := getCallerInfo(3)
	level.Warn(l.logger).Log(
		"query", msg,
		"from", fmt.Sprintf("%s:%d", file, line),
		"func", fn,
	)
}

func ColorInit(keyvals ...interface{}) term.FgBgColor {
	for i := 0; i < len(keyvals)-1; i += 2 {
		if keyvals[i] != "level" {
			continue
		}
		level := common.ParseString(keyvals[i+1])
		switch level {
		case "debug":
			return term.FgBgColor{Fg: term.White, Bg: term.DarkGray}
		case "info":
			return term.FgBgColor{Fg: term.Green}
		case "warn":
			return term.FgBgColor{Fg: term.Yellow}
		case "error":
			return term.FgBgColor{Fg: term.Red}
		}
	}
	return term.FgBgColor{} // Default, no color
}

func getCallerInfo(skip int) (file string, line int, fn string) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "???", 0, "???"
	}

	fnDetails := runtime.FuncForPC(pc)
	fnName := "???"
	if fnDetails != nil {
		parts := strings.Split(fnDetails.Name(), "/")
		fnName = parts[len(parts)-1]
	}

	return file, line, fnName
}
