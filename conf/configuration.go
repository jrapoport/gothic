package conf

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Configuration holds all the configuration that applies to all instances.
type Configuration struct {
	Host       string `json:"host" default:"localhost" `
	RestPort   int    `json:"rest_port" split_words:"true" default:"8081" `
	RpcPort    int    `json:"rpc_port" split_words:"true" default:"3001" `
	RpcWebPort int    `json:"rpcweb_port" envconfig:"RPCWEB_PORT"  default:"6001" `

	SiteURL       string `json:"site_url" split_words:"true" required:"true"`
	RateLimit     string `json:"rate_limit" split_words:"true"`
	RequestID     string `json:"request_id" split_words:"true"`
	DisableSignup bool   `json:"disable_signup" split_words:"true"`
	PasswordRegex string `json:"_" split_words:"true" default:"^[a-zA-Z0-9[:punct:]]{8,28}$"`

	DB        DatabaseConfig  `json:"db"`
	External  ExternalConfig  `json:"external"`
	Log       LogConfig       `json:"log"`
	Tracing   TracingConfig   `json:"tracing"`
	JWT       JWTConfig       `json:"jwt"`
	Webhook   WebhookConfig   `json:"webhook"`
	Cookies   CookieConfig    `json:"cookies"`
	Recaptcha RecaptchaConfig `json:"recaptcha"`

	MailConfig
}

func loadEnvironment(filename string) error {
	var err error
	if filename != "" {
		err = godotenv.Load(filename)
	} else {
		err = godotenv.Load()
		// handle if .env file does not exist, this is OK
		if os.IsNotExist(err) {
			return nil
		}
	}
	return err
}

// LoadConfiguration loads configuration from file and environment variables.
func LoadConfiguration(filename string) (*Configuration, error) {
	if err := loadEnvironment(filename); err != nil {
		return nil, err
	}

	config := new(Configuration)
	if err := envconfig.Process("gothic", config); err != nil {
		return nil, err
	}

	if _, err := ConfigureLog(&config.Log); err != nil {
		return nil, err
	}

	ConfigureTracing(&config.Tracing)

	if config.SMTP.MaxFrequency == 0 {
		config.SMTP.MaxFrequency = 15 * time.Minute
	}

	config.ApplyDefaults()

	return config, nil
}

// ApplyDefaults sets defaults for a Configuration
func (config *Configuration) ApplyDefaults() {
	if config.JWT.AdminGroup == "" {
		config.JWT.AdminGroup = "admin"
	}

	if config.JWT.Exp == 0 {
		config.JWT.Exp = 3600
	}

	if config.Mailer.URLPaths.Invite == "" {
		config.Mailer.URLPaths.Invite = "/"
	}
	if config.Mailer.URLPaths.Confirmation == "" {
		config.Mailer.URLPaths.Confirmation = "/"
	}
	if config.Mailer.URLPaths.Recovery == "" {
		config.Mailer.URLPaths.Recovery = "/"
	}
	if config.Mailer.URLPaths.EmailChange == "" {
		config.Mailer.URLPaths.EmailChange = "/"
	}

	if config.SMTP.MaxFrequency == 0 {
		config.SMTP.MaxFrequency = 15 * time.Minute
	}

	if config.Cookies.Key == "" {
		config.Cookies.Key = "nf_jwt"
	}

	if config.Cookies.Duration == 0 {
		config.Cookies.Duration = 86400
	}
}

func (config *Configuration) Value() (driver.Value, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return driver.Value(""), err
	}
	return driver.Value(string(data)), nil
}

func (config *Configuration) Scan(src interface{}) error {
	var source []byte
	switch v := src.(type) {
	case string:
		source = []byte(v)
	case []byte:
		source = v
	default:
		return errors.New("Invalid data type for Configuration")
	}

	if len(source) == 0 {
		source = []byte("{}")
	}
	return json.Unmarshal(source, &config)
}
