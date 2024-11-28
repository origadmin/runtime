// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: config/v1/tlsconfig.proto

package configv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// TLSConfig
type TLSConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	File *TLSConfig_File `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	Pem  *TLSConfig_PEM  `protobuf:"bytes,2,opt,name=pem,proto3" json:"pem,omitempty"`
}

func (x *TLSConfig) Reset() {
	*x = TLSConfig{}
	mi := &file_config_v1_tlsconfig_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TLSConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TLSConfig) ProtoMessage() {}

func (x *TLSConfig) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_tlsconfig_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TLSConfig.ProtoReflect.Descriptor instead.
func (*TLSConfig) Descriptor() ([]byte, []int) {
	return file_config_v1_tlsconfig_proto_rawDescGZIP(), []int{0}
}

func (x *TLSConfig) GetFile() *TLSConfig_File {
	if x != nil {
		return x.File
	}
	return nil
}

func (x *TLSConfig) GetPem() *TLSConfig_PEM {
	if x != nil {
		return x.Pem
	}
	return nil
}

type TLSConfig_File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cert string `protobuf:"bytes,1,opt,name=cert,proto3" json:"cert,omitempty"`
	Key  string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Ca   string `protobuf:"bytes,3,opt,name=ca,proto3" json:"ca,omitempty"`
}

func (x *TLSConfig_File) Reset() {
	*x = TLSConfig_File{}
	mi := &file_config_v1_tlsconfig_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TLSConfig_File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TLSConfig_File) ProtoMessage() {}

func (x *TLSConfig_File) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_tlsconfig_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TLSConfig_File.ProtoReflect.Descriptor instead.
func (*TLSConfig_File) Descriptor() ([]byte, []int) {
	return file_config_v1_tlsconfig_proto_rawDescGZIP(), []int{0, 0}
}

func (x *TLSConfig_File) GetCert() string {
	if x != nil {
		return x.Cert
	}
	return ""
}

func (x *TLSConfig_File) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *TLSConfig_File) GetCa() string {
	if x != nil {
		return x.Ca
	}
	return ""
}

type TLSConfig_PEM struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cert []byte `protobuf:"bytes,1,opt,name=cert,proto3" json:"cert,omitempty"`
	Key  []byte `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Ca   []byte `protobuf:"bytes,3,opt,name=ca,proto3" json:"ca,omitempty"`
}

func (x *TLSConfig_PEM) Reset() {
	*x = TLSConfig_PEM{}
	mi := &file_config_v1_tlsconfig_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TLSConfig_PEM) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TLSConfig_PEM) ProtoMessage() {}

func (x *TLSConfig_PEM) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_tlsconfig_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TLSConfig_PEM.ProtoReflect.Descriptor instead.
func (*TLSConfig_PEM) Descriptor() ([]byte, []int) {
	return file_config_v1_tlsconfig_proto_rawDescGZIP(), []int{0, 1}
}

func (x *TLSConfig_PEM) GetCert() []byte {
	if x != nil {
		return x.Cert
	}
	return nil
}

func (x *TLSConfig_PEM) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *TLSConfig_PEM) GetCa() []byte {
	if x != nil {
		return x.Ca
	}
	return nil
}

var File_config_v1_tlsconfig_proto protoreflect.FileDescriptor

var file_config_v1_tlsconfig_proto_rawDesc = []byte{
	0x0a, 0x19, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x6c, 0x73, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x22, 0xe1, 0x01, 0x0a, 0x09, 0x54, 0x4c, 0x53, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x12, 0x2d, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x54,
	0x4c, 0x53, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x04, 0x66,
	0x69, 0x6c, 0x65, 0x12, 0x2a, 0x0a, 0x03, 0x70, 0x65, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x18, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x4c, 0x53,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x50, 0x45, 0x4d, 0x52, 0x03, 0x70, 0x65, 0x6d, 0x1a,
	0x3c, 0x0a, 0x04, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x65, 0x72, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x65, 0x72, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x0e, 0x0a,
	0x02, 0x63, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x63, 0x61, 0x1a, 0x3b, 0x0a,
	0x03, 0x50, 0x45, 0x4d, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x65, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x04, 0x63, 0x65, 0x72, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x63, 0x61,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x63, 0x61, 0x42, 0x9f, 0x01, 0x0a, 0x0d, 0x63,
	0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x42, 0x0e, 0x54, 0x6c,
	0x73, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x36,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x72, 0x69, 0x67, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2f, 0x67, 0x65, 0x6e,
	0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x76, 0x31, 0xf8, 0x01, 0x01, 0xa2, 0x02, 0x03, 0x43, 0x58, 0x58, 0xaa,
	0x02, 0x09, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x09, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x15, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea,
	0x02, 0x0a, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_config_v1_tlsconfig_proto_rawDescOnce sync.Once
	file_config_v1_tlsconfig_proto_rawDescData = file_config_v1_tlsconfig_proto_rawDesc
)

func file_config_v1_tlsconfig_proto_rawDescGZIP() []byte {
	file_config_v1_tlsconfig_proto_rawDescOnce.Do(func() {
		file_config_v1_tlsconfig_proto_rawDescData = protoimpl.X.CompressGZIP(file_config_v1_tlsconfig_proto_rawDescData)
	})
	return file_config_v1_tlsconfig_proto_rawDescData
}

var file_config_v1_tlsconfig_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_config_v1_tlsconfig_proto_goTypes = []any{
	(*TLSConfig)(nil),      // 0: config.v1.TLSConfig
	(*TLSConfig_File)(nil), // 1: config.v1.TLSConfig.File
	(*TLSConfig_PEM)(nil),  // 2: config.v1.TLSConfig.PEM
}
var file_config_v1_tlsconfig_proto_depIdxs = []int32{
	1, // 0: config.v1.TLSConfig.file:type_name -> config.v1.TLSConfig.File
	2, // 1: config.v1.TLSConfig.pem:type_name -> config.v1.TLSConfig.PEM
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_config_v1_tlsconfig_proto_init() }
func file_config_v1_tlsconfig_proto_init() {
	if File_config_v1_tlsconfig_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_config_v1_tlsconfig_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_v1_tlsconfig_proto_goTypes,
		DependencyIndexes: file_config_v1_tlsconfig_proto_depIdxs,
		MessageInfos:      file_config_v1_tlsconfig_proto_msgTypes,
	}.Build()
	File_config_v1_tlsconfig_proto = out.File
	file_config_v1_tlsconfig_proto_rawDesc = nil
	file_config_v1_tlsconfig_proto_goTypes = nil
	file_config_v1_tlsconfig_proto_depIdxs = nil
}
