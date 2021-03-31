package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	lg := Log
	assert.NotNil(t, lg)
	lg = WithName("test")
	assert.NotNil(t, lg)
	file := t.TempDir() + "test.log"
	lg = UseFileOutput(file)
	Debug("test", " ", "debug")
	Debugf("%s %s", "test", "debug")
	// info
	Info("test", " ", "info")
	Infof("%s %s", "test", "info")
	// warn
	Warn("test", " ", "warn")
	Warnf("%s %s", "test", "warn")
	// error
	Error("test", " ", "error")
	Errorf("%s %s", "test", "error")
	// panic
	assert.Panics(t, func() {
		Panic("test", " ", "panic")
	})
	assert.Panics(t, func() {
		Panicf("%s %s", "test", "panic")
	})
	Log = &noop{}
	Fatal("test", " ", "error")
	Fatalf("%s %s", "test", "error")
}

type noop struct{}

var _ Logger = (*noop)(nil)

func (n noop) WithName(string) Logger { return n }

func (n noop) Debug(...interface{}) {}

func (n noop) Debugf(string, ...interface{}) {}

func (n noop) Info(...interface{}) {}

func (n noop) Infof(string, ...interface{}) {}

func (n noop) Warn(...interface{}) {}

func (n noop) Warnf(string, ...interface{}) {}

func (n noop) Error(...interface{}) {}

func (n noop) Errorf(string, ...interface{}) {}

func (n noop) Panic(...interface{}) {}

func (n noop) Panicf(string, ...interface{}) {}

func (n noop) Fatal(...interface{}) {}

func (n noop) Fatalf(string, ...interface{}) {}

func (n noop) Print(v ...interface{}) {}

func (n noop) UseFileOutput(string) Logger {
	return n
}
