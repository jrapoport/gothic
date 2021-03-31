package log

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"log"
)

// A Level is a logging priority. Higher levels are more important.
// This is here as a convenience when using the various log options.
type Level int

const (
	// DebugLevel logs are typically voluminous,
	// and are usually disabled in production.
	DebugLevel Level = iota - 1

	// InfoLevel is the default logging priority.
	InfoLevel

	// WarnLevel logs are more important than Info,
	// but don't need individual human review.
	WarnLevel

	// ErrorLevel logs are high-priority. If an application runs
	// smoothly, it shouldn't generate any error-level logs.
	ErrorLevel

	// PanicLevel logs a message, then panics.
	PanicLevel

	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
)

// WithLevel returns a new copy of the logger with the level
func WithLevel(logger interface{}, level Level) Logger {
	switch logger.(type) {
	case *logrus.Logger:
		return NewLogrusLoggerWithLevel(level)
	case *logrus.Entry:
		return NewLogrusLoggerWithLevel(level)
	case *log.Logger:
		return NewStdLoggerWithLevel(level)
	case *stdLogger:
		return NewStdLoggerWithLevel(level)
	case *zap.SugaredLogger:
		return NewZapLoggerWithLevel(level)
	case *zap.Logger:
		return NewZapLoggerWithLevel(level)
	default:
		return NewStdLoggerWithLevel(level)
	}
}

// Log Levels
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
	LevelPanic = "panic"
	// LevelTrace = "trace"
)

// LevelFromString returns a normalized log level from the string
func LevelFromString(lvl string) Level {
	switch lvl {
	case LevelDebug:
		return DebugLevel
	case LevelInfo:
		return InfoLevel
	case LevelWarn:
		return WarnLevel
	case LevelError:
		return ErrorLevel
	case LevelFatal:
		return FatalLevel
	case LevelPanic:
		return PanicLevel
	}
	return PanicLevel
}
