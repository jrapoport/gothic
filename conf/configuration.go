package conf

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Configuration holds all the configuration that applies to all instances.
type Configuration struct {
	// Host is adapter to listen on.
	Host       string `json:"host" default:"localhost" `
	// RestPort is the port for the REST server to listen on.
	RestPort   int    `json:"rest_port" split_words:"true" default:"8081" `
	// RpcPort is the port for the gRPC server to listen on.
	RpcPort    int    `json:"rpc_port" split_words:"true" default:"3001" `
	// RpcWebPort is the port for the gRPC-Web server to listen on.
	RpcWebPort int    `json:"rpcweb_port" envconfig:"RPCWEB_PORT"  default:"6001" `

	SiteURL       string `json:"site_url" split_words:"true" required:"true"`
	RateLimit     string `json:"rate_limit" split_words:"true"`
	RequestID     string `json:"request_id" split_words:"true"`
	DisableSignup bool   `json:"disable_signup" split_words:"true"`

	// DB is the database configuration.
	DB         DatabaseConfig   `json:"db"`
	// External is the configuration for external OAuth providers.
	External   ExternalConfig   `json:"external"`
	// Logger is the log configuration.
	Logger     LogConfig        `json:"logger"`
	// Tracing is the Data Dog trace configuration.
	Tracing    TracingConfig    `json:"tracing"`
	// JWT is the JWT configuration.
	JWT        JWTConfig        `json:"jwt"`
	Webhook    WebhookConfig    `json:"webhook"`
	Cookies    CookieConfig     `json:"cookies"`
	Recaptcha  RecaptchaConfig  `json:"recaptcha"`
	Validation ValidationConfig `json:"validation"`

	MailConfig

	Log *logrus.Entry `json:"-" ignored:"true"`
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

	if err := ConfigureLog(config); err != nil {
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
func (c *Configuration) ApplyDefaults() {
	if c.JWT.Aud == "" {
		c.JWT.Aud = c.SiteURL
	}

	if c.Mailer.URLPaths.Invite == "" {
		c.Mailer.URLPaths.Invite = "/"
	}
	if c.Mailer.URLPaths.Confirmation == "" {
		c.Mailer.URLPaths.Confirmation = "/"
	}
	if c.Mailer.URLPaths.Recovery == "" {
		c.Mailer.URLPaths.Recovery = "/"
	}
	if c.Mailer.URLPaths.EmailChange == "" {
		c.Mailer.URLPaths.EmailChange = "/"
	}

	if c.SMTP.MaxFrequency == 0 {
		c.SMTP.MaxFrequency = 15 * time.Minute
	}

	if c.Cookies.Key == "" {
		c.Cookies.Key = "nf_jwt"
	}
	if c.Cookies.Duration == 0 {
		c.Cookies.Duration = 86400
	}
}

/*
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
		return errors.New("invalid data type for Configuration")
	}

	if len(source) == 0 {
		source = []byte("{}")
	}
	return json.Unmarshal(source, &config)
}
*/
