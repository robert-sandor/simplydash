package internal

import (
	"io"
	"log"
	"os"
	"strings"
)

type LogLevel int

const (
	Error = 1
	Info  = 2
	Debug = 3
)

type Logger struct {
	Error *log.Logger
	Info  *log.Logger
	Debug *log.Logger
}

var Log *Logger

func NewLogger(logLevel LogLevel) *Logger {
	return &Logger{
		Error: logger(Error, logLevel, "ERROR "),
		Info:  logger(Info, logLevel, "INFO "),
		Debug: logger(Debug, logLevel, "DEBUG "),
	}
}

func NewLoggerStr(logLevel string) *Logger {
	switch strings.ToLower(logLevel) {
	case "debug":
		return NewLogger(Debug)
	case "error":
		return NewLogger(Error)
	default:
		return NewLogger(Info)
	}
}

func logger(desired LogLevel, actual LogLevel, prefix string) *log.Logger {
	if desired <= actual {
		return log.New(os.Stdout, prefix, log.LstdFlags|log.Lmsgprefix)
	}

	return log.New(io.Discard, prefix, log.LstdFlags|log.Lmsgprefix)
}
