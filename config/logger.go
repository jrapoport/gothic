package config

import (
	"github.com/jrapoport/gothic/log"
)

// Logger config
type Logger struct {
	Package   string   `json:"package"`
	Level     string   `json:"level"`
	File      string   `json:"file"`
	Colors    bool     `json:"colors"`
	Timestamp string   `json:"timestamp"`
	Fields    []string `json:"fields"`

	// Tracer is the Data Dog trace configuration.
	Tracer Tracer `json:"tracer"`

	// log is a configured instance of a log based on the Logger settings.
	log log.Logger
}

func (c *Logger) load(srv Service) error {
	lg := c.NewLogger()
	c.log = lg.WithName(srv.Name)
	return c.Tracer.StartTracer(srv.Name, srv.Version())
}

// NewLogger returns a new instance of the configured logger.
func (c *Logger) NewLogger() log.Logger {
	lvl := log.LevelFromString(c.Level)
	var l log.Logger
	switch c.Package {
	case log.LogrusLogger:
		l = log.NewLogrusLoggerWithLevel(lvl)
	case log.ZapLogger:
		l = log.NewZapLoggerWithLevel(lvl)
	case log.StdLogger:
		fallthrough
	default:
		l = log.NewStdLoggerWithLevel(lvl)
	}
	// use a file if you want
	if c.File != "" {
		l = l.UseFileOutput(c.File)
		l.Infof("Set output file to %s", c.File)
	}
	return l
}

/*
// NewLogger returns a new instance of the configured logger.
func (l *Logger) NewLogger() (log.Logger, error) {
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
*/
