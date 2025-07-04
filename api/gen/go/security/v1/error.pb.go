// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: security/v1/error.proto

package securityv1

import (
	_ "github.com/go-kratos/kratos/v2/errors"
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

type SecurityErrorReason int32

const (
	SecurityErrorReason_SECURITY_ERROR_REASON_UNSPECIFIED SecurityErrorReason = 0
	// authentication starts at 1000, and ends at 1999
	SecurityErrorReason_SECURITY_ERROR_REASON_INVALID_AUTHENTICATION     SecurityErrorReason = 1000
	SecurityErrorReason_SECURITY_ERROR_REASON_INVALID_CLAIMS             SecurityErrorReason = 1001
	SecurityErrorReason_SECURITY_ERROR_REASON_INVALID_BEARER_TOKEN       SecurityErrorReason = 1002
	SecurityErrorReason_SECURITY_ERROR_REASON_INVALID_SUBJECT            SecurityErrorReason = 1003
	SecurityErrorReason_SECURITY_ERROR_REASON_INVALID_AUDIENCE           SecurityErrorReason = 1004
	SecurityErrorReason_SECURITY_ERROR_REASON_INVALID_ISSUER             SecurityErrorReason = 1005
	SecurityErrorReason_SECURITY_ERROR_REASON_INVALID_EXPIRATION         SecurityErrorReason = 1006
	SecurityErrorReason_SECURITY_ERROR_REASON_TOKEN_NOT_FOUND            SecurityErrorReason = 1007
	SecurityErrorReason_SECURITY_ERROR_REASON_BEARER_TOKEN_MISSING       SecurityErrorReason = 1010
	SecurityErrorReason_SECURITY_ERROR_REASON_TOKEN_EXPIRED              SecurityErrorReason = 1011
	SecurityErrorReason_SECURITY_ERROR_REASON_UNSUPPORTED_SIGNING_METHOD SecurityErrorReason = 1012
	SecurityErrorReason_SECURITY_ERROR_REASON_MISSING_KEY_FUNC           SecurityErrorReason = 1014
	SecurityErrorReason_SECURITY_ERROR_REASON_SIGN_TOKEN_FAILED          SecurityErrorReason = 1015
	SecurityErrorReason_SECURITY_ERROR_REASON_GET_KEY_FAILED             SecurityErrorReason = 1016
	// authorization starts at 2000, and ends at 2999
	SecurityErrorReason_SECURITY_ERROR_REASON_INVALID_AUTHORIZATION SecurityErrorReason = 2000
	SecurityErrorReason_SECURITY_ERROR_REASON_NO_AT_HASH            SecurityErrorReason = 1050
	SecurityErrorReason_SECURITY_ERROR_REASON_INVALID_AT_HASH       SecurityErrorReason = 1051
	SecurityErrorReason_SECURITY_ERROR_REASON_UNSECURITY_ENTICATED  SecurityErrorReason = 3000
)

// Enum value maps for SecurityErrorReason.
var (
	SecurityErrorReason_name = map[int32]string{
		0:    "SECURITY_ERROR_REASON_UNSPECIFIED",
		1000: "SECURITY_ERROR_REASON_INVALID_AUTHENTICATION",
		1001: "SECURITY_ERROR_REASON_INVALID_CLAIMS",
		1002: "SECURITY_ERROR_REASON_INVALID_BEARER_TOKEN",
		1003: "SECURITY_ERROR_REASON_INVALID_SUBJECT",
		1004: "SECURITY_ERROR_REASON_INVALID_AUDIENCE",
		1005: "SECURITY_ERROR_REASON_INVALID_ISSUER",
		1006: "SECURITY_ERROR_REASON_INVALID_EXPIRATION",
		1007: "SECURITY_ERROR_REASON_TOKEN_NOT_FOUND",
		1010: "SECURITY_ERROR_REASON_BEARER_TOKEN_MISSING",
		1011: "SECURITY_ERROR_REASON_TOKEN_EXPIRED",
		1012: "SECURITY_ERROR_REASON_UNSUPPORTED_SIGNING_METHOD",
		1014: "SECURITY_ERROR_REASON_MISSING_KEY_FUNC",
		1015: "SECURITY_ERROR_REASON_SIGN_TOKEN_FAILED",
		1016: "SECURITY_ERROR_REASON_GET_KEY_FAILED",
		2000: "SECURITY_ERROR_REASON_INVALID_AUTHORIZATION",
		1050: "SECURITY_ERROR_REASON_NO_AT_HASH",
		1051: "SECURITY_ERROR_REASON_INVALID_AT_HASH",
		3000: "SECURITY_ERROR_REASON_UNSECURITY_ENTICATED",
	}
	SecurityErrorReason_value = map[string]int32{
		"SECURITY_ERROR_REASON_UNSPECIFIED":                0,
		"SECURITY_ERROR_REASON_INVALID_AUTHENTICATION":     1000,
		"SECURITY_ERROR_REASON_INVALID_CLAIMS":             1001,
		"SECURITY_ERROR_REASON_INVALID_BEARER_TOKEN":       1002,
		"SECURITY_ERROR_REASON_INVALID_SUBJECT":            1003,
		"SECURITY_ERROR_REASON_INVALID_AUDIENCE":           1004,
		"SECURITY_ERROR_REASON_INVALID_ISSUER":             1005,
		"SECURITY_ERROR_REASON_INVALID_EXPIRATION":         1006,
		"SECURITY_ERROR_REASON_TOKEN_NOT_FOUND":            1007,
		"SECURITY_ERROR_REASON_BEARER_TOKEN_MISSING":       1010,
		"SECURITY_ERROR_REASON_TOKEN_EXPIRED":              1011,
		"SECURITY_ERROR_REASON_UNSUPPORTED_SIGNING_METHOD": 1012,
		"SECURITY_ERROR_REASON_MISSING_KEY_FUNC":           1014,
		"SECURITY_ERROR_REASON_SIGN_TOKEN_FAILED":          1015,
		"SECURITY_ERROR_REASON_GET_KEY_FAILED":             1016,
		"SECURITY_ERROR_REASON_INVALID_AUTHORIZATION":      2000,
		"SECURITY_ERROR_REASON_NO_AT_HASH":                 1050,
		"SECURITY_ERROR_REASON_INVALID_AT_HASH":            1051,
		"SECURITY_ERROR_REASON_UNSECURITY_ENTICATED":       3000,
	}
)

func (x SecurityErrorReason) Enum() *SecurityErrorReason {
	p := new(SecurityErrorReason)
	*p = x
	return p
}

func (x SecurityErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SecurityErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_security_v1_error_proto_enumTypes[0].Descriptor()
}

func (SecurityErrorReason) Type() protoreflect.EnumType {
	return &file_security_v1_error_proto_enumTypes[0]
}

func (x SecurityErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SecurityErrorReason.Descriptor instead.
func (SecurityErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_security_v1_error_proto_rawDescGZIP(), []int{0}
}

var File_security_v1_error_proto protoreflect.FileDescriptor

const file_security_v1_error_proto_rawDesc = "" +
	"\n" +
	"\x17security/v1/error.proto\x12\vsecurity.v1\x1a\x13errors/errors.proto*\xea\a\n" +
	"\x13SecurityErrorReason\x12%\n" +
	"!SECURITY_ERROR_REASON_UNSPECIFIED\x10\x00\x127\n" +
	",SECURITY_ERROR_REASON_INVALID_AUTHENTICATION\x10\xe8\a\x1a\x04\xa8E\x91\x03\x12/\n" +
	"$SECURITY_ERROR_REASON_INVALID_CLAIMS\x10\xe9\a\x1a\x04\xa8E\x91\x03\x125\n" +
	"*SECURITY_ERROR_REASON_INVALID_BEARER_TOKEN\x10\xea\a\x1a\x04\xa8E\x91\x03\x120\n" +
	"%SECURITY_ERROR_REASON_INVALID_SUBJECT\x10\xeb\a\x1a\x04\xa8E\x91\x03\x121\n" +
	"&SECURITY_ERROR_REASON_INVALID_AUDIENCE\x10\xec\a\x1a\x04\xa8E\x91\x03\x12/\n" +
	"$SECURITY_ERROR_REASON_INVALID_ISSUER\x10\xed\a\x1a\x04\xa8E\x91\x03\x123\n" +
	"(SECURITY_ERROR_REASON_INVALID_EXPIRATION\x10\xee\a\x1a\x04\xa8E\x91\x03\x120\n" +
	"%SECURITY_ERROR_REASON_TOKEN_NOT_FOUND\x10\xef\a\x1a\x04\xa8E\x91\x03\x125\n" +
	"*SECURITY_ERROR_REASON_BEARER_TOKEN_MISSING\x10\xf2\a\x1a\x04\xa8E\x91\x03\x12.\n" +
	"#SECURITY_ERROR_REASON_TOKEN_EXPIRED\x10\xf3\a\x1a\x04\xa8E\x91\x03\x12;\n" +
	"0SECURITY_ERROR_REASON_UNSUPPORTED_SIGNING_METHOD\x10\xf4\a\x1a\x04\xa8E\x91\x03\x121\n" +
	"&SECURITY_ERROR_REASON_MISSING_KEY_FUNC\x10\xf6\a\x1a\x04\xa8E\x91\x03\x122\n" +
	"'SECURITY_ERROR_REASON_SIGN_TOKEN_FAILED\x10\xf7\a\x1a\x04\xa8E\x91\x03\x12/\n" +
	"$SECURITY_ERROR_REASON_GET_KEY_FAILED\x10\xf8\a\x1a\x04\xa8E\x91\x03\x126\n" +
	"+SECURITY_ERROR_REASON_INVALID_AUTHORIZATION\x10\xd0\x0f\x1a\x04\xa8E\x93\x03\x12+\n" +
	" SECURITY_ERROR_REASON_NO_AT_HASH\x10\x9a\b\x1a\x04\xa8E\x93\x03\x120\n" +
	"%SECURITY_ERROR_REASON_INVALID_AT_HASH\x10\x9b\b\x1a\x04\xa8E\x93\x03\x125\n" +
	"*SECURITY_ERROR_REASON_UNSECURITY_ENTICATED\x10\xb8\x17\x1a\x04\xa8E\x93\x03\x1a\x04\xa0E\xf4\x03B\xad\x01\n" +
	"\x0fcom.security.v1B\n" +
	"ErrorProtoP\x01Z>github.com/origadmin/runtime/api/gen/go/security/v1;securityv1\xf8\x01\x01\xa2\x02\x03SXX\xaa\x02\vSecurity.V1\xca\x02\vSecurity\\V1\xe2\x02\x17Security\\V1\\GPBMetadata\xea\x02\fSecurity::V1b\x06proto3"

var (
	file_security_v1_error_proto_rawDescOnce sync.Once
	file_security_v1_error_proto_rawDescData []byte
)

func file_security_v1_error_proto_rawDescGZIP() []byte {
	file_security_v1_error_proto_rawDescOnce.Do(func() {
		file_security_v1_error_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_security_v1_error_proto_rawDesc), len(file_security_v1_error_proto_rawDesc)))
	})
	return file_security_v1_error_proto_rawDescData
}

var file_security_v1_error_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_security_v1_error_proto_goTypes = []any{
	(SecurityErrorReason)(0), // 0: security.v1.SecurityErrorReason
}
var file_security_v1_error_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_security_v1_error_proto_init() }
func file_security_v1_error_proto_init() {
	if File_security_v1_error_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_security_v1_error_proto_rawDesc), len(file_security_v1_error_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_security_v1_error_proto_goTypes,
		DependencyIndexes: file_security_v1_error_proto_depIdxs,
		EnumInfos:         file_security_v1_error_proto_enumTypes,
	}.Build()
	File_security_v1_error_proto = out.File
	file_security_v1_error_proto_goTypes = nil
	file_security_v1_error_proto_depIdxs = nil
}
