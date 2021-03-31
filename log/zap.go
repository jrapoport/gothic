package log

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger is a std logger
const ZapLogger = "zap"

type zapLogger struct {
	*zap.SugaredLogger
	level Level
	name  string
}

var _ Logger = (*zapLogger)(nil)

// NewZapLoggerWithLevel returns a new production zap logger with the log level.
func NewZapLoggerWithLevel(lvl Level) Logger {
	level := levelToZapLevel(lvl)
	cfg := zapConfig(level, "")
	l, err := cfg.Build()
	if err != nil {
		log.Panic(err.Error())
		return nil
	}
	return &zapLogger{
		l.Sugar(),
		lvl,
		"",
	}
}

func zapConfig(lvl zapcore.Level, file string) zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(lvl)
	if file == "" {
		return cfg
	}
	f, err := ensureLogFile(file)
	if err != nil {
		return cfg
	}
	_ = f.Close()
	cfg.OutputPaths = append(cfg.OutputPaths, file)
	return cfg
}

// WithName returns a new named logger
func (z *zapLogger) WithName(name string) Logger {
	if name == "" {
		return z
	}
	return named(z.SugaredLogger, name, z.level)
}

func (z *zapLogger) Print(v ...interface{}) {
	z.Info(v...)
}

func (z *zapLogger) UseFileOutput(name string) Logger {
	level := levelToZapLevel(z.level)
	cfg := zapConfig(level, name)
	l, err := cfg.Build()
	if err != nil {
		return z
	}
	logger := &zapLogger{
		l.Sugar(),
		z.level,
		z.name,
	}
	if z.name == "" {
		return logger
	}
	return logger.WithName(z.name)
}

func levelToZapLevel(lvl Level) zapcore.Level {
	switch lvl {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case PanicLevel:
		return zapcore.PanicLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
