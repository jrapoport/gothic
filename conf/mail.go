package conf

import "time"

type MailConfig struct {
	SMTP   SMTPConfig   `json:"smtp"`
	Mailer MailerConfig `json:"mailer"`
}

type SMTPConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port,omitempty" default:"587"`
	Username     string        `json:"username"`
	Password     string        `json:"-"`
	MaxFrequency time.Duration `json:"max_frequency" split_words:"true"`
	AdminEmail   string        `json:"admin_email" split_words:"true"`
}

type MailerConfig struct {
	Autoconfirm bool        `json:"autoconfirm"`
	Subjects    EmailConfig `json:"subjects"`
	Templates   EmailConfig `json:"templates"`
	URLPaths    EmailConfig `json:"url_paths"`
}

// EmailConfig holds the configuration for emails, both subjects and template URLs.
type EmailConfig struct {
	Invite       string `json:"invite"`
	Confirmation string `json:"confirmation"`
	Recovery     string `json:"recovery"`
	EmailChange  string `json:"email_change" split_words:"true"`
}
