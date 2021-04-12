package system

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/api/grpc/rpc/system"
	"github.com/jrapoport/gothic/models/account"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testServer(t *testing.T) *systemServer {
	srv, _ := tsrv.RPCServer(t, false)
	srv.Config().Signup.AutoConfirm = true
	return newSystemServer(srv)
}

func TestSystemServer_GetUser(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	srv.Config().MaskEmails = false
	ctx := context.Background()
	// no id or email
	req := &system.UserAccountRequest{}
	_, err := srv.GetUserAccount(ctx, req)
	assert.Error(t, err)
	// bad id
	req.Id = &system.UserAccountRequest_UserId{UserId: "1"}
	_, err = srv.GetUserAccount(ctx, req)
	assert.Error(t, err)
	// id not found
	req.Id = &system.UserAccountRequest_UserId{
		UserId: uuid.New().String(),
	}
	_, err = srv.GetUserAccount(ctx, req)
	assert.Error(t, err)
	// success
	u, _ := tcore.TestUser(t, srv.API, "", false)
	req.Id = &system.UserAccountRequest_UserId{
		UserId: u.ID.String(),
	}
	res, err := srv.GetUserAccount(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, u.ID.String(), res.Id)
	assert.Equal(t, u.Email, res.Email)
	assert.Equal(t, u.Username, res.Username)
	// bad email
	req = &system.UserAccountRequest{}
	req.Id = &system.UserAccountRequest_Email{
		Email: "@",
	}
	_, err = srv.GetUserAccount(ctx, req)
	assert.Error(t, err)
	// email not found
	req.Id = &system.UserAccountRequest_Email{
		Email: tutils.RandomEmail(),
	}
	_, err = srv.GetUserAccount(ctx, req)
	assert.Error(t, err)
	// success
	req.Id = &system.UserAccountRequest_Email{
		Email: u.Email,
	}
	res, err = srv.GetUserAccount(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, u.ID.String(), res.Id)
	assert.Equal(t, u.Email, res.Email)
	assert.Equal(t, u.Username, res.Username)
}

func TestSystemServer_LinkAccount(t *testing.T) {
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
	srv := testServer(t)
	srv.Config().MaskEmails = false
	ctx := context.Background()
	u, _ := tcore.TestUser(t, srv.API, "", false)
	// errors
	req := &system.LinkAccountRequest{
		UserId: uuid.Nil.String(),
	}
	_, err := srv.LinkAccount(ctx, req)
	assert.Error(t, err)
	req = &system.LinkAccountRequest{
		UserId: "1",
	}
	_, err = srv.LinkAccount(ctx, req)
	assert.Error(t, err)
	bad := account.NewAccount(p, aid, email, data)
	req = &system.LinkAccountRequest{
		UserId:  uuid.Nil.String(),
		Account: NewAccount(bad),
	}
	_, err = srv.LinkAccount(ctx, req)
	assert.Error(t, err)
	bad = account.NewAccount(provider.Unknown, aid, email, data)
	req = &system.LinkAccountRequest{
		UserId:  u.ID.String(),
		Account: NewAccount(bad),
	}
	_, err = srv.LinkAccount(ctx, req)
	assert.Error(t, err)
	req.Account.Type = uint32(account.Wallet)
	_, err = srv.LinkAccount(ctx, req)
	assert.Error(t, err)
	bad = account.NewAccount(p, "", email, data)
	req = &system.LinkAccountRequest{
		UserId:  u.ID.String(),
		Account: NewAccount(bad),
	}
	_, err = srv.LinkAccount(ctx, req)
	assert.Error(t, err)
	// success
	la1 := account.NewAccount(p, aid, email, data)
	req = &system.LinkAccountRequest{
		UserId:  u.ID.String(),
		Account: NewAccount(la1),
	}
	_, err = srv.LinkAccount(ctx, req)
	assert.NoError(t, err)
	links, err := srv.API.GetLinkedAccounts(nil, u.ID, account.Any, nil)
	require.NoError(t, err)
	require.Len(t, links, 1)
	assert.Equal(t, u.ID, links[0].UserID)
	assert.Equal(t, aid, links[0].AccountID)
	assert.Equal(t, email, links[0].Email)
	// relink error
	la1 = account.NewAccount(p, aid, email, data)
	req = &system.LinkAccountRequest{
		UserId:  u.ID.String(),
		Account: NewAccount(la1),
	}
	_, err = srv.LinkAccount(ctx, req)
	assert.Error(t, err)
	la2 := account.NewAccount(provider.GitHub, aid, email, data)
	req = &system.LinkAccountRequest{
		UserId:  u.ID.String(),
		Account: NewAccount(la2),
	}
	_, err = srv.LinkAccount(ctx, req)
	assert.NoError(t, err)
	links, err = srv.API.GetLinkedAccounts(nil, u.ID, account.Any, nil)
	require.NoError(t, err)
	require.Len(t, links, 2)
	assert.NotEqual(t, links[0].ID, links[1].ID)
}

func TestSystemServer_GetLinkedAccounts(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	srv.Config().MaskEmails = false
	ctx := context.Background()
	u, _ := tcore.TestUser(t, srv.API, "", false)
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
	for _, link := range links {
		req := &system.LinkAccountRequest{
			UserId:  u.ID.String(),
			Account: NewAccount(link),
		}
		_, err := srv.LinkAccount(ctx, req)
		require.NoError(t, err)
	}
	type Filters map[string]string
	tests := []struct {
		typ      account.Type
		filters  Filters
		expected []int
	}{
		{account.Any, nil, []int{0, 1, 2, 3}},
		{account.Payment, nil, []int{2, 3}},
		{account.Wallet, nil, []int{3}},
		{account.Payment | account.Wallet, nil, []int{2, 3}},
		{account.Any, Filters{
			key.Provider: links[0].Provider.String(),
		}, []int{0}},
		{account.Any, Filters{
			key.Email: links[2].Email,
		}, []int{2}},
		{account.Any, Filters{
			key.Token: links[1].Provider.ID().String(),
		}, []int{1}},
		{account.Any, Filters{
			"hello": "world",
		}, []int{0, 1, 2, 3}},
		{account.Payment, Filters{
			"hello": "world",
		}, []int{2, 3}},
		{account.Any, Filters{
			"bad": "unknown",
		}, []int{}},
	}
	for _, test := range tests {
		var reqType = uint32(test.typ)
		req := &system.LinkedAccountsRequest{
			UserId:  u.ID.String(),
			Type:    &reqType,
			Filters: test.filters,
		}
		res, err := srv.GetLinkedAccounts(ctx, req)
		require.NoError(t, err)
		require.Len(t, res.Linked, len(test.expected))
		for i, r := range res.Linked {
			idx := test.expected[i]
			link := links[idx]
			typ := account.Type(r.Type)
			assert.Equal(t, link.AccountID, r.AccountId)
			assert.True(t, typ.Has(test.typ))
			assert.Equal(t, link.Email, r.Email)
			assert.EqualValues(t, link.Provider, r.Provider)
			data := r.Data.AsMap()
			for k, v := range link.Data {
				assert.Equal(t, v, data[k])
			}
		}
	}
}
