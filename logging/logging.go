package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorPurple = "\033[35m"
)

const defaultCallerDepth = 2

type Logger struct {
	verbosity LogLevel
	l         log.Logger
}

func NewLogger(v LogLevel) *Logger {
	return &Logger{
		verbosity: v,
		l:         *log.New(os.Stderr, "", log.Lmsgprefix),
	}
}

func (l *Logger) Print(v ...interface{}) {
	l.l.SetFlags(0)
	l.l.Println(v...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.l.SetFlags(0)
	l.l.Printf(format, v...)
}

func (l *Logger) Debug(v ...interface{}) {
	if l.verbosity <= DebugLevel {
		l.l.SetPrefix(colorPurple + "DEBUG: " + colorReset + getCaller())
		l.l.Println(v...)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.verbosity <= DebugLevel {
		l.l.SetPrefix(colorPurple + "DEBUG: " + colorReset + getCaller())
		l.l.Printf(format, v...)
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.verbosity <= InfoLevel {
		l.l.SetPrefix(colorGreen + "INFO: " + colorReset)
		l.l.Println(v...)
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.verbosity <= InfoLevel {
		l.l.SetPrefix(colorGreen + "INFO: " + colorReset)
		l.l.Printf(format, v...)
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if l.verbosity <= WarnLevel {
		l.l.SetPrefix(colorYellow + "WARN: " + colorReset)
		l.l.Println(v...)
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.verbosity <= WarnLevel {
		l.l.SetPrefix(colorYellow + "WARN: " + colorReset)
		l.l.Printf(format, v...)
	}
}

func (l *Logger) Error(v ...interface{}) {
	if l.verbosity <= ErrorLevel {
		l.l.SetPrefix(colorRed + "ERROR: " + colorReset + getCaller())
		l.l.SetFlags(log.Lmsgprefix | log.Ltime)
		l.l.Println(v...)
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.verbosity <= ErrorLevel {
		l.l.SetPrefix(colorRed + "ERROR: " + colorReset + getCaller())
		l.l.SetFlags(log.Lmsgprefix | log.Ltime)
		l.l.Printf(format, v...)
	}
}

func (l *Logger) Panic(v ...interface{}) {
	if l.verbosity <= PanicLevel {
		l.l.SetPrefix(colorRed + "PANIC: " + colorReset + getCaller())
		l.l.SetFlags(log.Lmsgprefix | log.Ltime)
		l.l.Panicln(v...)
	}
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	if l.verbosity <= PanicLevel {
		l.l.SetPrefix(colorRed + "PANIC: " + colorReset + getCaller())
		l.l.SetFlags(log.Lmsgprefix | log.Ltime)
		l.l.Panicf(format, v...)
	}
}

func getCaller() string {
	_, file, line, ok := runtime.Caller(defaultCallerDepth)
	if !ok {
		file = "???"
		line = 0
	}

	return fmt.Sprintf("%s:%d: ", file, line)
}
