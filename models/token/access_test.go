package token

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/migration"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testToken    = "1234567890asdfghjkl="
	tokenExp     = time.Hour * 1
	tokenExpired = tokenExp * -1
	tokenTests   = []struct {
		uses int
		exp  time.Duration
		use  Type
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

func accessTokenConn(t *testing.T) *store.Connection {
	conn, _ := tconn.TempConn(t)
	mg := migration.NewMigrationWithIndexes("1",
		AccessToken{}, AccessTokenIndexes)
	err := conn.RunMigration(mg)
	require.NoError(t, err)
	return conn
}

func TestNewAccessToken(t *testing.T) {
	conn := accessTokenConn(t)
	at := AccessToken{Type: 255}
	assert.Equal(t, "invalid", at.Usage().String())
	err := conn.Create(at).Error
	assert.Error(t, err)
	err = conn.Transaction(func(tx *store.Connection) error {
		for i, test := range tokenTests {
			testToken = testToken + strconv.Itoa(i)
			tk := NewAccessToken(testToken, test.uses, test.exp)
			assert.NotNil(t, tk)
			assert.Equal(t, user.SystemID, tk.UserID)
			assert.Equal(t, testToken, tk.Token)
			assert.Equal(t, test.use, tk.Usage())
			assert.Equal(t, Access.String(), tk.Class().String())
			assert.Equal(t, user.SystemID, tk.IssuedTo())
			assert.Equal(t, tk.CreatedAt, tk.Issued())
			assert.Equal(t, time.Time{}, tk.LastUsed())
			assert.NotEmpty(t, tk.String())
			assert.Equal(t, test.max, tk.MaxUses)
			assert.Equal(t, 0, tk.Used)
			assert.Nil(t, tk.UsedAt)
			exp := test.exp
			if exp < NoExpiration {
				exp = NoExpiration
			}
			assert.Equal(t, exp, tk.Expiration)
			err := tx.Create(tk).Error
			assert.NoError(t, err)
			assert.True(t, tk.Usable())
			if tk.ExpiredAt == nil {
				assert.Equal(t, time.Time{}, tk.ExpirationDate())
			} else {
				assert.Equal(t, *tk.ExpiredAt, tk.ExpirationDate())
			}
			err = tx.Delete(tk).Error
			assert.NoError(t, err)
			assert.Equal(t, tk.DeletedAt.Time, tk.Revoked())
		}
		return nil
	})
	require.NoError(t, err)
	tk := NewAccessToken("", SingleUse, NoExpiration)
	assert.Nil(t, tk)
}

type badToken struct {
	AccessToken
}

// Class returns the class of the access token.
func (t *badToken) Class() Class {
	return ""
}

func TestAccessToken_BeforeCreate(t *testing.T) {
	conn := accessTokenConn(t)
	tk := NewAccessToken(testToken, SingleUse, NoExpiration)
	err := conn.Create(tk).Error
	assert.NoError(t, err)
	bad := NewAccessToken(testToken, SingleUse, NoExpiration)
	bad.Token = ""
	err = conn.Create(tk).Error
	assert.Error(t, err)
	bt := &badToken{}
	assert.False(t, bt.Usable())
	bt.Class()
	err = conn.Create(bt).Error
	assert.Error(t, err)
}

func TestAccessToken_AfterCreate(t *testing.T) {
	conn := accessTokenConn(t)
	err := conn.Transaction(func(tx *store.Connection) error {
		for i, test := range tokenTests {
			testToken = fmt.Sprintf("%s-%d", testToken, i)
			tk := NewAccessToken(testToken, SingleUse, test.exp)
			err := tx.Create(tk).Error
			assert.NoError(t, err)
			assert.True(t, tk.Usable())
			if test.exp <= NoExpiration {
				assert.Nil(t, tk.ExpiredAt)
				assert.Equal(t, time.Duration(0), tk.Expiration)
			} else {
				assert.Equal(t, test.exp, tk.Expiration)
				assert.NotNil(t, tk.ExpiredAt)
				diff := tk.ExpiredAt.Sub(tk.CreatedAt)
				assert.Equal(t, test.exp, diff)
			}
			tk = NewAccessToken(testToken, SingleUse, test.exp)
			err = tx.Create(tk).Error
			assert.Error(t, err)
		}
		return nil
	})
	require.NoError(t, err)
}

func TestAccessToken_Usable(t *testing.T) {
	const token = "token"
	conn := accessTokenConn(t)
	err := conn.Transaction(func(tx *store.Connection) error {
		for i, test := range tokenTests {
			testToken = fmt.Sprintf("%s-%d", testToken, i)
			tk := NewAccessToken(testToken, test.uses, test.exp)
			err := conn.Create(tk).Error
			assert.NoError(t, err)
			assert.True(t, tk.Usable())
			if tk.Type == Timed {
				if test.exp > 0 {
					assert.True(t, tk.Usable())

				}
			}
			if tk.MaxUses != InfiniteUse {
				tk.Used = 3
				assert.False(t, tk.Usable())
			}
			if tk.Type == Timed {
				tm := time.Now().UTC()
				tm = tm.Add(-time.Hour)
				tk.ExpiredAt = &tm
				assert.False(t, tk.Usable())
			}
			err = conn.Delete(tk).Error
			assert.NoError(t, err)
			assert.False(t, tk.Usable())
		}
		return nil
	})
	require.NoError(t, err)
	bad := NewAccessToken(token, SingleUse, NoExpiration)
	bad.Token = ""
	assert.False(t, bad.Usable())
}

func TestAccessToken_Use(t *testing.T) {
	conn := accessTokenConn(t)
	tk := NewAccessToken(testToken, SingleUse, NoExpiration)
	err := conn.Create(tk).Error
	assert.NoError(t, err)
	assert.NotNil(t, tk)
	assert.True(t, tk.Usable())
	tk.Use()
	assert.Equal(t, 1, tk.Used)
	assert.NotNil(t, tk.UsedAt)
	assert.False(t, tk.Usable())
	assert.NotNil(t, tk.LastUsed())

}
