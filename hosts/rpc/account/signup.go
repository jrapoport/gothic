package account

import (
	"context"
	"errors"

	"github.com/jrapoport/gothic/api/grpc/rpc"
	"github.com/jrapoport/gothic/api/grpc/rpc/account"
	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc/codes"
)

func (s *accountServer) Signup(ctx context.Context,
	req *account.SignupRequest) (*api.UserResponse, error) {
	if req == nil {
		err := errors.New("request not found")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	rtx := rpc.RequestContext(ctx)
	rtx.SetCode(req.Code)
	rtx.SetProvider(s.Provider())
	s.Debugf("signup user: %v (%v)", req, rtx)
	u, err := s.API.Signup(rtx, req.Email, req.Username, req.Password, req.Data.AsMap())
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
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
