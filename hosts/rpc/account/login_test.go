package account

import (
	"context"
	"testing"

	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/core/tokens/jwt"
	"github.com/jrapoport/gothic/protobuf/grpc/rpc/account"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountServer_Login(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	u := testUser(t, srv)
	// invalid req
	ctx := context.Background()
	_, err := srv.Login(ctx, nil)
	assert.Error(t, err)
	// no email
	_, err = srv.Login(ctx, &account.LoginRequest{})
	assert.Error(t, err)
	// bad email
	_, err = srv.Login(ctx, &account.LoginRequest{
		Email: "bad",
	})
	assert.Error(t, err)
	// not found
	_, err = srv.Login(ctx, &account.LoginRequest{
		Email: "bad@example.com",
	})
	assert.Error(t, err)
	// bad password
	_, err = srv.Login(ctx, &account.LoginRequest{
		Email:    u.Email,
		Password: "",
	})
	assert.Error(t, err)
	// login
	_, err = srv.GetAuthenticatedUser(u.ID)
	assert.Error(t, err)
	res, err := srv.Login(ctx, &account.LoginRequest{
		Email:    u.Email,
		Password: testPass,
	})
	assert.NoError(t, err)
	claims, err := jwt.ParseUserClaims(srv.Config().JWT, res.Token.Access)
	require.NoError(t, err)
	require.NotNil(t, claims)
	assert.EqualValues(t, tokens.Bearer, res.Token.Type)
	assert.Equal(t, u.ID.String(), claims.Subject)
	assert.Equal(t, res.Email, utils.MaskEmail(u.Email))
	au, err := srv.GetAuthenticatedUser(u.ID)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, au.ID)
}
