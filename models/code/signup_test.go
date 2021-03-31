package code

import (
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testUser(c *config.Config) *user.User {
	p := c.Provider()
	em := tutils.RandomEmail()
	r := user.RoleUser
	return user.NewUser(p, r, em, "", []byte(""), nil, nil)
}

func TestSignupCode_HasCode(t *testing.T) {
	t.Parallel()
	for _, f := range testFormats {
		t.Run(testName(f), func(t *testing.T) {
			t.Parallel()
			testSignupCodeHasCode(t, f)
		})
	}
}

func testSignupCodeHasCode(t *testing.T, f Format) {
	conn, _ := tconn.TempConn(t)
	testCode := func() *SignupCode {
		sc := NewSignupCode(user.SystemID, f, SingleUse)
		err := conn.Create(sc).Error
		require.NoError(t, err)
		return sc
	}
	deletedCode := testCode()
	err := conn.Delete(deletedCode).Error
	require.NoError(t, err)
	tests := []struct {
		sc  *SignupCode
		Err assert.ErrorAssertionFunc
		Has assert.BoolAssertionFunc
	}{
		{&SignupCode{}, assert.Error, assert.False},
		{NewSignupCode(user.SystemID, f, SingleUse), assert.NoError, assert.False},
		{testCode(), assert.NoError, assert.True},
		{deletedCode, assert.NoError, assert.False},
	}
	var has bool
	for _, test := range tests {
		has, err = test.sc.HasCode(conn)
		test.Err(t, err)
		test.Has(t, has)
	}
}

func TestSignupCode_UseCode(t *testing.T) {
	t.Parallel()
	for _, f := range testFormats {
		t.Run(testName(f), func(t *testing.T) {
			t.Parallel()
			testSignupCodeUseCode(t, f)
		})
	}
	conn, c := tconn.TempConn(t)
	sc := NewSignupCode(user.SystemID, PIN, SingleUse)
	assert.Equal(t, Signup, sc.Class())
	sc.Token = utils.DebugPIN
	u := testUser(c)
	err := conn.Create(u).Error
	assert.NoError(t, err)
	err = sc.UseCode(conn, u)
	assert.NoError(t, err)
	// reusing a single use code would normally fail
	err = sc.UseCode(conn, u)
	assert.NoError(t, err)
}

func testSignupCodeUseCode(t *testing.T, f Format) {
	conn, c := tconn.TempConn(t)
	su := user.NewSystemUser()
	iu := testUser(c)
	vu := testUser(c)
	err := conn.Create(vu).Error
	assert.NoError(t, err)
	du := testUser(c)
	err = conn.Create(du).Error
	assert.NoError(t, err)
	err = conn.Delete(du).Error
	assert.NoError(t, err)
	testCode := func() *SignupCode {
		sc := NewSignupCode(user.SystemID, f, SingleUse)
		err = conn.Create(sc).Error
		require.NoError(t, err)
		return sc
	}
	usedCode := testCode()
	usedCode.Used = SingleUse
	deletedCode := testCode()
	err = conn.Delete(deletedCode).Error
	require.NoError(t, err)
	tests := []struct {
		sc  *SignupCode
		u   *user.User
		Err assert.ErrorAssertionFunc
	}{
		{&SignupCode{}, vu, assert.Error},
		{NewSignupCode(user.SystemID, f, SingleUse), vu, assert.Error},
		{usedCode, vu, assert.Error},
		{deletedCode, vu, assert.Error},
		{testCode(), nil, assert.Error},
		{testCode(), su, assert.Error},
		{testCode(), iu, assert.Error},
		{testCode(), du, assert.Error},
		{testCode(), vu, assert.NoError},
	}
	err = conn.Transaction(func(tx *store.Connection) error {
		for _, test := range tests {
			err = test.sc.UseCode(tx, test.u)
			test.Err(t, err)
			if err == nil {
				assert.False(t, test.sc.Usable())
				assert.Equal(t, SingleUse, test.sc.Used)
				u := &user.User{ID: test.u.ID}
				err = tx.Find(u).Error
				assert.NoError(t, err)
				assert.Equal(t, &test.sc.ID, u.SignupCode)
			}
		}
		return nil
	})
	require.NoError(t, err)
}
