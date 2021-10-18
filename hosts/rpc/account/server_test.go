package account

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/api/grpc/rpc"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/jwt"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPass = "SXJAm7qJ4?3dH!aN8T3f5p!oNnpXbaRy#Gtx#8jG"

func testServer(t *testing.T) *server {
	srv, _ := tsrv.RPCServer(t, false)
	return newServer(srv)
}

func testUser(t *testing.T, srv *server) *user.User {
	em := tutils.RandomEmail()
	ctx := rpc.RequestContext(nil)
	ctx.SetProvider(srv.Provider())
	u, err := srv.API.Signup(ctx, em, "", testPass, nil)
	require.NoError(t, err)
	require.NotNil(t, u)
	return u
}

func testAuthCtx(t *testing.T, srv *server, u *user.User) context.Context {
	bt, err := srv.GrantBearerToken(context.Background(), u)
	require.NoError(t, err)
	claims, err := jwt.ParseUserClaims(srv.Config().JWT, bt.Token)
	require.NoError(t, err)
	require.NotNil(t, claims)
	ctx := context.Background()
	return context.WithContext(rpc.WithClaims(ctx, claims))
}

func assertUserResponse(t *testing.T, srv *server, test *rpc.UserResponse, res *api.UserResponse) {
	assert.Equal(t, test.Username, res.Username)
	assert.Equal(t, test.Role, res.Role)
	em := test.Email
	if srv.Config().MaskEmails {
		em = utils.MaskEmail(em)
	}
	assert.Equal(t, em, res.Email)
	data := res.Data.AsMap()
	for k, v := range test.Data.AsMap() {
		assert.Equal(t, v, data[k])
	}
	require.NotNil(t, res.Token)
	require.NotEmpty(t, res.Token.Access)
	claims, err := jwt.ParseUserClaims(srv.Config().JWT, res.Token.Access)
	require.NoError(t, err)
	require.NotNil(t, claims)
	assert.EqualValues(t, tokens.Bearer, res.Token.Type)
	u, err := srv.GetUser(uuid.MustParse(claims.Subject()))
	require.NoError(t, err)
	require.NotNil(t, u)
	assert.Equal(t, user.RoleUser, u.Role)
	assert.Equal(t, test.Email, u.Email)
}
