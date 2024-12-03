// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: pwt/v1/pwt.proto

package pwtv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
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
type PWT struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExpirationTime *durationpb.Duration `protobuf:"bytes,1,opt,name=expiration_time,proto3" json:"expiration_time,omitempty"`
	IssuedAt       *durationpb.Duration `protobuf:"bytes,2,opt,name=issued_at,proto3" json:"issued_at,omitempty"`
	NotBefore      *durationpb.Duration `protobuf:"bytes,3,opt,name=not_before,proto3" json:"not_before,omitempty"`
	Issuer         string               `protobuf:"bytes,4,opt,name=issuer,proto3" json:"issuer,omitempty"`
	Audience       []string             `protobuf:"bytes,5,rep,name=audience,proto3" json:"audience,omitempty"`
	Subject        string               `protobuf:"bytes,6,opt,name=subject,proto3" json:"subject,omitempty"`
	JwtId          string               `protobuf:"bytes,7,opt,name=jwt_id,proto3" json:"jwt_id,omitempty"`
	ClientId       string               `protobuf:"bytes,8,opt,name=client_id,proto3" json:"client_id,omitempty"`
	ClientSecret   string               `protobuf:"bytes,9,opt,name=client_secret,proto3" json:"client_secret,omitempty"`
}

func (x *PWT) Reset() {
	*x = PWT{}
	mi := &file_pwt_v1_pwt_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PWT) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PWT) ProtoMessage() {}

func (x *PWT) ProtoReflect() protoreflect.Message {
	mi := &file_pwt_v1_pwt_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PWT.ProtoReflect.Descriptor instead.
func (*PWT) Descriptor() ([]byte, []int) {
	return file_pwt_v1_pwt_proto_rawDescGZIP(), []int{0}
}

func (x *PWT) GetExpirationTime() *durationpb.Duration {
	if x != nil {
		return x.ExpirationTime
	}
	return nil
}

func (x *PWT) GetIssuedAt() *durationpb.Duration {
	if x != nil {
		return x.IssuedAt
	}
	return nil
}

func (x *PWT) GetNotBefore() *durationpb.Duration {
	if x != nil {
		return x.NotBefore
	}
	return nil
}

func (x *PWT) GetIssuer() string {
	if x != nil {
		return x.Issuer
	}
	return ""
}

func (x *PWT) GetAudience() []string {
	if x != nil {
		return x.Audience
	}
	return nil
}

func (x *PWT) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *PWT) GetJwtId() string {
	if x != nil {
		return x.JwtId
	}
	return ""
}

func (x *PWT) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *PWT) GetClientSecret() string {
	if x != nil {
		return x.ClientSecret
	}
	return ""
}

var File_pwt_v1_pwt_proto protoreflect.FileDescriptor

var file_pwt_v1_pwt_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x77, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x77, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x70, 0x77, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe8, 0x02, 0x0a, 0x03, 0x50,
	0x57, 0x54, 0x12, 0x43, 0x0a, 0x0f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x37, 0x0a, 0x09, 0x69, 0x73, 0x73, 0x75, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x12, 0x39, 0x0a, 0x0a, 0x6e, 0x6f, 0x74, 0x5f, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x0a, 0x6e, 0x6f, 0x74, 0x5f, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x69,
	0x73, 0x73, 0x75, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x69, 0x73, 0x73,
	0x75, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x61, 0x75, 0x64, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x18,
	0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x61, 0x75, 0x64, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6a, 0x77, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6a, 0x77, 0x74, 0x5f, 0x69,
	0x64, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x12,
	0x24, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73,
	0x65, 0x63, 0x72, 0x65, 0x74, 0x42, 0x84, 0x01, 0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x77,
	0x74, 0x2e, 0x76, 0x31, 0x42, 0x08, 0x50, 0x77, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x72, 0x69,
	0x67, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2f, 0x67,
	0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x77, 0x74, 0x2f, 0x76, 0x31, 0x3b, 0x70, 0x77, 0x74,
	0x76, 0x31, 0xf8, 0x01, 0x01, 0xa2, 0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02, 0x06, 0x50, 0x77,
	0x74, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x06, 0x50, 0x77, 0x74, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x12,
	0x50, 0x77, 0x74, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0xea, 0x02, 0x07, 0x50, 0x77, 0x74, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pwt_v1_pwt_proto_rawDescOnce sync.Once
	file_pwt_v1_pwt_proto_rawDescData = file_pwt_v1_pwt_proto_rawDesc
)

func file_pwt_v1_pwt_proto_rawDescGZIP() []byte {
	file_pwt_v1_pwt_proto_rawDescOnce.Do(func() {
		file_pwt_v1_pwt_proto_rawDescData = protoimpl.X.CompressGZIP(file_pwt_v1_pwt_proto_rawDescData)
	})
	return file_pwt_v1_pwt_proto_rawDescData
}

var file_pwt_v1_pwt_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pwt_v1_pwt_proto_goTypes = []any{
	(*PWT)(nil),                 // 0: pwt.v1.PWT
	(*durationpb.Duration)(nil), // 1: google.protobuf.Duration
}
var file_pwt_v1_pwt_proto_depIdxs = []int32{
	1, // 0: pwt.v1.PWT.expiration_time:type_name -> google.protobuf.Duration
	1, // 1: pwt.v1.PWT.issued_at:type_name -> google.protobuf.Duration
	1, // 2: pwt.v1.PWT.not_before:type_name -> google.protobuf.Duration
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_pwt_v1_pwt_proto_init() }
func file_pwt_v1_pwt_proto_init() {
	if File_pwt_v1_pwt_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pwt_v1_pwt_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pwt_v1_pwt_proto_goTypes,
		DependencyIndexes: file_pwt_v1_pwt_proto_depIdxs,
		MessageInfos:      file_pwt_v1_pwt_proto_msgTypes,
	}.Build()
	File_pwt_v1_pwt_proto = out.File
	file_pwt_v1_pwt_proto_rawDesc = nil
	file_pwt_v1_pwt_proto_goTypes = nil
	file_pwt_v1_pwt_proto_depIdxs = nil
}
