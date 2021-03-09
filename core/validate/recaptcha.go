package validate

import (
	"errors"
	"fmt"

	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/jrapoport/gothic/config"
)

const (
	// ReCaptchaDebugKey debug key
	ReCaptchaDebugKey = "RECAPTCHA-DEBUG-KEY"
	// ReCaptchaDebugToken debug token
	ReCaptchaDebugToken = "RECAPTCHA-DEBUG-TOKEN"
)

// ReCaptcha validates a ReCaptcha token.
func ReCaptcha(c *config.Config, ip string, token string) error {
	if c.Security.Recaptcha.Key == "" {
		return nil
	}
	if ip == "" {
		return errors.New("invalid recaptcha ip address")
	} else if token == "" {
		return errors.New("invalid recaptcha token")
	}
	if c.IsDebug() {
		if c.Security.Recaptcha.Key == ReCaptchaDebugKey && token == ReCaptchaDebugToken {
			return nil
		}
	}
	recaptcha.Init(c.Security.Recaptcha.Key)
	rc, err := recaptcha.Confirm(ip, token)
	if err != nil {
		return fmt.Errorf("recaptcha error: %w", err)
	}
	if !rc {
		return errors.New("recaptcha failed")
	}
	return nil
}
