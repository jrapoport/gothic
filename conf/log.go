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

func ConfigureLog(config *LogConfig) (*logrus.Entry, error) {
	logger := logrus.New()

	tsFormat := time.RFC3339Nano
	if config.TimestampFormat != "" {
		tsFormat = config.TimestampFormat
	}
	// always use the full timestamp
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		DisableTimestamp: false,
		TimestampFormat:  tsFormat,
		DisableColors:    config.DisableColors,
		QuoteEmptyFields: config.QuoteEmpty,
	})

	// use a file if you want
	if config.File != "" {
		f, errOpen := os.OpenFile(config.File, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
		if errOpen != nil {
			return nil, errOpen
		}
		logger.SetOutput(f)
		logger.Infof("Set output file to %s", config.File)
	}

	if config.Level != "" {
		level, err := logrus.ParseLevel(config.Level)
		if err != nil {
			return nil, err
		}
		logger.SetLevel(level)
		logger.Debug("Set log level to: " + logger.GetLevel().String())
	}

	f := logrus.Fields{}
	for k, v := range config.Fields {
		f[k] = v
	}

	return logger.WithFields(f), nil
}
