package tokens

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrantConfirmToken(t *testing.T) {
	t.Parallel()
	conn, _ := tconn.TempConn(t)
	uid := uuid.New()
	ct, err := GrantConfirmToken(conn, uid, token.NoExpiration)
	assert.NoError(t, err)
	require.NotNil(t, ct)
	assert.NotEmpty(t, ct.AccessToken)
	assert.Equal(t, uid, ct.UserID)
	// system user id
	_, err = GrantConfirmToken(conn, user.SystemID, token.NoExpiration)
	assert.Error(t, err)
}

func TestGetConfirmToken(t *testing.T) {
	t.Parallel()
	conn, _ := tconn.TempConn(t)
	uid := uuid.New()
	test, err := GrantConfirmToken(conn, uid, token.NoExpiration)
	assert.NoError(t, err)
	assert.NotNil(t, test)
	ct, err := GetConfirmToken(conn, test.String())
	assert.NoError(t, err)
	assert.Equal(t, test.UserID, ct.UserID)
	assert.Equal(t, test.Token, ct.Token)
	_, err = GetConfirmToken(conn, "")
	assert.Error(t, err)
}

func TestConfirmTokenSent(t *testing.T) {
	t.Parallel()
	conn, _ := tconn.TempConn(t)
	uid := uuid.New()
	ct, err := GrantConfirmToken(conn, uid, token.NoExpiration)
	assert.NoError(t, err)
	err = ConfirmTokenSent(conn, ct)
	assert.NoError(t, err)
	assert.NotNil(t, ct.SentAt)
}
