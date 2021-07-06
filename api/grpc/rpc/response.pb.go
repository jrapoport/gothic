// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.0
// 	protoc        v3.17.3
// source: response.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   string           `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Role     string           `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
	Email    string           `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Username string           `protobuf:"bytes,4,opt,name=username,proto3" json:"username,omitempty"`
	Data     *structpb.Struct `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	Token    *BearerResponse  `protobuf:"bytes,6,opt,name=token,proto3,oneof" json:"token,omitempty"`
}

func (x *UserResponse) Reset() {
	*x = UserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_response_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserResponse) ProtoMessage() {}

func (x *UserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_response_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserResponse.ProtoReflect.Descriptor instead.
func (*UserResponse) Descriptor() ([]byte, []int) {
	return file_response_proto_rawDescGZIP(), []int{0}
}

func (x *UserResponse) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UserResponse) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *UserResponse) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *UserResponse) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *UserResponse) GetData() *structpb.Struct {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *UserResponse) GetToken() *BearerResponse {
	if x != nil {
		return x.Token
	}
	return nil
}

type BearerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      string                 `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Access    string                 `protobuf:"bytes,2,opt,name=access,proto3" json:"access,omitempty"`
	Refresh   string                 `protobuf:"bytes,3,opt,name=refresh,proto3" json:"refresh,omitempty"`
	Id        string                 `protobuf:"bytes,4,opt,name=id,proto3" json:"id,omitempty"`
	ExpiresAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
}

func (x *BearerResponse) Reset() {
	*x = BearerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_response_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BearerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BearerResponse) ProtoMessage() {}

func (x *BearerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_response_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BearerResponse.ProtoReflect.Descriptor instead.
func (*BearerResponse) Descriptor() ([]byte, []int) {
	return file_response_proto_rawDescGZIP(), []int{1}
}

func (x *BearerResponse) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *BearerResponse) GetAccess() string {
	if x != nil {
		return x.Access
	}
	return ""
}

func (x *BearerResponse) GetRefresh() string {
	if x != nil {
		return x.Refresh
	}
	return ""
}

func (x *BearerResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *BearerResponse) GetExpiresAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpiresAt
	}
	return nil
}

type PagedResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index int64  `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Size  int64  `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	First int64  `protobuf:"varint,3,opt,name=first,proto3" json:"first,omitempty"`
	Prev  int64  `protobuf:"varint,4,opt,name=prev,proto3" json:"prev,omitempty"`
	Next  int64  `protobuf:"varint,5,opt,name=next,proto3" json:"next,omitempty"`
	Last  int64  `protobuf:"varint,6,opt,name=last,proto3" json:"last,omitempty"`
	Count int64  `protobuf:"varint,7,opt,name=count,proto3" json:"count,omitempty"`
	Total uint64 `protobuf:"varint,8,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *PagedResponse) Reset() {
	*x = PagedResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_response_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PagedResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PagedResponse) ProtoMessage() {}

func (x *PagedResponse) ProtoReflect() protoreflect.Message {
	mi := &file_response_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PagedResponse.ProtoReflect.Descriptor instead.
func (*PagedResponse) Descriptor() ([]byte, []int) {
	return file_response_proto_rawDescGZIP(), []int{2}
}

func (x *PagedResponse) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *PagedResponse) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *PagedResponse) GetFirst() int64 {
	if x != nil {
		return x.First
	}
	return 0
}

func (x *PagedResponse) GetPrev() int64 {
	if x != nil {
		return x.Prev
	}
	return 0
}

func (x *PagedResponse) GetNext() int64 {
	if x != nil {
		return x.Next
	}
	return 0
}

func (x *PagedResponse) GetLast() int64 {
	if x != nil {
		return x.Last
	}
	return 0
}

func (x *PagedResponse) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *PagedResponse) GetTotal() uint64 {
	if x != nil {
		return x.Total
	}
	return 0
}

var File_response_proto protoreflect.FileDescriptor

var file_response_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0a, 0x67, 0x6f, 0x74, 0x68, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x69, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74,
	0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdb, 0x01, 0x0a, 0x0c,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12,
	0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75,
	0x63, 0x74, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x35, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x74, 0x68, 0x69, 0x63,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x42, 0x65, 0x61, 0x72, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x88, 0x01, 0x01, 0x42,
	0x08, 0x0a, 0x06, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0xa1, 0x01, 0x0a, 0x0e, 0x42, 0x65,
	0x61, 0x72, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x66, 0x72,
	0x65, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x72, 0x65, 0x66, 0x72, 0x65,
	0x73, 0x68, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x61, 0x74,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x41, 0x74, 0x22, 0xb7, 0x01,
	0x0a, 0x0d, 0x50, 0x61, 0x67, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x69, 0x72,
	0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x66, 0x69, 0x72, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x70, 0x72, 0x65, 0x76, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x70,
	0x72, 0x65, 0x76, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x61, 0x73, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x6c, 0x61, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x72, 0x61, 0x70, 0x6f, 0x70, 0x6f, 0x72, 0x74, 0x2f,
	0x67, 0x6f, 0x74, 0x68, 0x69, 0x63, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f,
	0x72, 0x70, 0x63, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_response_proto_rawDescOnce sync.Once
	file_response_proto_rawDescData = file_response_proto_rawDesc
)

func file_response_proto_rawDescGZIP() []byte {
	file_response_proto_rawDescOnce.Do(func() {
		file_response_proto_rawDescData = protoimpl.X.CompressGZIP(file_response_proto_rawDescData)
	})
	return file_response_proto_rawDescData
}

var file_response_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_response_proto_goTypes = []interface{}{
	(*UserResponse)(nil),          // 0: gothic.api.UserResponse
	(*BearerResponse)(nil),        // 1: gothic.api.BearerResponse
	(*PagedResponse)(nil),         // 2: gothic.api.PagedResponse
	(*structpb.Struct)(nil),       // 3: google.protobuf.Struct
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_response_proto_depIdxs = []int32{
	3, // 0: gothic.api.UserResponse.data:type_name -> google.protobuf.Struct
	1, // 1: gothic.api.UserResponse.token:type_name -> gothic.api.BearerResponse
	4, // 2: gothic.api.BearerResponse.expires_at:type_name -> google.protobuf.Timestamp
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_response_proto_init() }
func file_response_proto_init() {
	if File_response_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_response_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_response_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BearerResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_response_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PagedResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_response_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_response_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_response_proto_goTypes,
		DependencyIndexes: file_response_proto_depIdxs,
		MessageInfos:      file_response_proto_msgTypes,
	}.Build()
	File_response_proto = out.File
	file_response_proto_rawDesc = nil
	file_response_proto_goTypes = nil
	file_response_proto_depIdxs = nil
}
