package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserToken(t *testing.T) {
	c := tconf.Config(t)
	c.JWT.Audience = "test"
	c.JWT.Expiration = 100 * time.Minute
	now := time.Now()
	u := &user.User{
		ID:          uuid.New(),
		Provider:    c.Provider(),
		Role:        user.RoleUser,
		Status:      user.Active,
		Email:       tutils.RandomEmail(),
		Username:    utils.RandomUsername(),
		CreatedAt:   now,
		ConfirmedAt: &now,
		VerifiedAt:  &now,
	}
	tok := NewUserToken(c.JWT, u)
	assert.NotNil(t, tok)
	b, err := tok.Bearer()
	require.NoError(t, err)
	claims, err := ParseUserClaims(c.JWT, b)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, claims.UserID())
	// bad sub
	u.ID = uuid.Nil
	tok = NewUserToken(c.JWT, u)
	assert.NotNil(t, tok)
	b, err = tok.Bearer()
	require.NoError(t, err)
	claims, err = ParseUserClaims(c.JWT, b)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, claims.UserID())
	// invalid sub
	tok = NewUserToken(c.JWT, u)
	assert.NotNil(t, tok)
	// error
	_, err = ParseUserClaims(c.JWT, "bad")
	assert.Error(t, err)
	// bad subject
	usr := UserClaims{}
	usr.Subject = "1"
	assert.Equal(t, uuid.Nil, usr.UserID())
}
