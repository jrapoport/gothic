package validate

import (
	"testing"

	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/assert"
)

func TestReCaptcha(t *testing.T) {
	t.Parallel()
	c := tconf.Config(t)
	key := ReCaptchaDebugKey
	token := ReCaptchaDebugToken
	ip := "0.0.0.0"
	tests := []struct {
		key   string
		ip    string
		token string
		Err   assert.ErrorAssertionFunc
	}{
		{"", "", "", assert.NoError},
		{"", ip, "", assert.NoError},
		{"", "", token, assert.NoError},
		{"", ip, token, assert.NoError},
		{key, "", "", assert.Error},
		{key, ip, "", assert.Error},
		{key, "", token, assert.Error},
		{key, ip, token, assert.NoError},
	}
	for _, test := range tests {
		c.Recaptcha.Key = test.key
		err := ReCaptcha(c, test.ip, test.token)
		test.Err(t, err, "key: %s ip: %s token: %s", test.key, test.ip, test.token)
	}
}
