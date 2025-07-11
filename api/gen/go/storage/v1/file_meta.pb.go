// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: storage/v1/file_meta.proto

package storagev1

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

type FileMeta struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Hash          string                 `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	Size          int64                  `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
	MimeType      string                 `protobuf:"bytes,5,opt,name=mime_type,proto3" json:"mime_type,omitempty"`
	ModTime       int64                  `protobuf:"varint,6,opt,name=mod_time,proto3" json:"mod_time,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FileMeta) Reset() {
	*x = FileMeta{}
	mi := &file_storage_v1_file_meta_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileMeta) ProtoMessage() {}

func (x *FileMeta) ProtoReflect() protoreflect.Message {
	mi := &file_storage_v1_file_meta_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileMeta.ProtoReflect.Descriptor instead.
func (*FileMeta) Descriptor() ([]byte, []int) {
	return file_storage_v1_file_meta_proto_rawDescGZIP(), []int{0}
}

func (x *FileMeta) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *FileMeta) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *FileMeta) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *FileMeta) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *FileMeta) GetMimeType() string {
	if x != nil {
		return x.MimeType
	}
	return ""
}

func (x *FileMeta) GetModTime() int64 {
	if x != nil {
		return x.ModTime
	}
	return 0
}

var File_storage_v1_file_meta_proto protoreflect.FileDescriptor

const file_storage_v1_file_meta_proto_rawDesc = "" +
	"\n" +
	"\x1astorage/v1/file_meta.proto\x12\n" +
	"storage.v1\x1a$gnostic/openapi/v3/annotations.proto\x1a\x17validate/validate.proto\"\x90\x01\n" +
	"\bFileMeta\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x12\n" +
	"\x04hash\x18\x03 \x01(\tR\x04hash\x12\x12\n" +
	"\x04size\x18\x04 \x01(\x03R\x04size\x12\x1c\n" +
	"\tmime_type\x18\x05 \x01(\tR\tmime_type\x12\x1a\n" +
	"\bmod_time\x18\x06 \x01(\x03R\bmod_timeB\xa9\x01\n" +
	"\x0ecom.storage.v1B\rFileMetaProtoP\x01Z<github.com/origadmin/runtime/api/gen/go/storage/v1;storagev1\xf8\x01\x01\xa2\x02\x03SXX\xaa\x02\n" +
	"Storage.V1\xca\x02\n" +
	"Storage\\V1\xe2\x02\x16Storage\\V1\\GPBMetadata\xea\x02\vStorage::V1b\x06proto3"

var (
	file_storage_v1_file_meta_proto_rawDescOnce sync.Once
	file_storage_v1_file_meta_proto_rawDescData []byte
)

func file_storage_v1_file_meta_proto_rawDescGZIP() []byte {
	file_storage_v1_file_meta_proto_rawDescOnce.Do(func() {
		file_storage_v1_file_meta_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_storage_v1_file_meta_proto_rawDesc), len(file_storage_v1_file_meta_proto_rawDesc)))
	})
	return file_storage_v1_file_meta_proto_rawDescData
}

var file_storage_v1_file_meta_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_storage_v1_file_meta_proto_goTypes = []any{
	(*FileMeta)(nil), // 0: storage.v1.FileMeta
}
var file_storage_v1_file_meta_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_storage_v1_file_meta_proto_init() }
func file_storage_v1_file_meta_proto_init() {
	if File_storage_v1_file_meta_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_storage_v1_file_meta_proto_rawDesc), len(file_storage_v1_file_meta_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_storage_v1_file_meta_proto_goTypes,
		DependencyIndexes: file_storage_v1_file_meta_proto_depIdxs,
		MessageInfos:      file_storage_v1_file_meta_proto_msgTypes,
	}.Build()
	File_storage_v1_file_meta_proto = out.File
	file_storage_v1_file_meta_proto_goTypes = nil
	file_storage_v1_file_meta_proto_depIdxs = nil
}
