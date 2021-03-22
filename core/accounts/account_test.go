package accounts

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/account"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const name = provider.Google

func testAccount(t *testing.T, conn *store.Connection, data types.Map) *account.Account {
	var aid = uuid.New().String()
	var mail = tutils.RandomEmail()
	var uid = uuid.New()
	la := account.NewAccount(name, aid, mail, data)
	la.UserID = uid
	err := conn.Create(la).Error
	require.NoError(t, err)
	return la
}

func TestGetAccount(t *testing.T) {
	conn, _ := tconn.TempConn(t)
	la := testAccount(t, conn, nil)
	got, err := GetAccount(conn, name, la.AccountID)
	assert.NoError(t, err)
	assert.Equal(t, la.Email, got.Email)
	got, err = GetAccount(conn, name, "")
	assert.Error(t, err)
	assert.Nil(t, got)
	got, err = GetAccount(conn, provider.Unknown, la.AccountID)
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestHasAccount(t *testing.T) {
	conn, _ := tconn.TempConn(t)
	la := testAccount(t, conn, nil)
	has, err := HasAccount(conn, name, la.AccountID)
	assert.NoError(t, err)
	assert.True(t, has)
	has, err = HasAccount(conn, name, "")
	assert.NoError(t, err)
	assert.False(t, has)
	has, err = HasAccount(conn, provider.Unknown, la.AccountID)
	assert.NoError(t, err)
	assert.False(t, has)
}

func TestUpdateAccount(t *testing.T) {
	conn, _ := tconn.TempConn(t)
	data := types.Map{
		"hello":  "world",
		"foobar": 13.37,
	}
	la := testAccount(t, conn, data)
	email := la.Email
	aid := la.AccountID
	emailIn := tutils.RandomEmail()
	dataIn := types.Map{
		"hello": "bar",
		"quack": "ok",
	}
	dataOut := types.Map{
		"hello":  "bar",
		"foobar": 13.37,
		"quack":  "ok",
	}
	tests := []struct {
		email       *string
		data        types.Map
		assertEmail string
		assertData  types.Map
		Ok          assert.BoolAssertionFunc
	}{
		{nil, nil, email, data, assert.False},
		{&email, nil, email, data, assert.False},
		{nil, types.Map{}, email, data, assert.False},
		{&emailIn, nil, emailIn, data, assert.True},
		{nil, dataIn, email, dataOut, assert.True},
		{&emailIn, dataIn, emailIn, dataOut, assert.True},
	}
	for _, test := range tests {
		ok, err := UpdateAccount(conn, la, test.email, test.data)
		assert.NoError(t, err)
		test.Ok(t, ok)
		assert.Equal(t, aid, la.AccountID)
		assert.Equal(t, test.assertEmail, la.Email)
		assert.Equal(t, test.assertData, la.Data)
		la.Email = email
		la.Data = data
		err = conn.Save(la).Error
		require.NoError(t, err)
	}
}
