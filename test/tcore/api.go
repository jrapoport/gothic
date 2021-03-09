package tcore

import (
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/require"
)

// API for testing.
func API(t *testing.T, smtp bool) (*core.API, *config.Config, *tconf.SMTPMock) {
	c := tconf.TempDB(t)
	c.Mail.SpamProtection = false
	c.Signup.AutoConfirm = true
	c.Signup.Default.Username = false
	c.Signup.Default.Color = false
	c.Validation.UsernameRegex = ""
	c.Signup.Username = false
	var mock *tconf.SMTPMock
	if smtp {
		c, mock = tconf.MockSMTP(t, c)
		require.NotNil(t, mock)
	}
	a, err := core.NewAPI(c)
	require.NoError(t, err)
	require.NotNil(t, a)
	t.Cleanup(func() {
		err = a.Shutdown()
		require.NoError(t, err)
	})
	return a, c, mock
}
