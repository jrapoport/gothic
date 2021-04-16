package admin

import (
	"context"
	"testing"

	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/core/codes"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/models/code"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestSignupServer_CreateSignupCodes(t *testing.T) {
	t.Parallel()
	const testLen = 10
	s, _ := tsrv.RPCServer(t, false)
	srv := newAdminServer(s)
	ctx := context.Background()
	_, err := srv.CreateSignupCodes(ctx, nil)
	assert.Error(t, err)
	req := &admin.CreateSignupCodesRequest{
		Uses:  code.InfiniteUse,
		Count: testLen,
	}
	_, tok := tcore.TestUser(t, srv.API, "", false)
	ctx = tsrv.RPCAuthContext(t, srv.Config(), tok)
	_, err = srv.CreateSignupCodes(ctx, req)
	assert.Error(t, err)
	ctx = metadata.NewIncomingContext(context.Background(),
		metadata.Pairs(rpc.RootPassword, "bad"))
	_, err = srv.CreateSignupCodes(ctx, req)
	assert.Error(t, err)
	pw := s.Config().RootPassword
	ctx = metadata.NewIncomingContext(context.Background(),
		metadata.Pairs(rpc.RootPassword, pw))
	res, err := srv.CreateSignupCodes(ctx, req)
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, res.GetCodes(), testLen)
}

func TestSignupServer_CheckSignupCode(t *testing.T) {
	t.Parallel()
	s, _ := tsrv.RPCServer(t, false)
	srv := newAdminServer(s)
	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.Pairs(rpc.RootPassword, s.Config().RootPassword))
	cr, err := srv.CreateSignupCodes(ctx, &admin.CreateSignupCodesRequest{
		Uses:  code.SingleUse,
		Count: 1,
	})
	require.NoError(t, err)
	require.NotNil(t, cr)
	require.Len(t, cr.GetCodes(), 1)
	_, err = srv.CheckSignupCode(ctx, nil)
	assert.Error(t, err)
	req := &admin.CheckSignupCodeRequest{
		Code: "bad",
	}
	_, err = srv.CheckSignupCode(ctx, req)
	assert.Error(t, err)
	test := cr.GetCodes()[0]
	req = &admin.CheckSignupCodeRequest{
		Code: test,
	}
	res, err := srv.CheckSignupCode(ctx, req)
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.True(t, res.Valid)
	assert.Equal(t, test, res.Code)
	conn := tconn.Conn(t, srv.Config())
	sc, err := codes.GetSignupCode(conn, test)
	require.NoError(t, err)
	require.NotNil(t, sc)
	sc.Used = 1
	err = conn.Save(sc).Error
	require.NoError(t, err)
	res, err = srv.CheckSignupCode(ctx, req)
	assert.NoError(t, err)
	assert.False(t, res.Valid)
	err = conn.Delete(sc).Error
	require.NoError(t, err)
	_, err = srv.CheckSignupCode(ctx, req)
	assert.Error(t, err)
}

func TestSignupServer_DeleteSignupCode(t *testing.T) {
	t.Parallel()
	s, _ := tsrv.RPCServer(t, false)
	srv := newAdminServer(s)
	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.Pairs(rpc.RootPassword, s.Config().RootPassword))
	_, err := srv.DeleteSignupCode(ctx, nil)
	assert.Error(t, err)
	req := &admin.DeleteSignupCodeRequest{
		Code: "bad",
	}
	_, err = srv.DeleteSignupCode(ctx, req)
	assert.Error(t, err)
	cr, err := srv.CreateSignupCodes(ctx, &admin.CreateSignupCodesRequest{
		Uses:  code.SingleUse,
		Count: 1,
	})
	require.NoError(t, err)
	require.NotNil(t, cr)
	require.Len(t, cr.GetCodes(), 1)
	test := cr.GetCodes()[0]
	req = &admin.DeleteSignupCodeRequest{
		Code: test,
	}
	_, err = srv.DeleteSignupCode(ctx, req)
	assert.NoError(t, err)
	conn := tconn.Conn(t, srv.Config())
	_, err = codes.GetSignupCode(conn, test)
	require.Error(t, err)
}
