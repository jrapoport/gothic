package admin

import (
	"context"
	"errors"

	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/models/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateSignupCodes returns the settings for a server.
func (s *adminServer) CreateSignupCodes(ctx context.Context,
	req *admin.CreateSignupCodesRequest) (*admin.SignupCodesResponse, error) {
	if req == nil {
		return nil, s.RPCError(codes.InvalidArgument, nil)
	}
	rtx, err := s.adminRequestContext(ctx)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	uses := int(req.GetUses())
	count := int(req.GetCount())
	s.Debugf("create signup codes: %v", req)
	list, err := s.API.CreateSignupCodes(rtx, uses, count)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	if len(list) < count {
		s.Warnf("expected %d but created %d", count, len(list))
	}
	s.Debugf("created %d codes", len(list))
	res := &admin.SignupCodesResponse{
		Codes: list,
	}
	return res, nil
}

func (s *adminServer) CheckSignupCode(ctx context.Context,
	req *admin.CheckSignupCodeRequest) (*admin.SignupCodeResponse, error) {
	if req == nil {
		return nil, s.RPCError(codes.InvalidArgument, nil)
	}
	_, err := s.adminRequestContext(ctx)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	cd := req.GetCode()
	s.Debugf("check signup code: %s", cd)
	sc, err := s.API.CheckSignupCode(cd)
	if err != nil && errors.Is(err, code.ErrUnusableCode) {
		s.Debugf("checked signup code is invalid: %s", cd)
		res := &admin.SignupCodeResponse{
			Valid: false,
			Code:  cd,
		}
		return res, nil
	} else if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	s.Debugf("checked signup code is valid: %s", cd)
	res := &admin.SignupCodeResponse{
		Valid:      true,
		Code:       sc.Token,
		Format:     admin.CodeFormat(sc.Format),
		Type:       admin.CodeType(sc.Type),
		Expiration: durationpb.New(sc.Expiration),
		UserId:     sc.UserID.String(),
	}
	return res, nil
}

func (s *adminServer) DeleteSignupCode(ctx context.Context, req *admin.DeleteSignupCodeRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, s.RPCError(codes.InvalidArgument, nil)
	}
	_, err := s.adminRequestContext(ctx)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	cd := req.GetCode()
	s.Debugf("delete signup code: %s", cd)
	err = s.API.DeleteSignupCode(cd)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	s.Debugf("deleted signup code: %s", cd)
	return &emptypb.Empty{}, nil
}
