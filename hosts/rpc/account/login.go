package account

import (
	"context"
	"errors"

	"github.com/jrapoport/gothic/api/grpc/rpc"
	"github.com/jrapoport/gothic/api/grpc/rpc/account"
	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc/codes"
)

func (s *server) Login(ctx context.Context,
	req *account.LoginRequest) (*api.UserResponse, error) {
	if req == nil {
		err := errors.New("request not found")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	rtx := rpc.RequestContext(ctx)
	rtx.SetProvider(s.Provider())
	s.Debugf("login user: %v (%v)", req, rtx)
	u, err := s.API.Login(rtx, req.Email, req.Password)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	res, err := rpc.NewUserResponse(u)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	if s.Config().MaskEmails {
		res.MaskEmail()
	}
	bt, err := s.GrantBearerToken(rtx, u)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	res.Token = rpc.NewBearerResponse(bt)
	return (*api.UserResponse)(res), nil
}
