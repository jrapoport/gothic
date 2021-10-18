package auth

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/api/grpc/rpc/auth"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/jwt"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthServer_RefreshBearerToken(t *testing.T) {
	t.Parallel()
	srv := testServer(t)
	ctx := context.Background()
	// invalid req
	_, err := srv.RefreshBearerToken(ctx, nil)
	assert.Error(t, err)
	// empty email
	req := &auth.RefreshTokenRequest{}
	_, err = srv.RefreshBearerToken(ctx, req)
	assert.Error(t, err)
	// bad email
	req.Token = "bad"
	_, err = srv.RefreshBearerToken(ctx, req)
	assert.Error(t, err)
	u, _ := tcore.TestUser(t, srv.API, "", false)
	bt, err := srv.GrantBearerToken(context.Background(), u)
	require.NoError(t, err)
	req = &auth.RefreshTokenRequest{
		Token: bt.RefreshToken.Token,
	}
	res, err := srv.RefreshBearerToken(ctx, req)
	assert.NoError(t, err)
	claims, err := jwt.ParseUserClaims(srv.Config().JWT, res.Access)
	assert.NoError(t, err)
	require.NotNil(t, claims)
	uid, err := uuid.Parse(claims.Subject())
	assert.NoError(t, err)
	u2, err := srv.GetUser(uid)
	assert.NoError(t, err)
	assert.Equal(t, u2.ID.String(), claims.Subject())
	au, err := srv.GetAuthenticatedUser(u2.ID)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, au.ID)
	res, err = srv.RefreshBearerToken(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, res)
}
