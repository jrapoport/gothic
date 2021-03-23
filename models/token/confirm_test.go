package token

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfirmToken_Kind(t *testing.T) {
	t.Parallel()
	assert.NotPanics(t, func() {
		tk := NewConfirmToken(uuid.New(), 0)
		cls := tk.Class()
		assert.Equal(t, Confirm, cls)
	})
}

func TestConfirmToken_HasToken(t *testing.T) {
	t.Parallel()
	conn, _ := tconn.TempConn(t)
	createToken := func() *ConfirmToken {
		tk := NewConfirmToken(uuid.New(), 0)
		assert.False(t, tk.Usable())
		err := conn.Create(tk).Error
		require.NoError(t, err)
		assert.True(t, tk.Usable())
		return tk
	}
	deletedToken := createToken()
	err := conn.Delete(deletedToken).Error
	require.NoError(t, err)
	tests := []struct {
		ct  *ConfirmToken
		Err assert.ErrorAssertionFunc
		Has assert.BoolAssertionFunc
	}{
		{&ConfirmToken{}, assert.Error, assert.False},
		{NewConfirmToken(uuid.New(), 0), assert.NoError, assert.False},
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
