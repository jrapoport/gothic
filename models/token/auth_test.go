package token

import (
	"testing"

	"github.com/jrapoport/gothic/store/types/provider"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthToken_Kind(t *testing.T) {
	t.Parallel()
	assert.NotPanics(t, func() {
		tk := NewAuthToken(provider.Google, 0)
		cls := tk.Class()
		assert.Equal(t, Auth, cls)
	})
}

func TestAuthToken_HasToken(t *testing.T) {
	t.Parallel()
	conn, _ := tconn.TempConn(t)
	createToken := func() *AuthToken {
		tk := NewAuthToken(provider.Google, 0)
		err := conn.Create(tk).Error
		require.NoError(t, err)
		return tk
	}
	deletedToken := createToken()
	err := conn.Delete(deletedToken).Error
	require.NoError(t, err)
	tests := []struct {
		ct  *AuthToken
		Err assert.ErrorAssertionFunc
		Has assert.BoolAssertionFunc
	}{
		{&AuthToken{}, assert.Error, assert.False},
		{NewAuthToken(provider.Google, 0), assert.NoError, assert.False},
		{createToken(), assert.NoError, assert.True},
		{deletedToken, assert.NoError, assert.False},
	}
	var has bool
	for _, test := range tests {
		has, err = test.ct.HasToken(conn)
		test.Err(t, err)
		test.Has(t, has)
	}
}
