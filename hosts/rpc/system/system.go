package system

import (
	"context"
	"errors"
	"net/mail"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/api/grpc/rpc/system"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/models/account"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type systemServer struct {
	system.UnimplementedSystemServer
	*rpc.Server
}

var _ system.SystemServer = (*systemServer)(nil)

func newSystemServer(srv *rpc.Server) *systemServer {
	srv.FieldLogger = srv.WithField("module", "user")
	return &systemServer{Server: srv}
}

// RegisterServer registers a new admin server.
func RegisterServer(s *grpc.Server, srv *rpc.Server) {
	system.RegisterSystemServer(s, newSystemServer(srv))
}

func (s *systemServer) GetUser(_ context.Context, req *system.UserRequest) (*system.UserResponse, error) {
	var u *user.User
	switch msg := req.GetId().(type) {
	case *system.UserRequest_UserId:
		s.Debugf("get user %s", msg.UserId)
		uid, err := uuid.Parse(msg.UserId)
		if err != nil {
			return nil, s.RPCError(codes.InvalidArgument, err)
		}
		u, err = s.API.GetUser(uid)
		if err != nil {
			return nil, s.RPCError(codes.InvalidArgument, err)
		}
	case *system.UserRequest_Email:
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

func (s *systemServer) LinkAccount(ctx context.Context,
	req *system.LinkAccountRequest) (*emptypb.Empty, error) {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	act := req.Account
	t := account.Type(act.GetType())
	if !t.Has(account.All) {
		err = errors.New("invalid type")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	p := provider.Name(act.GetProvider())
	if p == provider.Unknown {
		err = errors.New("invalid provider")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	aid := act.GetAccountId()
	if aid == "" {
		err = errors.New("invalid account id")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	l := &account.Account{
		Type:      t,
		Provider:  p,
		AccountID: aid,
		Email:     act.Email,
	}
	if act.Data != nil {
		l.Data = act.Data.AsMap()
	}
	rtx := rpc.RequestContext(ctx)
	err = s.API.LinkAccount(rtx, uid, l)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *systemServer) GetLinkedAccounts(ctx context.Context,
	req *system.LinkedAccountsRequest) (*system.LinkedAccountsResponse, error) {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	filters := store.FiltersFromMap(req.GetFilters())
	p := provider.Name(req.GetProvider())
	if p != provider.Unknown {
		filters[key.Provider] = p
	}
	t := account.Type(req.GetType())
	rtx := rpc.RequestContext(ctx)
	linked, err := s.API.GetLinkedAccounts(rtx, uid, t, filters)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	list := make([]*system.Account, len(linked))
	for i, link := range linked {
		list[i] = NewAccount(link)
	}
	res := &system.LinkedAccountsResponse{
		Linked: list,
	}
	return res, nil
}

// NewUserResponse returns a UserResponse for a user.
func NewUserResponse(u *user.User) *system.UserResponse {
	ur := &system.UserResponse{
		Id:        u.ID.String(),
		Provider:  u.Provider.String(),
		Role:      u.Role.String(),
		Status:    system.UserResponse_Status(u.Status),
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

func NewAccount(a *account.Account) *system.Account {
	res := &system.Account{
		Type:      uint32(a.Type),
		Provider:  a.Provider.String(),
		AccountId: a.AccountID,
		Email:     a.Email,
	}
	res.Data, _ = structpb.NewStruct(a.Data)
	return res
}
