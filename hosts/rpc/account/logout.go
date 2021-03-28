package account

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/api/grpc/rpc/account"
	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *accountServer) Logout(ctx context.Context,
	_ *account.LogoutRequest) (*emptypb.Empty, error) {
	rtx := rpc.RequestContext(ctx)
	uid := rtx.UserID()
	if uid == uuid.Nil {
		err := errors.New("invalid user id")
		return nil, s.RPCError(codes.Unauthenticated, err)
	}
	s.Debugf("logout user: %s (%v)", uid, rtx.Provider())
	err := s.API.Logout(rtx, uid)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	return &emptypb.Empty{}, nil
}
