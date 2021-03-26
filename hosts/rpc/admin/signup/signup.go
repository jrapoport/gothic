package signup

import (
	"context"
	"errors"

	"github.com/jrapoport/gothic/api/grpc/rpc/admin/signup"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/models/code"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type signupServer struct {
	signup.UnimplementedSignupServer
	*rpc.Server
}

var _ signup.SignupServer = (*signupServer)(nil)

func newSignupServer(srv *rpc.Server) *signupServer {
	srv.FieldLogger = srv.WithField("module", "codes")
	return &signupServer{Server: srv}
}

// RegisterServer registers a new admin server.
func RegisterServer(s *grpc.Server, srv *rpc.Server) {
	signup.RegisterSignupServer(s, newSignupServer(srv))
}

// CreateSignupCodes returns the settings for a server.
func (s *signupServer) CreateSignupCodes(ctx context.Context,
	req *signup.CreateSignupCodesRequest) (*signup.CreateSignupCodesResponse, error) {
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
	res := &signup.CreateSignupCodesResponse{
		Codes: list,
	}
	return res, nil
}

func (s *signupServer) CheckSignupCode(_ context.Context,
	req *signup.CheckSignupCodeRequest) (*signup.CheckSignupCodeResponse, error) {
	if req == nil {
		return nil, s.RPCError(codes.InvalidArgument, nil)
	}
	cd := req.GetCode()
	s.Debugf("check signup code: %s", cd)
	sc, err := s.API.CheckSignupCode(cd)
	if err != nil && errors.Is(err, code.ErrUnusableCode) {
		s.Debugf("checked signup code is invalid: %s", cd)
		res := &signup.CheckSignupCodeResponse{
			Usable: false,
			Code:   cd,
		}
		return res, nil
	} else if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	s.Debugf("checked signup code is valid: %s", cd)
	res := &signup.CheckSignupCodeResponse{
		Usable:     true,
		Code:       sc.Token,
		CodeFormat: signup.CodeFormat(sc.Format),
		CodeType:   signup.CodeUsage(sc.Usage()),
		Expiration: durationpb.New(sc.Expiration),
		UserId:     sc.UserID.String(),
	}
	return res, nil
}

func (s *signupServer) VoidSignupCode(_ context.Context,
	req *signup.VoidSignupCodeRequest) (*emptypb.Empty, error) {
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
