package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger config
type Logger struct {
	Level     string   `json:"level"`
	File      string   `json:"file"`
	Colors    bool     `json:"colors"`
	Timestamp string   `json:"timestamp"`
	Fields    []string `json:"fields"`

	// Tracer is the Data Dog trace configuration.
	Tracer Tracer `json:"tracer"`

	// log is a configured instance of a log based on the Logger settings.
	log logrus.FieldLogger
}

func (l *Logger) load(srv Service) error {
	log, err := l.NewLogger()
	if err != nil {
		return err
	}
	l.log = log.WithField("logger", srv.Name)
	return l.Tracer.StartTracer(srv.Name, srv.Version())
}

// NewLogger returns a new instance of the configured logger.
func (l *Logger) NewLogger() (logrus.FieldLogger, error) {
	log := logrus.New()
	// always use the full timestamp
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		DisableTimestamp: false,
		TimestampFormat:  l.Timestamp,
		DisableColors:    !l.Colors,
		QuoteEmptyFields: true,
	})
	// use a file if you want
	if l.File != "" {
		f, err := os.OpenFile(l.File, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		log.SetOutput(f)
		log.Infof("Set output file to %s", l.File)
	}
	if l.Level != "" {
		level, err := logrus.ParseLevel(l.Level)
		if err != nil {
			return nil, err
		}
		log.SetLevel(level)
		log.Debug("set log level to: " + log.GetLevel().String())
	}
	f := logrus.Fields{}
	fields := newKeyValueMap(l.Fields)
	for k, v := range fields {
		f[k] = v
	}
	return log.WithFields(f), nil
}
