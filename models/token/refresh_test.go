package token

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshToken_Kind(t *testing.T) {
	assert.NotPanics(t, func() {
		tk := NewRefreshToken(uuid.New())
		cls := tk.Class()
		assert.Equal(t, Refresh, cls)
	})
}

func TestRefreshToken_HasToken(t *testing.T) {
	conn, _ := tconn.TempConn(t)
	createToken := func() *RefreshToken {
		tk := NewRefreshToken(uuid.New())
		err := conn.Create(tk).Error
		require.NoError(t, err)
		return tk
	}
	deletedToken := createToken()
	err := conn.Delete(deletedToken).Error
	require.NoError(t, err)
	tests := []struct {
		rt  *RefreshToken
		Err assert.ErrorAssertionFunc
		Has assert.BoolAssertionFunc
	}{
		{&RefreshToken{}, assert.Error, assert.False},
		{NewRefreshToken(uuid.New()), assert.NoError, assert.False},
		{createToken(), assert.NoError, assert.True},
		{deletedToken, assert.NoError, assert.False},
	}
	var has bool
	for _, test := range tests {
		has, err = test.rt.HasToken(conn)
		test.Err(t, err)
		test.Has(t, has)
	}
}
