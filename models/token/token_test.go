package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testToken    = "1234567890asdfghjkl="
	tokenExp     = time.Hour * 1
	tokenExpired = tokenExp * -1
	tokenTests   = []struct {
		uses int
		exp  time.Duration
		use  Usage
		max  int
	}{
		{-2, NoExpiration, Infinite, InfiniteUse},
		{InfiniteUse, NoExpiration, Infinite, InfiniteUse},
		{0, NoExpiration, Single, SingleUse},
		{SingleUse, NoExpiration, Single, SingleUse},
		{2, NoExpiration, Multi, 2},
		{-2, tokenExpired, Infinite, InfiniteUse},
		{InfiniteUse, tokenExpired, Infinite, InfiniteUse},
		{0, tokenExpired, Single, SingleUse},
		{SingleUse, tokenExpired, Single, SingleUse},
		{2, tokenExpired, Multi, 2},
		{-2, tokenExp, Timed, InfiniteUse},
		{InfiniteUse, tokenExp, Timed, InfiniteUse},
		{0, tokenExp, Timed, SingleUse},
		{SingleUse, tokenExp, Timed, SingleUse},
		{2, tokenExp, Timed, 2},
	}
)

func TestUseToken(t *testing.T) {
	conn := accessTokenConn(t)
	tk := NewAccessToken(testToken, SingleUse, NoExpiration)
	err := conn.Create(tk).Error
	assert.NoError(t, err)
	assert.NotNil(t, tk)
	assert.True(t, tk.Usable())
	err = UseToken(conn, tk)
	assert.NoError(t, err)
	assert.Equal(t, 1, tk.Used)
	assert.NotNil(t, tk.UsedAt)
	assert.True(t, tk.DeletedAt.Valid)
	assert.False(t, tk.Usable())
	err = UseToken(conn, tk)
	assert.Error(t, err)
}
