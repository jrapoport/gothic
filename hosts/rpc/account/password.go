package account

import (
	"context"
	"errors"

	"github.com/jrapoport/gothic/api/grpc/rpc"
	"github.com/jrapoport/gothic/api/grpc/rpc/account"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) SendResetPassword(ctx context.Context,
	req *account.ResetPasswordRequest) (*emptypb.Empty, error) {
	if req == nil {
		err := errors.New("request not found")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	if req.Email == "" {
		err := errors.New("email not found")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	email, err := s.ValidateEmail(req.Email)
	if err != nil {
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	u, err := s.GetUserWithEmail(email)
	if err != nil {
		return nil, s.RPCError(codes.InvalidArgument, err)

	}
	rtx := rpc.RequestContext(ctx)
	rtx.SetProvider(s.Provider())
	err = s.API.SendResetPassword(rtx, u.ID)
	if errors.Is(err, config.ErrRateLimitExceeded) {
		return nil, s.RPCError(codes.DeadlineExceeded, err)
	}
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *server) ConfirmResetPassword(ctx context.Context,
	req *account.ConfirmPasswordRequest) (*api.BearerResponse, error) {
	if req == nil {
		err := errors.New("request not found")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	if req.Password == "" {
		err := errors.New("password not found")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	if req.Token == "" {
		err := errors.New("token not found")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	s.Debugf("change password: %v", req)
	rtx := rpc.RequestContext(ctx)
	u, err := s.API.ConfirmResetPassword(rtx, req.Token, req.Password)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	bt, err := s.GrantBearerToken(rtx, u)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	s.Debugf("password changed: %s", bt.UserID)
	res := rpc.NewBearerResponse(bt)
	return res, nil
}
