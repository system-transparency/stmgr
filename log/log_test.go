package log

import (
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	for _, table := range []struct {
		name string
		v    logLevel
	}{
		{
			name: "debug",
			v:    DebugLevel,
		},
		{
			name: "info",
			v:    InfoLevel,
		},
	} {
		t.Run(table.name, func(t *testing.T) {
			logger := newLogger(table.v, os.Stderr)
			if logger.verbosity != table.v {
				t.Errorf("newLogger: Expected %q, Got %q", table.v, logger.verbosity)
			}
		})
	}
}

func TestSetLoglevel(t *testing.T) {
	for _, table := range []struct {
		name string
		v    logLevel
	}{
		{
			name: "debug",
			v:    DebugLevel,
		},
		{
			name: "info",
			v:    InfoLevel,
		},
	} {
		t.Run(table.name, func(t *testing.T) {
			SetLoglevel(table.v)
			if stdLogger.verbosity != table.v {
				t.Errorf("newLogger: Expected %q, Got %q", table.v, stdLogger.verbosity)
			}
		})
	}
}

func TestPrint(t *testing.T) {
	Print("test")
	Printf("test")
}

func TestDebug(t *testing.T) {
	SetLoglevel(DebugLevel)
	Debug("test")
	Debugf("test")
}

func TestInfo(t *testing.T) {
	SetLoglevel(InfoLevel)
	Info("test")
	Infof("test")
}

func TestWarn(t *testing.T) {
	SetLoglevel(WarnLevel)
	Warn("test")
	Warnf("test")
}

func TestError(t *testing.T) {
	SetLoglevel(ErrorLevel)
	Error("test")
	Errorf("test")
}

func TestPanic(t *testing.T) {
	SetLoglevel(PanicLevel)

	defer func(t *testing.T) {
		t.Helper()

		if err := recover(); err == nil {
			t.Error("Expected a panic call")
		}
	}(t)

	Panic("test")
}

func TestPanicf(t *testing.T) {
	SetLoglevel(PanicLevel)

	defer func(t *testing.T) {
		t.Helper()

		if err := recover(); err == nil {
			t.Error("Expected a panic call")
		}
	}(t)

	Panicf("test")
}
