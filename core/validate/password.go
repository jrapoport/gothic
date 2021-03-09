package validate

import (
	"errors"
	"regexp"

	"github.com/jrapoport/gothic/config"
)

// Password validates a password
func Password(c *config.Config, password string) error {
	if c.Security.Validation.PasswordRegex == "" {
		return nil
	} else if password == "" {
		return errors.New("password required")
	}
	rx, err := regexp.Compile(c.Security.Validation.PasswordRegex)
	if err != nil {
		return err
	}
	if !rx.MatchString(password) {
		return errors.New("invalid password")
	}
	return nil
}
