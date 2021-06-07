package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/models/user"
	"google.golang.org/grpc/codes"
)

func (s *adminServer) CreateUser(ctx context.Context, req *admin.CreateUserRequest) (*admin.CreateUserResponse, error) {
	if req == nil {
		return nil, s.RPCError(codes.InvalidArgument, nil)
	}
	email := req.GetEmail()
	username := req.GetUsername()
	pw := req.GetPassword()
	data := req.GetData().AsMap()
	adminRole := req.GetAdmin()
	if email == "" {
		err := errors.New("email required")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	rtx, err := s.adminRequestContext(ctx)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	role := adminRoleFromContext(rtx)
	if adminRole && role != user.RoleSuper {
		err = fmt.Errorf("super admin access required to create admin users")
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	// we don't need to validate the admin here since CreateUser will do it.
	u, err := s.API.CreateUser(rtx, email, username, pw, data, adminRole)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	res := &admin.CreateUserResponse{
		Role:   u.Role.String(),
		UserId: u.ID.String(),
		Email:  u.Email,
	}
	return res, nil
}

func (s *adminServer) ChangeUserRole(ctx context.Context, req *admin.ChangeUserRoleRequest) (*admin.ChangeUserRoleResponse, error) {
	if req == nil {
		return nil, s.RPCError(codes.InvalidArgument, nil)
	}
	if req.GetUserId() == "" && req.GetEmail() == "" {
		err := errors.New("user id or email is required")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	role := user.ToRole(req.GetRole())
	if role == user.InvalidRole {
		err := errors.New("invalid user role")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	rtx, err := s.adminRequestContext(ctx)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	adminRole := adminRoleFromContext(rtx)
	if adminRole == user.RoleAdmin && role == user.RoleAdmin {
		err = errors.New("super admin user required")
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	var uid uuid.UUID
	switch req.GetUser().(type) {
	case *admin.ChangeUserRoleRequest_UserId:
		userID := req.GetUserId()
		id, err := uuid.Parse(userID)
		if err != nil || id == uuid.Nil {
			err = fmt.Errorf("invalid user id '%s': %w", userID, err)
			return nil, s.RPCError(codes.InvalidArgument, err)
		}
		uid = id
	case *admin.ChangeUserRoleRequest_Email:
		email := req.GetEmail()
		u, err := s.API.GetUserWithEmail(email)
		if err != nil {
			err = fmt.Errorf("user not found '%s': %w", email, err)
			return nil, s.RPCError(codes.InvalidArgument, err)
		}
		uid = u.ID
	}
	u, err := s.API.ChangeRole(rtx, uid, role)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	res := &admin.ChangeUserRoleResponse{
		UserId: u.ID.String(),
		Role:   u.Role.String(),
	}
	return res, nil
}

func (s *adminServer) DeleteUser(ctx context.Context, req *admin.DeleteUserRequest) (*admin.DeleteUserResponse, error) {
	if req == nil {
		return nil, s.RPCError(codes.InvalidArgument, nil)
	}
	if req.GetUserId() == "" && req.GetEmail() == "" {
		err := errors.New("user id or email is required")
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	rtx, err := s.adminRequestContext(ctx)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	var uid uuid.UUID
	switch req.GetUser().(type) {
	case *admin.DeleteUserRequest_UserId:
		userID := req.GetUserId()
		id, err := uuid.Parse(userID)
		if err != nil || id == uuid.Nil {
			err = fmt.Errorf("invalid user id '%s': %w", userID, err)
			return nil, s.RPCError(codes.InvalidArgument, err)
		}
		uid = id
	case *admin.DeleteUserRequest_Email:
		email := req.GetEmail()
		u, err := s.API.GetUserWithEmail(email)
		if err != nil {
			err = fmt.Errorf("user not found '%s': %w", email, err)
			return nil, s.RPCError(codes.InvalidArgument, err)
		}
		uid = u.ID
	}
	err = s.API.DeleteUser(rtx, uid, req.GetHard())
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	res := &admin.DeleteUserResponse{
		UserId: uid.String(),
	}
	return res, nil
}
