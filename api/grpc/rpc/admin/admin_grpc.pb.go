// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package admin

import (
	context "context"
	rpc "github.com/jrapoport/gothic/api/grpc/rpc"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AdminClient is the client API for Admin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdminClient interface {
	CreateSignupCodes(ctx context.Context, in *CreateSignupCodesRequest, opts ...grpc.CallOption) (*SignupCodesResponse, error)
	CheckSignupCode(ctx context.Context, in *CheckSignupCodeRequest, opts ...grpc.CallOption) (*SignupCodeResponse, error)
	DeleteSignupCode(ctx context.Context, in *DeleteSignupCodeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error)
	UpdateUserMetadata(ctx context.Context, in *UpdateUserMetadataRequest, opts ...grpc.CallOption) (*UpdateUserMetadataResponse, error)
	ChangeUserRole(ctx context.Context, in *ChangeUserRoleRequest, opts ...grpc.CallOption) (*ChangeUserRoleResponse, error)
	SearchAuditLogs(ctx context.Context, in *rpc.SearchRequest, opts ...grpc.CallOption) (*AuditLogsResult, error)
	Settings(ctx context.Context, in *SettingsRequest, opts ...grpc.CallOption) (*SettingsResponse, error)
}

type adminClient struct {
	cc grpc.ClientConnInterface
}

func NewAdminClient(cc grpc.ClientConnInterface) AdminClient {
	return &adminClient{cc}
}

func (c *adminClient) CreateSignupCodes(ctx context.Context, in *CreateSignupCodesRequest, opts ...grpc.CallOption) (*SignupCodesResponse, error) {
	out := new(SignupCodesResponse)
	err := c.cc.Invoke(ctx, "/gothic.api.Admin/CreateSignupCodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) CheckSignupCode(ctx context.Context, in *CheckSignupCodeRequest, opts ...grpc.CallOption) (*SignupCodeResponse, error) {
	out := new(SignupCodeResponse)
	err := c.cc.Invoke(ctx, "/gothic.api.Admin/CheckSignupCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) DeleteSignupCode(ctx context.Context, in *DeleteSignupCodeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/gothic.api.Admin/DeleteSignupCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, "/gothic.api.Admin/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error) {
	out := new(DeleteUserResponse)
	err := c.cc.Invoke(ctx, "/gothic.api.Admin/DeleteUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) UpdateUserMetadata(ctx context.Context, in *UpdateUserMetadataRequest, opts ...grpc.CallOption) (*UpdateUserMetadataResponse, error) {
	out := new(UpdateUserMetadataResponse)
	err := c.cc.Invoke(ctx, "/gothic.api.Admin/UpdateUserMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) ChangeUserRole(ctx context.Context, in *ChangeUserRoleRequest, opts ...grpc.CallOption) (*ChangeUserRoleResponse, error) {
	out := new(ChangeUserRoleResponse)
	err := c.cc.Invoke(ctx, "/gothic.api.Admin/ChangeUserRole", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) SearchAuditLogs(ctx context.Context, in *rpc.SearchRequest, opts ...grpc.CallOption) (*AuditLogsResult, error) {
	out := new(AuditLogsResult)
	err := c.cc.Invoke(ctx, "/gothic.api.Admin/SearchAuditLogs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) Settings(ctx context.Context, in *SettingsRequest, opts ...grpc.CallOption) (*SettingsResponse, error) {
	out := new(SettingsResponse)
	err := c.cc.Invoke(ctx, "/gothic.api.Admin/Settings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdminServer is the server API for Admin service.
// All implementations must embed UnimplementedAdminServer
// for forward compatibility
type AdminServer interface {
	CreateSignupCodes(context.Context, *CreateSignupCodesRequest) (*SignupCodesResponse, error)
	CheckSignupCode(context.Context, *CheckSignupCodeRequest) (*SignupCodeResponse, error)
	DeleteSignupCode(context.Context, *DeleteSignupCodeRequest) (*emptypb.Empty, error)
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error)
	UpdateUserMetadata(context.Context, *UpdateUserMetadataRequest) (*UpdateUserMetadataResponse, error)
	ChangeUserRole(context.Context, *ChangeUserRoleRequest) (*ChangeUserRoleResponse, error)
	SearchAuditLogs(context.Context, *rpc.SearchRequest) (*AuditLogsResult, error)
	Settings(context.Context, *SettingsRequest) (*SettingsResponse, error)
	mustEmbedUnimplementedAdminServer()
}

// UnimplementedAdminServer must be embedded to have forward compatible implementations.
type UnimplementedAdminServer struct {
}

func (UnimplementedAdminServer) CreateSignupCodes(context.Context, *CreateSignupCodesRequest) (*SignupCodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSignupCodes not implemented")
}
func (UnimplementedAdminServer) CheckSignupCode(context.Context, *CheckSignupCodeRequest) (*SignupCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckSignupCode not implemented")
}
func (UnimplementedAdminServer) DeleteSignupCode(context.Context, *DeleteSignupCodeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSignupCode not implemented")
}
func (UnimplementedAdminServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedAdminServer) DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedAdminServer) UpdateUserMetadata(context.Context, *UpdateUserMetadataRequest) (*UpdateUserMetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserMetadata not implemented")
}
func (UnimplementedAdminServer) ChangeUserRole(context.Context, *ChangeUserRoleRequest) (*ChangeUserRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeUserRole not implemented")
}
func (UnimplementedAdminServer) SearchAuditLogs(context.Context, *rpc.SearchRequest) (*AuditLogsResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchAuditLogs not implemented")
}
func (UnimplementedAdminServer) Settings(context.Context, *SettingsRequest) (*SettingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Settings not implemented")
}
func (UnimplementedAdminServer) mustEmbedUnimplementedAdminServer() {}

// UnsafeAdminServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdminServer will
// result in compilation errors.
type UnsafeAdminServer interface {
	mustEmbedUnimplementedAdminServer()
}

func RegisterAdminServer(s grpc.ServiceRegistrar, srv AdminServer) {
	s.RegisterService(&Admin_ServiceDesc, srv)
}

func _Admin_CreateSignupCodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSignupCodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).CreateSignupCodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gothic.api.Admin/CreateSignupCodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).CreateSignupCodes(ctx, req.(*CreateSignupCodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_CheckSignupCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckSignupCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).CheckSignupCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gothic.api.Admin/CheckSignupCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).CheckSignupCode(ctx, req.(*CheckSignupCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_DeleteSignupCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSignupCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).DeleteSignupCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gothic.api.Admin/DeleteSignupCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).DeleteSignupCode(ctx, req.(*DeleteSignupCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gothic.api.Admin/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gothic.api.Admin/DeleteUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_UpdateUserMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).UpdateUserMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gothic.api.Admin/UpdateUserMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).UpdateUserMetadata(ctx, req.(*UpdateUserMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_ChangeUserRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeUserRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).ChangeUserRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gothic.api.Admin/ChangeUserRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).ChangeUserRole(ctx, req.(*ChangeUserRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_SearchAuditLogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(rpc.SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).SearchAuditLogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gothic.api.Admin/SearchAuditLogs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).SearchAuditLogs(ctx, req.(*rpc.SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_Settings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SettingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).Settings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gothic.api.Admin/Settings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).Settings(ctx, req.(*SettingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Admin_ServiceDesc is the grpc.ServiceDesc for Admin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Admin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gothic.api.Admin",
	HandlerType: (*AdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateSignupCodes",
			Handler:    _Admin_CreateSignupCodes_Handler,
		},
		{
			MethodName: "CheckSignupCode",
			Handler:    _Admin_CheckSignupCode_Handler,
		},
		{
			MethodName: "DeleteSignupCode",
			Handler:    _Admin_DeleteSignupCode_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _Admin_CreateUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _Admin_DeleteUser_Handler,
		},
		{
			MethodName: "UpdateUserMetadata",
			Handler:    _Admin_UpdateUserMetadata_Handler,
		},
		{
			MethodName: "ChangeUserRole",
			Handler:    _Admin_ChangeUserRole_Handler,
		},
		{
			MethodName: "SearchAuditLogs",
			Handler:    _Admin_SearchAuditLogs_Handler,
		},
		{
			MethodName: "Settings",
			Handler:    _Admin_Settings_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "admin.proto",
}
