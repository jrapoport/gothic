package code

import (
	"github.com/jrapoport/gothic/utils"
	"testing"
	"time"

	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/migration"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testFormats = []Format{Invite, PIN}

func testName(f Format) string {
	switch f {
	case Invite:
		return "Token"
	case PIN:
		return "PIN"
	default:
		return ""
	}
}

func accessCodeConn(t *testing.T) *store.Connection {
	conn, _ := tconn.TempConn(t)
	mg := migration.NewMigrationWithIndexes("1",
		AccessCode{}, token.AccessTokenIndexes)
	err := conn.RunMigration(mg)
	require.NoError(t, err)
	return conn
}

func TestNewAccessCode(t *testing.T) {
	t.Parallel()
	for _, f := range testFormats {
		t.Run(testName(f), func(t *testing.T) {
			t.Parallel()
			testNewAccessCode(t, f)
		})
	}
}

func testNewAccessCode(t *testing.T, f Format) {
	conn := accessCodeConn(t)
	futureExp := time.Hour * 1
	pastExp := futureExp * -1
	tests := []struct {
		uses int
		exp  time.Duration
		use  token.Type
		max  int
	}{
		{-2, NoExpiration, token.Infinite, InfiniteUse},
		{InfiniteUse, NoExpiration, token.Infinite, InfiniteUse},
		{0, NoExpiration, token.Single, SingleUse},
		{SingleUse, NoExpiration, token.Single, SingleUse},
		{2, NoExpiration, token.Multi, 2},
		{-2, pastExp, token.Infinite, InfiniteUse},
		{InfiniteUse, pastExp, token.Infinite, InfiniteUse},
		{0, pastExp, token.Single, SingleUse},
		{SingleUse, pastExp, token.Single, SingleUse},
		{2, pastExp, token.Multi, 2},
		{-2, futureExp, token.Timed, InfiniteUse},
		{InfiniteUse, futureExp, token.Timed, InfiniteUse},
		{0, futureExp, token.Timed, SingleUse},
		{SingleUse, futureExp, token.Timed, SingleUse},
		{2, futureExp, token.Timed, 2},
	}
	err := conn.Transaction(func(tx *store.Connection) error {
		for _, test := range tests {
			code := NewAccessCode(f, test.uses, test.exp)
			assert.NotNil(t, code)
			assert.Equal(t, user.SystemID, code.UserID)
			assert.Equal(t, test.use, code.Usage())
			assert.NotEmpty(t, code.String())
			assert.Equal(t, test.max, code.MaxUses)
			assert.Equal(t, 0, code.Used)
			assert.Nil(t, code.UsedAt)
			exp := test.exp
			if exp < NoExpiration {
				exp = NoExpiration
			}
			assert.Equal(t, exp, code.Expiration)
			assert.False(t, code.Usable())
			err := tx.Create(code).Error
			assert.NoError(t, err)
			assert.True(t, code.Usable())
		}
		return nil
	})
	require.NoError(t, err)
	code := NewAccessCode(255, SingleUse, NoExpiration)
	assert.NotNil(t, code)
}

func TestAccessCode_BeforeCreate(t *testing.T) {
	t.Parallel()
	conn := accessCodeConn(t)
	for _, f := range testFormats {
		t.Run(testName(f), func(t *testing.T) {
			t.Parallel()
			testAccessCodeBeforeCreate(t, conn, f)
		})
	}
	ac := NewAccessCode(Invite, SingleUse, NoExpiration)
	err := conn.Create(ac).Error
	assert.NoError(t, err)
	ac = NewAccessCode(PIN, SingleUse, NoExpiration)
	ac.Token = utils.DebugPIN
	err = conn.Create(ac).Error
	assert.Error(t, err)
	assert.True(t, ac.Usable())
}

func testAccessCodeBeforeCreate(t *testing.T, conn *store.Connection, f Format) {
	ac := NewAccessCode(f, SingleUse, NoExpiration)
	err := conn.Create(ac).Error
	assert.NoError(t, err)
	bad := NewAccessCode(f, SingleUse, NoExpiration)
	bad.Token = ""
	err = conn.Create(bad).Error
	assert.Error(t, err)
}

func TestAccessCode_Usable(t *testing.T) {
	t.Parallel()
	for _, f := range testFormats {
		t.Run(testName(f), func(t *testing.T) {
			t.Parallel()
			testAccessCodeUsable(t, f)
		})
	}
}

func testAccessCodeUsable(t *testing.T, f Format) {
	conn := accessCodeConn(t)
	expiration := time.Hour * 1
	tests := []struct {
		uses int
		exp  time.Duration
	}{
		{-2, NoExpiration},
		{InfiniteUse, NoExpiration},
		{0, NoExpiration},
		{SingleUse, NoExpiration},
		{2, NoExpiration},
		{-2, expiration},
		{InfiniteUse, expiration},
		{0, expiration},
		{SingleUse, expiration},
		{2, expiration},
	}
	err := conn.Transaction(func(tx *store.Connection) error {
		for _, test := range tests {
			code := NewAccessCode(f, test.uses, test.exp)
			err := tx.Create(code).Error
			assert.NoError(t, err)
			assert.True(t, code.Usable())
			if code.MaxUses != InfiniteUse {
				code.Used = 3
				assert.False(t, code.Usable())
			}
			if code.Type == token.Timed {
				tm := time.Now().UTC().Add(-1 * time.Hour)
				code.ExpiredAt = &tm
				assert.False(t, code.Usable())
			}
			err = tx.Delete(code).Error
			assert.NoError(t, err)
			assert.False(t, code.Usable())
		}
		return nil
	})
	require.NoError(t, err)
}
