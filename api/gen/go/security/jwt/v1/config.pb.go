// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: security/jwt/v1/config.proto

package jwtv1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "github.com/google/gnostic/openapiv3"
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

// Config contains configuration parameters for creating and validating a JWT.
type Config struct {
	state                protoimpl.MessageState `protogen:"open.v1"`
	SigningMethod        string                 `protobuf:"bytes,1,opt,name=signing_method,proto3" json:"signing_method,omitempty"`
	Key                  string                 `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Key2                 string                 `protobuf:"bytes,3,opt,name=key2,proto3" json:"key2,omitempty"`
	AccessTokenLifetime  int64                  `protobuf:"varint,5,opt,name=access_token_lifetime,proto3" json:"access_token_lifetime,omitempty"`
	RefreshTokenLifetime int64                  `protobuf:"varint,6,opt,name=refresh_token_lifetime,proto3" json:"refresh_token_lifetime,omitempty"`
	Issuer               string                 `protobuf:"bytes,7,opt,name=issuer,proto3" json:"issuer,omitempty"`
	Audience             []string               `protobuf:"bytes,8,rep,name=audience,proto3" json:"audience,omitempty"` // Audience
	TokenType            string                 `protobuf:"bytes,9,opt,name=token_type,proto3" json:"token_type,omitempty"`
	unknownFields        protoimpl.UnknownFields
	sizeCache            protoimpl.SizeCache
}

func (x *Config) Reset() {
	*x = Config{}
	mi := &file_security_jwt_v1_config_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_security_jwt_v1_config_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_security_jwt_v1_config_proto_rawDescGZIP(), []int{0}
}

func (x *Config) GetSigningMethod() string {
	if x != nil {
		return x.SigningMethod
	}
	return ""
}

func (x *Config) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Config) GetKey2() string {
	if x != nil {
		return x.Key2
	}
	return ""
}

func (x *Config) GetAccessTokenLifetime() int64 {
	if x != nil {
		return x.AccessTokenLifetime
	}
	return 0
}

func (x *Config) GetRefreshTokenLifetime() int64 {
	if x != nil {
		return x.RefreshTokenLifetime
	}
	return 0
}

func (x *Config) GetIssuer() string {
	if x != nil {
		return x.Issuer
	}
	return ""
}

func (x *Config) GetAudience() []string {
	if x != nil {
		return x.Audience
	}
	return nil
}

func (x *Config) GetTokenType() string {
	if x != nil {
		return x.TokenType
	}
	return ""
}

var File_security_jwt_v1_config_proto protoreflect.FileDescriptor

const file_security_jwt_v1_config_proto_rawDesc = "" +
	"\n" +
	"\x1csecurity/jwt/v1/config.proto\x12\x0fsecurity.jwt.v1\x1a$gnostic/openapi/v3/annotations.proto\x1a\x17validate/validate.proto\"\xda\x05\n" +
	"\x06Config\x12\x80\x01\n" +
	"\x0esigning_method\x18\x01 \x01(\tBX\xfaB\x14r\x12\x10\x01\x18\x80\b2\v^[A-Z0-9]+$\xbaG>\x92\x02;The signing method used for the token (e.g., HS256, RS256).R\x0esigning_method\x12E\n" +
	"\x03key\x18\x02 \x01(\tB3\xfaB\ar\x05\x10\x01\x18\x80\b\xbaG&\x92\x02#The key used for signing the token.R\x03key\x12G\n" +
	"\x04key2\x18\x03 \x01(\tB3\xbaG0\x92\x02-The secondary key used for signing the token.R\x04key2\x12b\n" +
	"\x15access_token_lifetime\x18\x05 \x01(\x03B,\xfaB\t\"\a\x18\x80\xe7\x84\x0f(\x01\xbaG\x1d\x92\x02\x1aThe lifetime of the token.R\x15access_token_lifetime\x12l\n" +
	"\x16refresh_token_lifetime\x18\x06 \x01(\x03B4\xfaB\t\"\a\x18\x80\xe7\x84\x0f(\x01\xbaG%\x92\x02\"The lifetime of the refresh token.R\x16refresh_token_lifetime\x126\n" +
	"\x06issuer\x18\a \x01(\tB\x1e\xbaG\x1b\x92\x02\x18The issuer of the token.R\x06issuer\x12\\\n" +
	"\baudience\x18\b \x03(\tB@\xfaB\n" +
	"\x92\x01\a\b\x01\x10\x80\b\x18\x01\xbaG0\x92\x02-The audience for which the token is intended.R\baudience\x12U\n" +
	"\n" +
	"token_type\x18\t \x01(\tB5\xfaB\ar\x05\x10\x01\x18\x80\b\xbaG(\x92\x02%The type of the token (e.g., Bearer).R\n" +
	"token_typeB\xc2\x01\n" +
	"\x13com.security.jwt.v1B\vConfigProtoP\x01Z=github.com/origadmin/runtime/api/gen/go/security/jwt/v1;jwtv1\xf8\x01\x01\xa2\x02\x03SJX\xaa\x02\x0fSecurity.Jwt.V1\xca\x02\x0fSecurity\\Jwt\\V1\xe2\x02\x1bSecurity\\Jwt\\V1\\GPBMetadata\xea\x02\x11Security::Jwt::V1b\x06proto3"

var (
	file_security_jwt_v1_config_proto_rawDescOnce sync.Once
	file_security_jwt_v1_config_proto_rawDescData []byte
)

func file_security_jwt_v1_config_proto_rawDescGZIP() []byte {
	file_security_jwt_v1_config_proto_rawDescOnce.Do(func() {
		file_security_jwt_v1_config_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_security_jwt_v1_config_proto_rawDesc), len(file_security_jwt_v1_config_proto_rawDesc)))
	})
	return file_security_jwt_v1_config_proto_rawDescData
}

var file_security_jwt_v1_config_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_security_jwt_v1_config_proto_goTypes = []any{
	(*Config)(nil), // 0: security.jwt.v1.Config
}
var file_security_jwt_v1_config_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_security_jwt_v1_config_proto_init() }
func file_security_jwt_v1_config_proto_init() {
	if File_security_jwt_v1_config_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_security_jwt_v1_config_proto_rawDesc), len(file_security_jwt_v1_config_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_security_jwt_v1_config_proto_goTypes,
		DependencyIndexes: file_security_jwt_v1_config_proto_depIdxs,
		MessageInfos:      file_security_jwt_v1_config_proto_msgTypes,
	}.Build()
	File_security_jwt_v1_config_proto = out.File
	file_security_jwt_v1_config_proto_goTypes = nil
	file_security_jwt_v1_config_proto_depIdxs = nil
}
