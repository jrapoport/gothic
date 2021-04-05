package log

import (
	"io"
	"log"
	"os"
)

// StdLogger is a std logger
const StdLogger = "std"

// stdLogger is a wrapper of standard log.
type stdLogger struct {
	*log.Logger
	level Level
	name  string
}

var _ Logger = (*stdLogger)(nil)

// NewStdLoggerWithLevel is a stderr logger with the log level.
func NewStdLoggerWithLevel(lvl Level) Logger {
	return NewStdLogger(lvl, os.Stderr, "", log.LstdFlags)
}

// NewStdLogger returns a new standard logger with the log level.
func NewStdLogger(lvl Level, out io.Writer, prefix string, flag int) Logger {
	return &stdLogger{log.New(out, prefix, flag), lvl, ""}
}

// WithName returns a new named logger
func (l *stdLogger) WithName(name string) Logger {
	if name == "" {
		return l
	}
	return named(l.Logger, name, l.level)
}

// Debug logs args when the logger level is debug.
func (l *stdLogger) Debug(v ...interface{}) {
	if l.level > DebugLevel {
		return
	}
	l.Logger.Print(v...)
}

// Debugf formats args and logs the result when the logger level is debug.
func (l *stdLogger) Debugf(format string, v ...interface{}) {
	if l.level > DebugLevel {
		return
	}
	l.Logger.Printf(format, v...)
}

// Info logs args when the logger level is info.
func (l *stdLogger) Info(v ...interface{}) {
	if l.level > InfoLevel {
		return
	}
	l.Logger.Print(v...)
}

// Infof formats args and logs the result when the logger level is info.
func (l *stdLogger) Infof(format string, v ...interface{}) {
	if l.level > InfoLevel {
		return
	}
	l.Logger.Printf(format, v...)
}

// Warn logs args when the logger level is warn.
func (l *stdLogger) Warn(v ...interface{}) {
	if l.level > WarnLevel {
		return
	}
	l.Logger.Print(v...)
}

// Warnf formats args and logs the result when the logger level is warn.
func (l *stdLogger) Warnf(format string, v ...interface{}) {
	if l.level > WarnLevel {
		return
	}
	l.Logger.Printf(format, v...)
}

// Error logs args when the logger level is error.
func (l *stdLogger) Error(v ...interface{}) {
	if l.level > ErrorLevel {
		return
	}
	l.Logger.Print(v...)
}

// Errorf formats args and logs the result when the logger level is debug.
func (l *stdLogger) Errorf(format string, v ...interface{}) {
	if l.level > ErrorLevel {
		return
	}
	l.Logger.Printf(format, v...)
}

// Panic logs args on panic.
func (l *stdLogger) Panic(v ...interface{}) {
	l.Logger.Panic(v...)
}

// Panicf formats args and logs the result on panic.
func (l *stdLogger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(format, v...)
}

// Fatal logs args when the error is fatal.
func (l *stdLogger) Fatal(v ...interface{}) {
	l.Logger.Fatal(v...)
}

// Fatalf formats args and logs the result when the error is fatal.
func (l *stdLogger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatalf(format, v...)
}

// Fatal logs args when the error is fatal.
func (l *stdLogger) Print(v ...interface{}) {
	l.Info(v...)
}

// Fatalf formats args and logs the result when the error is fatal.
func (l *stdLogger) Printf(format string, v ...interface{}) {
	l.Infof(format, v...)
}

func (l *stdLogger) UseFileOutput(name string) Logger {
	f, err := ensureLogFile(name)
	if err != nil {
		return l
	}
	logger := NewStdLogger(l.level, f, "", log.LstdFlags)
	if l.name == "" {
		return logger
	}
	return logger.WithName(l.name)
}
