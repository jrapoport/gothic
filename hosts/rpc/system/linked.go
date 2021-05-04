package system

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/api/grpc/rpc/system"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/models/account"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

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
	list := make([]*system.LinkedAccount, len(linked))
	for i, link := range linked {
		list[i] = NewAccount(link)
	}
	res := &system.LinkedAccountsResponse{
		Linked: list,
	}
	return res, nil
}

// NewAccount returns a new system account
func NewAccount(a *account.Account) *system.LinkedAccount {
	res := &system.LinkedAccount{
		Type:      uint32(a.Type),
		Provider:  a.Provider.String(),
		AccountId: a.AccountID,
		Email:     a.Email,
	}
	res.Data, _ = structpb.NewStruct(a.Data)
	return res
}
