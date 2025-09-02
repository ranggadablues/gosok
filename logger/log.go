package logger

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/term"
	"github.com/ranggadablues/gosok/common"
)

type ILogLevel interface {
	LogLevel(keyvals ...interface{})
}

type LogLevel struct {
}

func NewLogger() log.Logger {
	logger := term.NewLogger(os.Stdout, log.NewLogfmtLogger, ColorInit)
	logger = log.With(logger, "ts", log.DefaultTimestamp, "caller", log.DefaultCaller)
	return logger
}

func ColorInit(keyvals ...interface{}) term.FgBgColor {
	for i := 0; i < len(keyvals)-1; i += 2 {
		if keyvals[i] != "level" {
			continue
		}
		level := common.ToString(keyvals[i+1])
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
