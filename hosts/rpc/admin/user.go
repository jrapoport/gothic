package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc/codes"
)

func (s *adminServer) CreateUser(ctx context.Context, req *admin.CreateUserRequest) (*admin.CreateUserResponse, error) {
	if req == nil {
		return nil, s.RPCError(codes.InvalidArgument, nil)
	}
	email := req.GetEmail()
	username := req.GetUsername()
	pw := req.GetPassword()
	data := req.GetData().AsMap()
	adminRole := req.GetAdmin()
	if email == "" {
		err := errors.New("email required")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}

	rtx := rpc.RequestContext(ctx)
	root := rpc.GetRootPassword(ctx)
	if root != "" {
		sa, err := s.API.GetSuperAdmin(root)
		if err != nil {
			return nil, s.RPCError(codes.PermissionDenied, err)
		}
		rtx.SetProvider(sa.Provider)
		rtx.SetAdminID(sa.ID)
	} else if adminRole {
		err := fmt.Errorf("super admin access required to create admin users")
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	// we don't need to validate the admin here since AdminCreateUser will do it.
	u, err := s.API.AdminCreateUser(rtx, email, username, pw, data, adminRole)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	res := &admin.CreateUserResponse{
		Role:  u.Role.String(),
		Email: u.Email,
	}
	return res, nil
}
