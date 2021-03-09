package code

import (
	"testing"

	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignupCode_HasCode(t *testing.T) {
	for _, f := range testFormats {
		t.Run(testName(f), func(t *testing.T) {
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
	for _, f := range testFormats {
		t.Run(testName(f), func(t *testing.T) {
			testSignupCodeUseCode(t, f)
		})
	}
}

func testSignupCodeUseCode(t *testing.T, f Format) {
	conn, c := tconn.TempConn(t)
	testUser := func() *user.User {
		p := c.Provider()
		em := tutils.RandomEmail()
		r := user.RoleUser
		return user.NewUser(p, r, em, "", []byte(""), nil, nil)
	}
	su := user.NewSystemUser()
	iu := testUser()
	vu := testUser()
	err := conn.Create(vu).Error
	assert.NoError(t, err)
	du := testUser()
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
	for _, test := range tests {
		err = test.sc.UseCode(conn, test.u)
		test.Err(t, err)
		if err == nil {
			assert.False(t, test.sc.Usable())
			assert.Equal(t, SingleUse, test.sc.Used)
			u := &user.User{ID: test.u.ID}
			err = conn.Find(u).Error
			assert.NoError(t, err)
			assert.Equal(t, &test.sc.ID, u.SignupCode)
		}
	}
}
