// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: middleware/v1/validator/validator.proto

package validatorv1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Validator struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Enabled       bool                   `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Version       int32                  `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	FailFast      bool                   `protobuf:"varint,3,opt,name=fail_fast,proto3" json:"fail_fast,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Validator) Reset() {
	*x = Validator{}
	mi := &file_middleware_v1_validator_validator_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Validator) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Validator) ProtoMessage() {}

func (x *Validator) ProtoReflect() protoreflect.Message {
	mi := &file_middleware_v1_validator_validator_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Validator.ProtoReflect.Descriptor instead.
func (*Validator) Descriptor() ([]byte, []int) {
	return file_middleware_v1_validator_validator_proto_rawDescGZIP(), []int{0}
}

func (x *Validator) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *Validator) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Validator) GetFailFast() bool {
	if x != nil {
		return x.FailFast
	}
	return false
}

var File_middleware_v1_validator_validator_proto protoreflect.FileDescriptor

const file_middleware_v1_validator_validator_proto_rawDesc = "" +
	"\n" +
	"'middleware/v1/validator/validator.proto\x12\x17middleware.v1.validator\x1a\x17validate/validate.proto\"h\n" +
	"\tValidator\x12\x18\n" +
	"\aenabled\x18\x01 \x01(\bR\aenabled\x12#\n" +
	"\aversion\x18\x02 \x01(\x05B\t\xfaB\x06\x1a\x04\x10\x03 \x00R\aversion\x12\x1c\n" +
	"\tfail_fast\x18\x03 \x01(\bR\tfail_fastB\xfb\x01\n" +
	"\x1bcom.middleware.v1.validatorB\x0eValidatorProtoP\x01ZKgithub.com/origadmin/runtime/api/gen/go/middleware/v1/validator;validatorv1\xf8\x01\x01\xa2\x02\x03MVV\xaa\x02\x17Middleware.V1.Validator\xca\x02\x17Middleware\\V1\\Validator\xe2\x02#Middleware\\V1\\Validator\\GPBMetadata\xea\x02\x19Middleware::V1::Validatorb\x06proto3"

var (
	file_middleware_v1_validator_validator_proto_rawDescOnce sync.Once
	file_middleware_v1_validator_validator_proto_rawDescData []byte
)

func file_middleware_v1_validator_validator_proto_rawDescGZIP() []byte {
	file_middleware_v1_validator_validator_proto_rawDescOnce.Do(func() {
		file_middleware_v1_validator_validator_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_middleware_v1_validator_validator_proto_rawDesc), len(file_middleware_v1_validator_validator_proto_rawDesc)))
	})
	return file_middleware_v1_validator_validator_proto_rawDescData
}

var file_middleware_v1_validator_validator_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_middleware_v1_validator_validator_proto_goTypes = []any{
	(*Validator)(nil), // 0: middleware.v1.validator.Validator
}
var file_middleware_v1_validator_validator_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_middleware_v1_validator_validator_proto_init() }
func file_middleware_v1_validator_validator_proto_init() {
	if File_middleware_v1_validator_validator_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_middleware_v1_validator_validator_proto_rawDesc), len(file_middleware_v1_validator_validator_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_middleware_v1_validator_validator_proto_goTypes,
		DependencyIndexes: file_middleware_v1_validator_validator_proto_depIdxs,
		MessageInfos:      file_middleware_v1_validator_validator_proto_msgTypes,
	}.Build()
	File_middleware_v1_validator_validator_proto = out.File
	file_middleware_v1_validator_validator_proto_goTypes = nil
	file_middleware_v1_validator_validator_proto_depIdxs = nil
}
