package admin

import (
	"context"

	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	core_ctx "github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/models/user"
	"google.golang.org/grpc"
)

type adminServer struct {
	admin.UnimplementedAdminServer
	*rpc.Server
}

var _ admin.AdminServer = (*adminServer)(nil)

func newAdminServer(srv *rpc.Server) *adminServer {
	srv.Logger = srv.WithName("admin")
	return &adminServer{Server: srv}
}

// RegisterServer registers a new admin server.
func RegisterServer(s *grpc.Server, srv *rpc.Server) {
	admin.RegisterAdminServer(s, newAdminServer(srv))
}

func (s *adminServer) rootContext(ctx context.Context) context.Context {
	rtx := rpc.RequestContext(ctx)
	pw := rpc.GetRootPassword(rtx)
	if pw == "" {
		return rtx
	}
	root, err := s.API.GetSuperAdmin(pw)
	if err != nil {
		return rtx
	}
	rtx.SetProvider(root.Provider)
	rtx.SetAdminID(root.ID)
	return rtx
}

type roleKey struct{}

func adminRoleFromContext(ctx context.Context) user.Role {
	role, _ := ctx.Value(roleKey{}).(user.Role)
	return role
}

func (s *adminServer) adminRequestContext(ctx context.Context) (core_ctx.Context, error) {
	ctx = s.rootContext(ctx)
	role, err := s.Server.ValidateAdmin(ctx)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, roleKey{}, role)
	return rpc.RequestContext(ctx), nil
}
