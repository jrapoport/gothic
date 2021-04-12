package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/imdario/mergo"
)

// Security config
type Security struct {
	// RootPassword the password for the super admin user. (there is no interactive login).
	RootPassword string `json:"root_password" yaml:"root_password" mapstructure:"root_password"`
	// MaskEmails user emails returned by api calls are masked by default.
	MaskEmails bool `json:"mask_emails" yaml:"mask_emails" mapstructure:"mask_emails"`
	// RateLimit is the rate limit per 100 requests to be enforced by ip address
	RateLimit time.Duration `json:"rate_limit" yaml:"rate_limit" mapstructure:"rate_limit"`
	// JWT is the JWT configuration.
	JWT JWT `json:"jwt"`
	// Recaptcha is the google CAPTCHA configuration.
	Recaptcha Recaptcha `json:"recaptcha"`
	// Validation is the validation to apply to user submitted data.
	Validation Validation `json:"validation"`
	// Cookies is the configuration for cookies
	Cookies Cookies `json:"cookies"`
}

func (s *Security) normalize(srv Service) error {
	s.JWT.normalize(srv, jwtDefaults)
	userRx := s.Validation.UsernameRegex
	if userRx != "" {
		_, err := regexp.Compile(userRx)
		if err != nil {
			err = fmt.Errorf("invalid username regex %s: %s", userRx, err)
			return err
		}
	}
	passRx := s.Validation.PasswordRegex
	if passRx != "" {
		_, err := regexp.Compile(passRx)
		if err != nil {
			err = fmt.Errorf("invalid password regex %s: %s", passRx, err)
			return err
		}
	}
	if s.Cookies.Duration == 0 {
		s.Cookies.Duration = cookieDuration
	}
	return nil
}

func (s *Security) CheckRequired() error {
	if s.RootPassword == "" {
		return errors.New("root password is required")
	}
	return nil
}

// JWT holds all the JWT related configuration.
type JWT struct {
	Secret    string `json:"secret"`
	Algorithm string `json:"algorithm"`
	// Issuer is the the entity that issued the token (default: Config.Service)
	Issuer string `json:"issuer"`
	// Audience is an optional comma separated list of resource
	// servers that should accept the token (default: n/a)
	Audience   string        `json:"audience"`
	Expiration time.Duration `json:"expiration"`
}

func (j *JWT) normalize(srv Service, def JWT) {
	if def.Issuer == "" {
		def.Issuer = strings.ToLower(srv.Name)
	}
	// no error is possible here since we
	// control the struct entirely
	_ = mergo.Merge(j, def)
}

func (j *JWT) CheckRequired() error {
	if j.Secret == "" {
		return errors.New("jwt secret is required")
	}
	return nil
}

// Recaptcha config
type Recaptcha struct {
	Key   string `json:"key"`
	Login bool   `json:"login"`
}

// Validation config
type Validation struct {
	UsernameRegex string `json:"username_regex" yaml:"username_regex" mapstructure:"username_regex"`
	PasswordRegex string `json:"password_regex" yaml:"password_regex" mapstructure:"password_regex"`
}

// Cookies config
type Cookies struct {
	Duration time.Duration `json:"duration"`
}
