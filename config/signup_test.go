package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	disabled    = true
	code        = true
	reqUsername = false
	username    = false
	color       = false
	autoconfirm = true
	invites     = Users
)

func TestSignup(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		s := c.Signup
		assert.Equal(t, disabled, s.Disabled)
		assert.Equal(t, autoconfirm, s.AutoConfirm)
		assert.Equal(t, code, s.Code)
		assert.Equal(t, invites, s.Invites)
		assert.Equal(t, reqUsername, s.Username)
		assert.Equal(t, username, s.Default.Username)
		assert.Equal(t, color, s.Default.Color)
	})
}

// tests the ENV vars are correctly taking precedence
func TestSignup_Env(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clearEnv()
			loadDotEnv(t)
			c, err := loadNormalized(test.file)
			assert.NoError(t, err)
			s := c.Signup
			assert.Equal(t, disabled, s.Disabled)
			assert.Equal(t, autoconfirm, s.AutoConfirm)
			assert.Equal(t, code, s.Code)
			assert.Equal(t, invites, s.Invites)
			assert.Equal(t, reqUsername, s.Username)
			assert.Equal(t, username, s.Default.Username)
			assert.Equal(t, color, s.Default.Color)
		})
	}
}

// test the *un-normalized* defaults with load
func TestSignup_Defaults(t *testing.T) {
	clearEnv()
	c, err := load("")
	assert.NoError(t, err)
	def := signupDefaults
	s := c.Signup
	assert.Equal(t, def, s)
}
