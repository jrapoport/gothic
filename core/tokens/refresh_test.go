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

func TestGrantRefreshToken(t *testing.T) {
	conn, _ := tconn.TempConn(t)
	uid := uuid.New()
	testGrant := func(userID uuid.UUID) *token.RefreshToken {
		rt, err := GrantRefreshToken(conn, userID)
		assert.NoError(t, err)
		require.NotNil(t, rt)
		assert.NotEmpty(t, rt.AccessToken)
		assert.Equal(t, uid, rt.UserID)
		return rt
	}
	rt1 := testGrant(uid)
	rt2 := testGrant(uid)
	assert.Equal(t, rt1.ID, rt2.ID)
	assert.Equal(t, rt1.Token, rt2.Token)
	// system user id
	_, err := GrantRefreshToken(conn, user.SystemID)
	assert.Error(t, err)
}

func TestSwapRefreshToken(t *testing.T) {
	conn, _ := tconn.TempConn(t)
	uid := uuid.New()
	_, err := SwapRefreshToken(conn, uid, "")
	assert.Error(t, err)
	rt, err := GrantRefreshToken(conn, uid)
	assert.NoError(t, err)
	assert.False(t, rt.DeletedAt.Valid)
	st, err := SwapRefreshToken(conn, uid, rt.String())
	assert.NoError(t, err)
	assert.Equal(t, uid, st.UserID)
	_, err = GetUsableRefreshToken(conn, rt.Token)
	assert.Error(t, err)
	assert.NotEqual(t, rt.ID, st.ID)
	assert.Equal(t, rt.UserID, st.UserID)
}

func TestRevokeAllRefreshTokens(t *testing.T) {
	conn, _ := tconn.TempConn(t)
	uid := uuid.New()
	_, err := GrantRefreshToken(conn, uid)
	assert.NoError(t, err)
	for i := 0; i < 10; i++ {
		_, err = GrantRefreshToken(conn, uid)
	}
	count := func() int {
		var count int64
		err = conn.Unscoped().
			Model(token.RefreshToken{}).
			Where("user_id = ?", uid).
			Count(&count).Error
		require.NoError(t, err)
		return int(count)
	}
	assert.Equal(t, 1, count())
	err = RevokeAllRefreshTokens(conn, uid)
	assert.NoError(t, err)
	assert.Equal(t, 0, count())
}

func TestHasUsableRefreshToken(t *testing.T) {
	conn, c := tconn.TempConn(t)
	uid := uuid.New()
	has, err := HasUsableRefreshToken(conn, uid)
	assert.NoError(t, err)
	assert.False(t, has)
	u := testUser(t, conn, c)
	rt, err := GrantRefreshToken(conn, u.ID)
	assert.NoError(t, err)
	has, err = HasUsableRefreshToken(conn, u.ID)
	assert.NoError(t, err)
	assert.True(t, has)
	err = conn.Delete(rt).Error
	assert.NoError(t, err)
	has, err = HasUsableRefreshToken(conn, u.ID)
	assert.NoError(t, err)
	assert.False(t, has)
}
