package user

import (
	"net/mail"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/store/types/provider"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	email := tutils.RandomEmail()
	conn, c := tconn.TempConn(t)
	p := c.Provider()
	u := NewUser(p, -2, email, "", []byte(""), nil, nil)
	assert.Nil(t, u)
	u = NewUser(p, RoleUser, "", "", []byte(""), nil, nil)
	assert.Nil(t, u)
	u = NewUser(p, RoleUser, email, "", nil, nil, nil)
	assert.Nil(t, u)
	u = NewUser(p, RoleSystem, email, "", []byte(""), nil, nil)
	assert.Nil(t, u)
	u = NewUser(p, RoleUser, email, "", []byte(""), nil, nil)
	assert.False(t, u.IsSystemUser())
	assert.Equal(t, RoleUser, u.Role)
	assert.Equal(t, email, u.Email)
	err := conn.Create(u).Error
	assert.NoError(t, err)
	// email unique index constraint
	u = NewUser(p, RoleUser, email, "", []byte(""), nil, nil)
	err = conn.Create(u).Error
	assert.Error(t, err)
}

func TestUser_BeforeSave(t *testing.T) {
	conn, c := tconn.TempConn(t)
	u := &User{ID: SystemID}
	err := conn.Create(u).Error
	assert.Error(t, err)
	u.ID = uuid.New()
	err = conn.Create(u).Error
	assert.Error(t, err)
	u.Provider = c.Provider()
	err = conn.Create(u).Error
	assert.NoError(t, err)
	u.Provider = provider.Unknown
	err = conn.Save(u).Error
	assert.Error(t, err)
	u.ID = SystemID
	u.Provider = c.Provider()
	err = conn.Save(u).Error
	assert.Error(t, err)
}

func TestUser_EmailAddress(t *testing.T) {
	email := tutils.RandomEmail()
	un := utils.RandomUsername()
	u := &User{
		Email:    email,
		Username: un,
	}
	addr := u.EmailAddress()
	assert.EqualValues(t, &mail.Address{
		Name:    un,
		Address: email,
	}, addr)

}

func TestUser_Authenticate(t *testing.T) {
	const testPass = "password"
	conn, c := tconn.TempConn(t)
	email := tutils.RandomEmail()
	hash, err := utils.HashPassword(testPass)
	u := NewUser(c.Provider(), RoleUser, email, "", hash, nil, nil)
	err = u.Authenticate(testPass)
	assert.Error(t, err)
	err = conn.Save(u).Error
	require.NoError(t, err)
	err = u.Authenticate(testPass)
	assert.NoError(t, err)
	err = u.Authenticate("")
	assert.Error(t, err)
}

func TestUser_Status(t *testing.T) {
	email := tutils.RandomEmail()
	conn, c := tconn.TempConn(t)
	p := c.Provider()
	u := NewUser(p, RoleUser, email, "", []byte(""), nil, nil)
	err := conn.Save(u).Error
	require.NoError(t, err)
	tests := []struct {
		change     func()
		valid      bool
		admin      bool
		banned     bool
		locked     bool
		restricted bool
		confirmed  bool
		active     bool
		verified   bool
	}{
		{
			// normal user
			func() {},
			true,
			false,
			false,
			false,
			true,
			false,
			false,
			false,
		},
		{
			// confirmed
			func() {
				now := time.Now()
				u.ConfirmedAt = &now
				err = conn.Save(u).Error
				require.NoError(t, err)
			},
			true,
			false,
			false,
			false,
			true,
			true,
			false,
			false,
		},
		{
			// active
			func() {
				u.Status = Active
				err = conn.Save(u).Error
				require.NoError(t, err)
			},
			true,
			false,
			false,
			false,
			false,
			true,
			true,
			false,
		},
		{
			// was verified
			func() {
				now := time.Now()
				u.VerifiedAt = &now
				err = conn.Save(u).Error
				require.NoError(t, err)
			},
			true,
			false,
			false,
			false,
			false,
			true,
			true,
			false,
		},
		{
			// is verified
			func() {
				u.Status = Verified
				err = conn.Save(u).Error
				require.NoError(t, err)
			},
			true,
			false,
			false,
			false,
			false,
			true,
			true,
			true,
		},
		{
			// admin
			func() {
				u.Role = RoleAdmin
				err = conn.Save(u).Error
				require.NoError(t, err)
			},
			true,
			true,
			false,
			false,
			false,
			true,
			true,
			true,
		},
		{
			// locked
			func() {
				u.Status = Locked
				err = conn.Save(u).Error
				require.NoError(t, err)
			},
			true,
			true,
			false,
			true,
			false,
			false,
			false,
			false,
		},
		{
			// banned
			func() {
				u.Status = Banned
				err = conn.Save(u).Error
				require.NoError(t, err)
			},
			true,
			true,
			true,
			true,
			false,
			false,
			false,
			false,
		},
		{
			// invalid
			func() {
				err = conn.Delete(u).Error
				require.NoError(t, err)
			},
			false,
			false,
			true,
			true,
			false,
			false,
			false,
			false,
		},
	}
	for _, test := range tests {
		test.change()
		assert.Equal(t, test.valid, u.Valid())
		assert.Equal(t, test.admin, u.IsAdmin())
		assert.Equal(t, test.banned, u.IsBanned())
		assert.Equal(t, test.locked, u.IsLocked())
		assert.Equal(t, test.restricted, u.IsRestricted())
		assert.Equal(t, test.confirmed, u.IsConfirmed())
		assert.Equal(t, test.active, u.IsActive())
		assert.Equal(t, test.verified, u.IsVerified())
	}
}
