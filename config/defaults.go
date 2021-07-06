package config

import (
	"github.com/jrapoport/gothic/core/events"
	"time"

	"github.com/jrapoport/gothic/log"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/store/drivers"
	"github.com/jrapoport/gothic/utils"
)

const (
	serviceName        = "gothic"
	cookieDuration     = 24 * 60 * time.Minute
	dbDriver           = drivers.MySQL
	dbMaxRetry         = 3
	jwtAlgorithm       = "HS256"
	jwtExpiration      = 60 * time.Minute
	logLevel           = log.LevelPanic
	logTimeFormat      = time.RFC3339Nano
	mailFrom           = ":name <do-not-reply@:link_hostname>"
	mailTheme          = "default"
	usernameRegex      = "^[a-zA-Z0-9_]{2,255}$"
	passwordRegex      = "^[a-zA-Z0-9[:punct:]]{8,40}$"
	secRateLimit       = 5 * time.Minute
	smtpAuthentication = "plain"
	smtpEncryption     = "none"
	smtpExpiration     = 60 * time.Minute
	smtpKeepalive      = true
	smtpPort           = 25
	smtpSendLimit      = 1 * time.Minute
	smtpSpamProtection = true
	webhookMaxRetry    = 3
	webhookTimeout     = 30 * time.Second
	host               = "localhost"
)

var configDefaults = &Config{
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
	HealthAddress: utils.MakeAddress(host, 7720),
	RPCAddress:    utils.MakeAddress(host, 7721),
	AdminAddress:  utils.MakeAddress(host, 7722),
	RESTAddress:   utils.MakeAddress(host, 7727),
	RPCWebAddress: utils.MakeAddress(host, 7729),
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
		list[name] = Provider{
			Scopes: []string{},
		}
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
		KeepAlive:      smtpKeepalive,
		Expiration:     smtpExpiration,
		SendLimit:      smtpSendLimit,
		SpamProtection: smtpSpamProtection,
	},
	MailFormat: MailFormat{
		From: mailFrom,
	},
	MailTemplates: templateDefaults,
}

var signupDefaults = Signup{
	Invites: Admins,
}

var webhooksDefaults = Webhooks{
	MaxRetries: webhookMaxRetry,
	Timeout:    webhookTimeout,
	JWT:        jwtDefaults,
	Events:     []events.Event{},
}

var loggerDefaults = Logger{
	Level:     logLevel,
	Timestamp: logTimeFormat,
	Fields:    []string{},
	Tracer: Tracer{
		Tags: []string{},
	},
}
