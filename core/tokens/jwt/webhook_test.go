package jwt

import (
	"testing"

	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWebhookClaims(t *testing.T) {
	c := tconf.Config(t)
	var checksum = utils.SecureToken()
	claims := NewWebhookClaims(checksum)
	assert.Equal(t, checksum, claims.Checksum())
	tok := NewToken(c.JWT, claims)
	bear, err := tok.Bearer()
	require.NoError(t, err)
	var web WebhookClaims
	err = ParseClaims(c.JWT, bear, &web)
	assert.NoError(t, err)
	assert.Equal(t, checksum, claims.Checksum())
}
