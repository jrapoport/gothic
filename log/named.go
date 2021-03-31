package log

import (
	"log"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// logrusName matches zap
const logrusName = "logger"

// named adds a name string to the logger. How the name is added is
// logger specific i.e. a logrus field or std logger prefix, etc.
func named(logger interface{}, name string, level Level) Logger {
	switch l := logger.(type) {
	case *logrus.Logger:
		return &logrusLogger{l.WithField(logrusName, name), level, name}
	case *logrus.Entry:
		return &logrusLogger{l.WithField(logrusName, name), level, name}
	case *log.Logger:
		l.SetPrefix(name + " ")
		return &stdLogger{l, level, name}
	case *stdLogger:
		l.SetPrefix(name + " ")
		return &stdLogger{l.Logger, level, name}
	case *zap.SugaredLogger:
		return &zapLogger{l.Named(name), level, name}
	case *zap.Logger:
		return &zapLogger{l.Sugar().Named(name), level, name}
	}
	return nil
}
