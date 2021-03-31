package log

import "github.com/sirupsen/logrus"

// LogrusLogger is a logrus logger
const LogrusLogger = "logrus"

type logrusLogger struct {
	logrus.FieldLogger
	level Level
	name  string
}

var _ Logger = (*logrusLogger)(nil)

// NewLogrusLoggerWithLevel returns a new production logrus logger with the log level.
func NewLogrusLoggerWithLevel(lvl Level) Logger {
	l := logrus.New()
	level := levelToLogrusLevel(lvl)
	l.SetLevel(level)
	return &logrusLogger{
		l.WithContext(nil),
		lvl,
		"",
	}
}

// WithName returns a new named logger
func (l *logrusLogger) WithName(name string) Logger {
	if name == "" {
		return l
	}
	return named(l.FieldLogger, name, l.level)
}

// UseFileOutput uses a log file
func (l *logrusLogger) UseFileOutput(name string) Logger {
	if name == "" {
		return l
	}
	f, err := ensureLogFile(name)
	if err != nil {
		return l
	}
	log := logrus.New()
	level := levelToLogrusLevel(l.level)
	log.SetLevel(level)
	log.SetOutput(f)
	logger := &logrusLogger{
		log.WithContext(nil),
		l.level,
		l.name,
	}
	if l.name == "" {
		return logger
	}
	return logger.WithName(l.name)
}

// NOTE: for logrus panic is a higher level than fatal.
func levelToLogrusLevel(lvl Level) logrus.Level {
	switch lvl {
	case DebugLevel:
		return logrus.DebugLevel
	case InfoLevel:
		return logrus.InfoLevel
	case WarnLevel:
		return logrus.WarnLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case PanicLevel:
		return logrus.PanicLevel
	case FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}
