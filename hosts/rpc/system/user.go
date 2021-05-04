package system

import (
	"context"
	"net/mail"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/api/grpc/rpc/system"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/models/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *systemServer) GetUserAccount(_ context.Context, req *system.UserAccountRequest) (*system.UserAccountResponse, error) {
	var u *user.User
	switch msg := req.GetId().(type) {
	case *system.UserAccountRequest_UserId:
		s.Debugf("get user %s", msg.UserId)
		uid, err := uuid.Parse(msg.UserId)
		if err != nil {
			return nil, s.RPCError(codes.InvalidArgument, err)
		}
		u, err = s.API.GetUser(uid)
		if err != nil {
			return nil, s.RPCError(codes.InvalidArgument, err)
		}
	case *system.UserAccountRequest_Email:
		s.Debugf("get user %s", req.GetEmail())
		addr, err := mail.ParseAddress(req.GetEmail())
		if err != nil {
			return nil, s.RPCError(codes.Internal, err)
		}
		u, err = s.API.GetUserWithEmail(addr.Address)
		if err != nil {
			return nil, s.RPCError(codes.Internal, err)
		}
	default:
		return nil, s.RPCError(codes.InvalidArgument, nil)
	}
	if u == nil { // this should never happen, bu jic
		return nil, s.RPCError(codes.Internal, nil)
	}
	res := NewUserResponse(u)
	s.Debugf("got user %s: %v", u.ID, res)
	return res, nil
}

// NewUserResponse returns a UserResponse for a user.
func NewUserResponse(u *user.User) *system.UserAccountResponse {
	ur := &system.UserAccountResponse{
		Id:        u.ID.String(),
		Provider:  u.Provider.String(),
		Role:      u.Role.String(),
		Status:    system.UserAccountResponse_Status(u.Status),
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: timestamppb.New(u.CreatedAt),
	}
	ur.Data, _ = structpb.NewStruct(u.Data)
	ur.Metadata, _ = structpb.NewStruct(u.Metadata)
	if u.ConfirmedAt != nil {
		ur.ConfirmedAt = timestamppb.New(*u.ConfirmedAt)
	}
	if u.VerifiedAt != nil {
		ur.VerifiedAt = timestamppb.New(*u.VerifiedAt)
	}
	return ur
}

func (s *systemServer) NotifyUser(ctx context.Context, req *system.NotificationRequest) (*system.NotificationResponse, error) {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	logo := req.GetLogo()
	sub := req.GetSubject()
	html := req.GetHtml()
	plain := req.GetPlain()
	rtx := rpc.RequestContext(ctx)
	sent, err := s.API.NotifyUser(rtx, uid, logo, sub, html, plain)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	res := &system.NotificationResponse{
		Sent: sent,
	}
	return res, nil
}
