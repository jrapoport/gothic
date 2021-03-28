package users

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/account"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkAccount(t *testing.T) {
	t.Parallel()
	var (
		p     = provider.Google
		aid   = uuid.New().String()
		email = tutils.RandomEmail()
		data  = types.Map{
			"hello":       "world",
			key.IPAddress: "127.0.0.1",
		}
	)
	conn, c := tconn.TempConn(t)
	u := testUser(t, conn, c.Provider())
	// errors
	err := LinkAccount(conn, uuid.Nil, nil)
	assert.Error(t, err)
	bad := account.NewAccount(p, aid, email, data)
	err = LinkAccount(conn, uuid.Nil, bad)
	assert.Error(t, err)
	bad = account.NewAccount(provider.Unknown, aid, email, data)
	err = LinkAccount(conn, u.ID, bad)
	assert.Error(t, err)
	bad = account.NewAccount(p, "", email, data)
	err = LinkAccount(conn, u.ID, bad)
	assert.Error(t, err)
	// success
	la1 := account.NewAccount(p, aid, email, data)
	err = LinkAccount(conn, u.ID, la1)
	assert.NoError(t, err)
	require.NotNil(t, la1)
	assert.Equal(t, aid, la1.AccountID)
	assert.Equal(t, email, la1.Email)
	assert.Equal(t, data, la1.Data)
	assert.Equal(t, u.ID, la1.UserID)
	cnt := conn.Model(u).Association("Linked").Count()
	assert.Equal(t, int64(1), cnt)
	// relink error
	la1 = account.NewAccount(p, aid, email, data)
	err = LinkAccount(conn, u.ID, la1)
	assert.Error(t, err)
	// 2nd account
	la2 := account.NewAccount(provider.GitHub, aid, email, data)
	err = LinkAccount(conn, u.ID, la2)
	assert.NoError(t, err)
	require.NotNil(t, la2)
	assert.NotEqual(t, la1.ID, la2.ID)
	cnt = conn.Model(u).Association("Linked").Count()
	assert.Equal(t, int64(2), cnt)
}

func TestGetLinkedAccounts(t *testing.T) {
	t.Parallel()
	testAccount := func(name provider.Name) *account.Account {
		link := account.NewAccount(name,
			uuid.New().String(),
			tutils.RandomEmail(),
			types.Map{
				"hello":   "world",
				key.Token: name.ID().String(),
			})
		return link
	}
	links := []*account.Account{
		testAccount(provider.Google),
		testAccount(provider.GitHub),
		testAccount(provider.Stripe),
		testAccount(provider.PayPal),
	}
	c := tconf.Config(t)
	conn, c := tconn.TempConn(t)
	u := testUser(t, conn, c.Provider())
	for _, link := range links {
		err := LinkAccount(conn, u.ID, link)
		require.NoError(t, err)
	}
	cnt := conn.Model(u).Association("Linked").Count()
	assert.Equal(t, int64(len(links)), cnt)
	tests := []struct {
		typ      account.Type
		filters  store.Filters
		expected []int
	}{
		{account.Any, nil, []int{0, 1, 2, 3}},
		{account.Payment, nil, []int{2, 3}},
		{account.Wallet, nil, []int{3}},
		{account.Payment | account.Wallet, nil, []int{2, 3}},
		{account.Any, store.Filters{
			key.Provider: links[0].Provider,
		}, []int{0}},
		{account.Any, store.Filters{
			key.Email: links[2].Email,
		}, []int{2}},
		{account.Any, store.Filters{
			key.Token: links[1].Provider.ID().String(),
		}, []int{1}},
		{account.Any, store.Filters{
			"hello": "world",
		}, []int{0, 1, 2, 3}},
		{account.Payment, store.Filters{
			"hello": "world",
		}, []int{2, 3}},
		{account.Any, store.Filters{
			"bad": "unknown",
		}, []int{}},
	}
	for _, test := range tests {
		res, err := GetLinkedAccounts(conn, u.ID, test.typ, test.filters)
		require.NoError(t, err)
		require.Len(t, res, len(test.expected))
		for i, r := range res {
			idx := test.expected[i]
			link := links[idx]
			assert.Equal(t, link.ID, r.ID)
			assert.True(t, r.Type.Has(test.typ))
			assert.Equal(t, link.Email, r.Email)
			assert.Equal(t, link.Provider, r.Provider)
			assert.EqualValues(t, link.Data, r.Data)
		}
	}
	// bad user id
	_, err := GetLinkedAccounts(conn, uuid.Nil, account.Any, nil)
	assert.Error(t, err)
	// user not found
	_, err = GetLinkedAccounts(conn, uuid.New(), account.Any, nil)
	assert.Error(t, err)
}
