package tokens

import (
	"testing"

	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/store/types/provider"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrantAuthToken(t *testing.T) {
	t.Parallel()
	conn, _ := tconn.TempConn(t)
	p := provider.Google
	at, err := GrantAuthToken(conn, p, token.NoExpiration)
	assert.NoError(t, err)
	require.NotNil(t, at)
	assert.NotEmpty(t, at.AccessToken)
	assert.Equal(t, p, at.Provider)
	assert.Equal(t, p.ID(), at.UserID)
}

func TestGetAuthToken(t *testing.T) {
	t.Parallel()
	conn, _ := tconn.TempConn(t)
	p := provider.Google
	test, err := GrantAuthToken(conn, p, token.NoExpiration)
	assert.NoError(t, err)
	assert.NotNil(t, test)
	ct, err := GetAuthToken(conn, test.String())
	assert.NoError(t, err)
	assert.Equal(t, test.UserID, ct.UserID)
	assert.Equal(t, test.Token, ct.Token)
	_, err = GetAuthToken(conn, "")
	assert.Error(t, err)
}
