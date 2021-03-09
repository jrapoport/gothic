package validate

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/jrapoport/gothic/config"
)

// Username validates a username.
func Username(c *config.Config, username string) error {
	if c.Security.Validation.UsernameRegex == "" {
		return nil
	} else if username == "" {
		return errors.New("invalid username")
	}
	rx, err := regexp.Compile(c.Security.Validation.UsernameRegex)
	if err != nil {
		return err
	}
	if !rx.MatchString(username) {
		return fmt.Errorf("invalid username: %s", username)
	}
	return nil
}
