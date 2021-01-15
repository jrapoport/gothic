package conf

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type LogConfig struct {
	Level           string                 `json:"log_level" mapstructure:"log_level"`
	File            string                 `json:"log_file" mapstructure:"log_file"`
	DisableColors   bool                   `json:"disable_colors" mapstructure:"disable_colors" split_words:"true"`
	QuoteEmpty      bool                   `json:"quote_empty" mapstructure:"quote_empty_fields" split_words:"true"`
	TimestampFormat string                 `json:"timestamp_format" mapstructure:"ts_format" split_words:"true"`
	Fields          map[string]interface{} `json:"fields" mapstructure:"fields"`
}

func ConfigureLog(c *Configuration) error {
	lc := c.Logger
	logger := logrus.New()

	tsFormat := time.RFC3339Nano
	if lc.TimestampFormat != "" {
		tsFormat = lc.TimestampFormat
	}
	// always use the full timestamp
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		DisableTimestamp: false,
		TimestampFormat:  tsFormat,
		DisableColors:    lc.DisableColors,
		QuoteEmptyFields: lc.QuoteEmpty,
	})

	// use a file if you want
	if lc.File != "" {
		f, err := os.OpenFile(lc.File, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
		if err != nil {
			return err
		}
		logger.SetOutput(f)
		logger.Infof("Set output file to %s", lc.File)
	}

	if lc.Level != "" {
		level, err := logrus.ParseLevel(lc.Level)
		if err != nil {
			return err
		}
		logger.SetLevel(level)
		logger.Debug("Set log level to: " + logger.GetLevel().String())
	}

	f := logrus.Fields{}
	for k, v := range lc.Fields {
		f[k] = v
	}
	e := logger.WithFields(f)
	c.Log = e
	return nil
}
