package codes

//go:generate protoc -I=. --go_out=plugins=grpc:. --go_opt=paths=source_relative codes.proto

import (
	"context"
	"errors"
	"github.com/jrapoport/gothic/models/code"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/durationpb"
)

type codesServer struct {
	*rpc.Server
}

var _ CodesServer = (*codesServer)(nil)

func newCodesServer(srv *rpc.Server) *codesServer {
	srv.FieldLogger = srv.WithField("module", "codes")
	return &codesServer{srv}
}

// RegisterServer registers a new admin server.
func RegisterServer(s *grpc.Server, srv *rpc.Server) {
	RegisterCodesServer(s, newCodesServer(srv))
}

// CreateSignupCodes returns the settings for a server.
func (s *codesServer) CreateSignupCodes(ctx context.Context, req *CreateSignupCodesRequest) (*CreateSignupCodesResponse, error) {
	if req == nil {
		return nil, s.RPCError(codes.InvalidArgument, nil)
	}
	rtx := rpc.RequestContext(ctx)
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
	res := &CreateSignupCodesResponse{
		Codes: list,
	}
	return res, nil
}

func (s *codesServer) CheckSignupCode(_ context.Context, req *CheckSignupCodeRequest) (*CheckSignupCodeResponse, error) {
	if req == nil {
		return nil, s.RPCError(codes.InvalidArgument, nil)
	}
	cd := req.GetCode()
	s.Debugf("check signup code: %s", cd)
	sc, err := s.API.CheckSignupCode(cd)
	if err != nil && errors.Is(err, code.ErrUnusableCode) {
		s.Debugf("checked signup code is invalid: %s", cd)
		res := &CheckSignupCodeResponse{
			Usable: false,
			Code:   cd,
		}
		return res, nil
	} else if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	s.Debugf("checked signup code is valid: %s", cd)
	res := &CheckSignupCodeResponse{
		Usable:     true,
		Code:       sc.Token,
		CodeFormat: CodeFormat(sc.Format),
		CodeType:   CodeUsage(sc.Usage()),
		Expiration: durationpb.New(sc.Expiration),
		UserId:     sc.UserID.String(),
	}
	return res, nil
}

func (s *codesServer) VoidSignupCode(_ context.Context, req *VoidSignupCodeRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, s.RPCError(codes.InvalidArgument, nil)
	}
	cd := req.GetCode()
	s.Debugf("void signup code: %s", cd)
	err := s.API.VoidSignupCode(cd)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	s.Debugf("voided signup code: %s", cd)
	return &emptypb.Empty{}, nil
}
