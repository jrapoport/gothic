package config

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/segmentio/encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config holds all the configuration that applies to all instances.
type Config struct {
	Service `yaml:",inline" mapstructure:",squash"`

	// Network are the network settings
	Network `yaml:",inline" mapstructure:",squash"`

	// Security
	Security `yaml:",inline" mapstructure:",squash"`

	// Authorization is the configuration for external auth providers.
	Authorization `yaml:",inline" mapstructure:",squash"`

	// Database is the database configuration.
	DB Database `json:"db"`

	// Mail is the mailer configuration.
	Mail Mail `json:"mail"`

	// Signup is the signup configuration.
	Signup Signup `json:"signup"`

	// Webhook is the configuration for webhooks
	Webhook Webhooks `json:"webhook"`

	// Logging is the log configuration.
	Logger `json:"log"  yaml:"log" mapstructure:"log"`
}

// LoadConfig loads a config from the provided path. If the
// path is empty it will still fallback to the enn vars.
func LoadConfig(filename string) (*Config, error) {
	var err error
	c, err := loadNormalized(filename)
	if err != nil {
		return nil, err
	}
	// do this first
	err = c.Logger.load(c.Service)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Log returns a configured log or a default log if one is not found.
func (c *Config) Log() logrus.FieldLogger {
	if c.log == nil {
		c.log = logrus.New()
	}
	return c.log
}

// Write writes a copy of the config to the supplied path.
func (c *Config) Write(path string) error {
	v := viper.New()
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	const envHackFormat = "yaml"
	v.SetConfigType(envHackFormat)
	err = v.MergeConfig(bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	return v.WriteConfigAs(path)
}

func (c *Config) normalize() error {
	err := c.Service.normalize()
	if err != nil {
		return err
	}
	err = c.Network.normalize(c.Service)
	if err != nil {
		return err
	}
	err = c.Security.normalize(c.Service)
	if err != nil {
		return err
	}
	err = c.Authorization.normalize(c.Service, c.REST)
	if err != nil {
		return err
	}
	err = c.DB.normalize(c.Service)
	if err != nil {
		return err
	}
	err = c.Mail.normalize(c.Service)
	if err != nil {
		return err
	}
	// do this last
	return c.Webhook.normalize(c.Service, c.JWT)
}

// loadFromFile loads a config file, the format will be inferred from the ext
// and with correct precedence (i.e. env vars will overwrite defaults etc.)
// If the path is empty, no error is returned and the env vars are used.
func loadFromFile(path string) (*viper.Viper, error) {
	v, err := loadEnv()
	if err != nil {
		return nil, err
	}
	if path != "" {
		// viper does not read "dotenv" or "env" files correctly &
		// uses '.' as the nested key delim vs. '_' as expected.
		// see: https://github.com/spf13/viper/issues/814
		ext := strings.ToLower(filepath.Ext(path))
		if len(ext) > 1 {
			ext = ext[1:]
		}
		switch ext {
		case "dotenv", "env":
			err = godotenv.Load(path)
			if err != nil {
				return nil, err
			}
			break
		case "json":
			// viper can't seem to figure out a
			// float is an int unless we do this.
			var b []byte
			b, err = ioutil.ReadFile(path)
			if err != nil {
				return nil, err
			}
			m := map[string]interface{}{}
			err = json.Unmarshal(b, &m)
			if err != nil {
				return nil, err
			}
			err = v.MergeConfigMap(m)
			if err != nil {
				return nil, err
			}
		default:
			v.SetConfigFile(path)
		}
	} else {
		v.AddConfigPath("/etc/gothic/")
		v.AddConfigPath("$HOME/.gothic")
		v.AddConfigPath(".")
		v.SetConfigName("gothic")
	}
	err = v.MergeInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return nil, err
		}
	}
	return v, nil
}

func load(path string) (*Config, error) {
	v, err := loadFromFile(path)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = v.Unmarshal(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// loadNormalized loads a config and normalizes the result.
func loadNormalized(path string) (*Config, error) {
	c, err := load(path)
	if err != nil {
		return nil, err
	}
	err = c.normalize()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func newKeyValueMap(s []string) map[string]string {
	m := make(map[string]string, len(s))
	for _, kv := range s {
		parts := strings.Split(kv, "=")
		if len(parts) < 2 {
			continue
		}
		m[parts[0]] = parts[1]
	}
	return m
}

// ErrRateLimitExceeded rate limited exceeded error.
var ErrRateLimitExceeded = errors.New("rate limit exceeded")
