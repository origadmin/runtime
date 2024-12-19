// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: security/jwt/v1/token.proto

package jwtv1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "github.com/google/gnostic/openapiv3"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/durationpb"
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

// PWT is a web token that can be used to authenticate a user with protobuf services.
type Token struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId       string                 `protobuf:"bytes,1,opt,name=client_id,proto3" json:"client_id,omitempty"`
	UserId         string                 `protobuf:"bytes,2,opt,name=user_id,proto3" json:"user_id,omitempty"`
	AccessToken    string                 `protobuf:"bytes,10,opt,name=access_token,proto3" json:"access_token,omitempty"`
	RefreshToken   string                 `protobuf:"bytes,11,opt,name=refresh_token,proto3" json:"refresh_token,omitempty"`
	ExpirationTime *timestamppb.Timestamp `protobuf:"bytes,12,opt,name=expiration_time,proto3" json:"expiration_time,omitempty"`
}

func (x *Token) Reset() {
	*x = Token{}
	mi := &file_security_jwt_v1_token_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Token) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Token) ProtoMessage() {}

func (x *Token) ProtoReflect() protoreflect.Message {
	mi := &file_security_jwt_v1_token_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Token.ProtoReflect.Descriptor instead.
func (*Token) Descriptor() ([]byte, []int) {
	return file_security_jwt_v1_token_proto_rawDescGZIP(), []int{0}
}

func (x *Token) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *Token) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Token) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *Token) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

func (x *Token) GetExpirationTime() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpirationTime
	}
	return nil
}

var File_security_jwt_v1_token_proto protoreflect.FileDescriptor

var file_security_jwt_v1_token_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2f, 0x6a, 0x77, 0x74, 0x2f, 0x76,
	0x31, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x73,
	0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2e, 0x6a, 0x77, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x24,
	0x67, 0x6e, 0x6f, 0x73, 0x74, 0x69, 0x63, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x2f,
	0x76, 0x33, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe3,
	0x03, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x4c, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x2e, 0xba, 0x47, 0x2b,
	0x92, 0x02, 0x28, 0x54, 0x68, 0x65, 0x20, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x20, 0x49, 0x44,
	0x20, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x64, 0x20, 0x77, 0x69, 0x74, 0x68,
	0x20, 0x74, 0x68, 0x65, 0x20, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x52, 0x09, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x12, 0x54, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x3a, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01,
	0xba, 0x47, 0x30, 0x92, 0x02, 0x2d, 0x54, 0x68, 0x65, 0x20, 0x49, 0x44, 0x20, 0x6f, 0x66, 0x20,
	0x74, 0x68, 0x65, 0x20, 0x75, 0x73, 0x65, 0x72, 0x20, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61,
	0x74, 0x65, 0x64, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x74, 0x68, 0x65, 0x20, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x2e, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x12, 0x5e, 0x0a, 0x0c,
	0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x3a, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0xba, 0x47, 0x30, 0x92, 0x02,
	0x2d, 0x54, 0x68, 0x65, 0x20, 0x77, 0x65, 0x62, 0x20, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x20,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x20, 0x75, 0x73, 0x65, 0x64, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x61,
	0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x52, 0x0c,
	0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x67, 0x0a, 0x0d,
	0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x41, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0xba, 0x47, 0x37, 0x92,
	0x02, 0x34, 0x54, 0x68, 0x65, 0x20, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x20, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x20, 0x75, 0x73, 0x65, 0x64, 0x20, 0x74, 0x6f, 0x20, 0x6f, 0x62, 0x74, 0x61,
	0x69, 0x6e, 0x20, 0x61, 0x20, 0x6e, 0x65, 0x77, 0x20, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x20,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x52, 0x0d, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x5f,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x6d, 0x0a, 0x0f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x27, 0xba, 0x47, 0x24, 0x92,
	0x02, 0x21, 0x54, 0x68, 0x65, 0x20, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x20, 0x74, 0x69, 0x6d, 0x65, 0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65, 0x20, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x2e, 0x52, 0x0f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x42, 0xbd, 0x01, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x65, 0x63,
	0x75, 0x72, 0x69, 0x74, 0x79, 0x2e, 0x6a, 0x77, 0x74, 0x2e, 0x76, 0x31, 0x42, 0x0a, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x39, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x72, 0x69, 0x67, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x2f, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f,
	0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2f, 0x6a, 0x77, 0x74, 0x2f, 0x76, 0x31, 0x3b,
	0x6a, 0x77, 0x74, 0x76, 0x31, 0xf8, 0x01, 0x01, 0xa2, 0x02, 0x03, 0x53, 0x4a, 0x58, 0xaa, 0x02,
	0x0f, 0x53, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2e, 0x4a, 0x77, 0x74, 0x2e, 0x56, 0x31,
	0xca, 0x02, 0x0f, 0x53, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x5c, 0x4a, 0x77, 0x74, 0x5c,
	0x56, 0x31, 0xe2, 0x02, 0x1b, 0x53, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x5c, 0x4a, 0x77,
	0x74, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x11, 0x53, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x3a, 0x3a, 0x4a, 0x77, 0x74,
	0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_security_jwt_v1_token_proto_rawDescOnce sync.Once
	file_security_jwt_v1_token_proto_rawDescData = file_security_jwt_v1_token_proto_rawDesc
)

func file_security_jwt_v1_token_proto_rawDescGZIP() []byte {
	file_security_jwt_v1_token_proto_rawDescOnce.Do(func() {
		file_security_jwt_v1_token_proto_rawDescData = protoimpl.X.CompressGZIP(file_security_jwt_v1_token_proto_rawDescData)
	})
	return file_security_jwt_v1_token_proto_rawDescData
}

var file_security_jwt_v1_token_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_security_jwt_v1_token_proto_goTypes = []any{
	(*Token)(nil),                 // 0: security.jwt.v1.Token
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_security_jwt_v1_token_proto_depIdxs = []int32{
	1, // 0: security.jwt.v1.Token.expiration_time:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_security_jwt_v1_token_proto_init() }
func file_security_jwt_v1_token_proto_init() {
	if File_security_jwt_v1_token_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_security_jwt_v1_token_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_security_jwt_v1_token_proto_goTypes,
		DependencyIndexes: file_security_jwt_v1_token_proto_depIdxs,
		MessageInfos:      file_security_jwt_v1_token_proto_msgTypes,
	}.Build()
	File_security_jwt_v1_token_proto = out.File
	file_security_jwt_v1_token_proto_rawDesc = nil
	file_security_jwt_v1_token_proto_goTypes = nil
	file_security_jwt_v1_token_proto_depIdxs = nil
}
