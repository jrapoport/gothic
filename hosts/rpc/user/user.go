package user

//go:generate protoc -I=. -I=.. --go_out=plugins=grpc:. --go_opt=paths=source_relative user.proto

import (
	"context"

	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type userServer struct {
	*rpc.Server
}

var _ UserServer = (*userServer)(nil)

func newUserServer(srv *rpc.Server) *userServer {
	srv.FieldLogger = srv.WithField("module", "user")
	return &userServer{srv}
}

// RegisterServer registers a new admin server.
func RegisterServer(s *grpc.Server, srv *rpc.Server) {
	RegisterUserServer(s, newUserServer(srv))
}

func (s *userServer) GetUser(ctx context.Context, _ *GetUserRequest) (*rpc.UserResponse, error) {
	uid, err := rpc.GetUserID(ctx)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	u, err := s.GetAuthenticatedUser(uid)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	s.Debugf("get user %s", u.ID)
	res, err := rpc.NewUserResponse(u)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	if s.Config().MaskEmails {
		res.MaskEmail()
	}
	s.Debugf("got user %s: %v", uid, res)
	return res, nil
}

func (s *userServer) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*rpc.UserResponse, error) {
	uid, err := rpc.GetUserID(ctx)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	u, err := s.GetAuthenticatedUser(uid)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	rtx := rpc.RequestContext(ctx)
	s.Debugf("update user %s: %v", uid.String(), req)
	u, err = s.API.UpdateUser(rtx, uid, &req.Username, req.Data.AsMap())
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
	s.Debugf("got user %s: %v", uid, res)
	return res, nil
}

func (s *userServer) ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*rpc.BearerResponse, error) {
	uid, err := rpc.GetUserID(ctx)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	u, err := s.GetAuthenticatedUser(uid)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	rtx := rpc.RequestContext(ctx)
	s.Debugf("change password %s: %v", u.ID, req)
	u, err = s.API.ChangePassword(rtx, u.ID, req.OldPassword, req.Password)
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
