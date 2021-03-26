package account

import (
	"context"
	"errors"

	"github.com/jrapoport/gothic/api/grpc/rpc"
	"github.com/jrapoport/gothic/api/grpc/rpc/account"
	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc/codes"
)

func (s *accountServer) RefreshBearerToken(ctx context.Context,
	req *account.RefreshTokenRequest) (*api.BearerResponse, error) {
	if req == nil {
		err := errors.New("request not found")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	if req.Token == "" {
		err := errors.New("token not found")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	s.Debugf("refresh token: %v", req)
	rtx := rpc.RequestContext(ctx)
	bt, err := s.API.RefreshBearerToken(rtx, req.Token)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	s.Debugf("password changed: %s", bt.UserID)
	res := rpc.NewBearerResponse(bt)
	return res, nil
}
