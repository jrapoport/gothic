package config

import (
	"time"

	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/store/drivers"
)

const (
	serviceName        = "gothic"
	cookieDuration     = 24 * 60 * time.Minute
	dbDriver           = drivers.MySQL
	dbMaxRetry         = 3
	jwtAlgorithm       = "HS256"
	jwtExpiration      = 60 * time.Minute
	logLevel           = "info"
	logTimeFormat      = time.RFC3339Nano
	mailFrom           = ":name <do-not-reply@:link_hostname>"
	mailTheme          = "default"
	passwordRegex      = "^[a-zA-Z0-9[:punct:]]{8,40}$"
	secRateLimit       = 5 * time.Minute
	smtpAuthentication = "plain"
	smtpEncryption     = "none"
	smtpExpiration     = 60 * time.Minute
	smtpPort           = 587
	smtpSendLimit      = 1 * time.Minute
	usernameRegex      = "^[a-zA-Z0-9_]{2,255}$"
	webhookMaxRetry    = 3
	webhookTimeout     = 30 * time.Second
)

var configDefaults = Config{
	Service:       serviceDefaults,
	Network:       networkDefaults,
	Security:      securityDefaults,
	Authorization: authorizationDefaults,
	DB:            databaseDefaults,
	Mail:          mailDefaults,
	Signup:        signupDefaults,
	Webhook:       webhooksDefaults,
	Logger:        loggerDefaults,
}

var serviceDefaults = Service{
	Name: serviceName,
}

var networkDefaults = Network{
	Host:   "localhost",
	REST:   "localhost:8081",
	RPC:    "localhost:3001",
	RPCWeb: "localhost:6001",
	Health: "localhost:10001",
}

var securityDefaults = Security{
	MaskEmails: true,
	RateLimit:  secRateLimit,
	JWT:        jwtDefaults,
	Recaptcha: Recaptcha{
		Login: true,
	},
	Validation: Validation{
		UsernameRegex: usernameRegex,
		PasswordRegex: passwordRegex,
	},
	Cookies: Cookies{
		Duration: cookieDuration,
	},
}

var jwtDefaults = JWT{
	Algorithm:  jwtAlgorithm,
	Expiration: jwtExpiration,
}

var authorizationDefaults = Authorization{
	UseInternal: true,
	Providers:   defaultProviders(),
}

func defaultProviders() Providers {
	list := Providers{}
	for name := range provider.External {
		list[name] = Provider{}
	}
	return list
}

var databaseDefaults = Database{
	Driver:     dbDriver,
	MaxRetries: dbMaxRetry,
}

var mailDefaults = Mail{
	SMTP: SMTP{
		Port:           smtpPort,
		Authentication: smtpAuthentication,
		Encryption:     smtpEncryption,
		KeepAlive:      true,
		Expiration:     smtpExpiration,
		SendLimit:      smtpSendLimit,
		SpamProtection: true,
	},
	From:          mailFrom,
	Theme:         mailTheme,
	MailTemplates: templateDefaults,
}

var signupDefaults = Signup{
	Invites:  Admins,
	Username: true,
	Default: SignupDefaults{
		Username: true,
		Color:    true,
	},
}

var webhooksDefaults = Webhooks{
	MaxRetries: webhookMaxRetry,
	Timeout:    webhookTimeout,
	JWT:        jwtDefaults,
}

var loggerDefaults = Logger{
	Level:     logLevel,
	Colors:    true,
	Timestamp: logTimeFormat,
}
