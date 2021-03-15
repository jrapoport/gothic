package account

import (
	"context"
	"errors"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *accountServer) SendConfirmUser(ctx context.Context,
	req *SendConfirmRequest) (*emptypb.Empty, error) {
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
	err = s.API.SendConfirmUser(rtx, u.ID)
	if errors.Is(err, config.ErrRateLimitExceeded) {
		return nil, s.RPCError(codes.DeadlineExceeded, err)
	}
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *accountServer) ConfirmUser(ctx context.Context,
	req *ConfirmUserRequest) (*rpc.BearerResponse, error) {
	if req == nil {
		err := errors.New("request not found")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	if req.Token == "" {
		err := errors.New("token not found")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	s.Debugf("confirm user: %v", req)
	rtx := rpc.RequestContext(ctx)
	u, err := s.API.ConfirmUser(rtx, req.Token)
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
