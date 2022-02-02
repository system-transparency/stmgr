package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

type logLevel int

const (
	DebugLevel logLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
)

// This might not work in every terminal though...
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorPurple = "\033[35m"
)

// This value is taken from the stdlib log pkg.
const defaultCallerDepth = 2

type logger struct {
	verbosity logLevel
	l         log.Logger
}

var stdLogger = newLogger(ErrorLevel, os.Stdout) //nolint:gochecknoglobals

func newLogger(v logLevel, w io.Writer) *logger {
	return &logger{
		verbosity: v,
		l:         *log.New(w, "", log.Lmsgprefix),
	}
}

// SetLoglevel changes verbosity of the logger.
// Allowed levels are the following and setting on
// enables all lower levels as well:
// Debug (Highest)
// Info
// Warn
// Error
// Panic (Lowest)
//
// The default level is "Error".
func SetLoglevel(v logLevel) {
	stdLogger.verbosity = v
}

// Print prints to stdout.
func Print(v ...interface{}) {
	stdLogger.l.SetFlags(0)
	stdLogger.l.Println(v...)
}

// Printf is the same as Print but with formatting.
func Printf(format string, v ...interface{}) {
	stdLogger.l.SetFlags(0)
	stdLogger.l.Printf(format, v...)
}

// Debug prints debug level output.
func Debug(v ...interface{}) {
	if stdLogger.verbosity <= DebugLevel {
		stdLogger.l.SetPrefix(colorPurple + "DEBUG: " + colorReset + getCaller())
		stdLogger.l.Println(v...)
	}
}

// Debugf is the same as Debug but with formatting.
func Debugf(format string, v ...interface{}) {
	if stdLogger.verbosity <= DebugLevel {
		stdLogger.l.SetPrefix(colorPurple + "DEBUG: " + colorReset + getCaller())
		stdLogger.l.Printf(format, v...)
	}
}

// Info prints info level output.
func Info(v ...interface{}) {
	if stdLogger.verbosity <= InfoLevel {
		stdLogger.l.SetPrefix(colorGreen + "INFO: " + colorReset)
		stdLogger.l.Println(v...)
	}
}

// Infof is the same as Info but with formatting.
func Infof(format string, v ...interface{}) {
	if stdLogger.verbosity <= InfoLevel {
		stdLogger.l.SetPrefix(colorGreen + "INFO: " + colorReset)
		stdLogger.l.Printf(format, v...)
	}
}

// Warn prints warn level output.
func Warn(v ...interface{}) {
	if stdLogger.verbosity <= WarnLevel {
		stdLogger.l.SetPrefix(colorYellow + "WARN: " + colorReset)
		stdLogger.l.Println(v...)
	}
}

// Warnf is the same as Warn but with formatting.
func Warnf(format string, v ...interface{}) {
	if stdLogger.verbosity <= WarnLevel {
		stdLogger.l.SetPrefix(colorYellow + "WARN: " + colorReset)
		stdLogger.l.Printf(format, v...)
	}
}

// Error prints error level output.
func Error(v ...interface{}) {
	if stdLogger.verbosity <= ErrorLevel {
		stdLogger.l.SetPrefix(colorRed + "ERROR: " + colorReset + getCaller())
		stdLogger.l.SetFlags(log.Lmsgprefix | log.Ltime)
		stdLogger.l.Println(v...)
	}
}

// Errorf is the same as Error but with formatting.
func Errorf(format string, v ...interface{}) {
	if stdLogger.verbosity <= ErrorLevel {
		stdLogger.l.SetPrefix(colorRed + "ERROR: " + colorReset + getCaller())
		stdLogger.l.SetFlags(log.Lmsgprefix | log.Ltime)
		stdLogger.l.Printf(format, v...)
	}
}

// Panic prints panic level output and calls panic().
func Panic(v ...interface{}) {
	if stdLogger.verbosity <= PanicLevel {
		stdLogger.l.SetPrefix(colorRed + "PANIC: " + colorReset + getCaller())
		stdLogger.l.SetFlags(log.Lmsgprefix | log.Ltime)
		stdLogger.l.Panicln(v...)
	}
}

// Panicf is the same as Panic but with formatting.
func Panicf(format string, v ...interface{}) {
	if stdLogger.verbosity <= PanicLevel {
		stdLogger.l.SetPrefix(colorRed + "PANIC: " + colorReset + getCaller())
		stdLogger.l.SetFlags(log.Lmsgprefix | log.Ltime)
		stdLogger.l.Panicf(format, v...)
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
