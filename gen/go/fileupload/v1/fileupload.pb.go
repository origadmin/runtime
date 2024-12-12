// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: fileupload/v1/fileupload.proto

package fileuploadv1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "github.com/google/gnostic/openapiv3"
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

// FileHeader defines the structure of a file header.
type FileHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filename      string            `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	Size          uint32            `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	ModTimeString string            `protobuf:"bytes,4,opt,name=mod_time_string,proto3" json:"mod_time_string,omitempty"`
	ModTime       uint32            `protobuf:"varint,5,opt,name=mod_time,proto3" json:"mod_time,omitempty"`
	ContentType   string            `protobuf:"bytes,2,opt,name=content_type,proto3" json:"content_type,omitempty"`
	Header        map[string]string `protobuf:"bytes,6,rep,name=header,proto3" json:"header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	IsDir         bool              `protobuf:"varint,7,opt,name=is_dir,proto3" json:"is_dir,omitempty"`
}

func (x *FileHeader) Reset() {
	*x = FileHeader{}
	mi := &file_fileupload_v1_fileupload_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileHeader) ProtoMessage() {}

func (x *FileHeader) ProtoReflect() protoreflect.Message {
	mi := &file_fileupload_v1_fileupload_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileHeader.ProtoReflect.Descriptor instead.
func (*FileHeader) Descriptor() ([]byte, []int) {
	return file_fileupload_v1_fileupload_proto_rawDescGZIP(), []int{0}
}

func (x *FileHeader) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *FileHeader) GetSize() uint32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *FileHeader) GetModTimeString() string {
	if x != nil {
		return x.ModTimeString
	}
	return ""
}

func (x *FileHeader) GetModTime() uint32 {
	if x != nil {
		return x.ModTime
	}
	return 0
}

func (x *FileHeader) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}

func (x *FileHeader) GetHeader() map[string]string {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *FileHeader) GetIsDir() bool {
	if x != nil {
		return x.IsDir
	}
	return false
}

// UploadRequest file block information
type UploadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsHeader bool   `protobuf:"varint,1,opt,name=is_header,proto3" json:"is_header,omitempty"`
	Data     []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *UploadRequest) Reset() {
	*x = UploadRequest{}
	mi := &file_fileupload_v1_fileupload_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadRequest) ProtoMessage() {}

func (x *UploadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_fileupload_v1_fileupload_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadRequest.ProtoReflect.Descriptor instead.
func (*UploadRequest) Descriptor() ([]byte, []int) {
	return file_fileupload_v1_fileupload_proto_rawDescGZIP(), []int{1}
}

func (x *UploadRequest) GetIsHeader() bool {
	if x != nil {
		return x.IsHeader
	}
	return false
}

func (x *UploadRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

// UploadResponse defines the structure of a file response.
type UploadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success    bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Hash       string `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	Path       string `protobuf:"bytes,3,opt,name=path,proto3" json:"path,omitempty"`
	Size       uint32 `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
	FailReason string `protobuf:"bytes,5,opt,name=fail_reason,proto3" json:"fail_reason,omitempty"`
}

func (x *UploadResponse) Reset() {
	*x = UploadResponse{}
	mi := &file_fileupload_v1_fileupload_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadResponse) ProtoMessage() {}

func (x *UploadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_fileupload_v1_fileupload_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadResponse.ProtoReflect.Descriptor instead.
func (*UploadResponse) Descriptor() ([]byte, []int) {
	return file_fileupload_v1_fileupload_proto_rawDescGZIP(), []int{2}
}

func (x *UploadResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *UploadResponse) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *UploadResponse) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *UploadResponse) GetSize() uint32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *UploadResponse) GetFailReason() string {
	if x != nil {
		return x.FailReason
	}
	return ""
}

var File_fileupload_v1_fileupload_proto protoreflect.FileDescriptor

var file_fileupload_v1_fileupload_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x66, 0x69, 0x6c, 0x65, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x2f, 0x76, 0x31, 0x2f,
	0x66, 0x69, 0x6c, 0x65, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0d, 0x66, 0x69, 0x6c, 0x65, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x76, 0x31, 0x1a,
	0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x24, 0x67, 0x6e, 0x6f, 0x73, 0x74, 0x69,
	0x63, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x33, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd5,
	0x03, 0x0a, 0x0a, 0x46, 0x69, 0x6c, 0x65, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x31, 0x0a,
	0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x15, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0xba, 0x47, 0x0b, 0x92, 0x02, 0x08, 0x66, 0x69,
	0x65, 0x20, 0x6e, 0x61, 0x6d, 0x65, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x23, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x0f,
	0xba, 0x47, 0x0c, 0x92, 0x02, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x20, 0x73, 0x69, 0x7a, 0x65, 0x52,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x44, 0x0a, 0x0f, 0x6d, 0x6f, 0x64, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x1a,
	0xba, 0x47, 0x17, 0x92, 0x02, 0x14, 0x66, 0x69, 0x6c, 0x65, 0x20, 0x6d, 0x6f, 0x64, 0x20, 0x74,
	0x69, 0x6d, 0x65, 0x20, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x0f, 0x6d, 0x6f, 0x64, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x34, 0x0a, 0x08, 0x6d,
	0x6f, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x18, 0xba,
	0x47, 0x15, 0x92, 0x02, 0x12, 0x66, 0x69, 0x6c, 0x65, 0x20, 0x6d, 0x6f, 0x64, 0x20, 0x74, 0x69,
	0x6d, 0x65, 0x20, 0x75, 0x6e, 0x69, 0x78, 0x52, 0x08, 0x6d, 0x6f, 0x64, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x12, 0x3b, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x17, 0xba, 0x47, 0x14, 0x92, 0x02, 0x11, 0x66,
	0x69, 0x6c, 0x65, 0x20, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x20, 0x74, 0x79, 0x70, 0x65,
	0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x12, 0x50,
	0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25,
	0x2e, 0x66, 0x69, 0x6c, 0x65, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x42, 0x11, 0xba, 0x47, 0x0e, 0x92, 0x02, 0x0b, 0x66, 0x69, 0x6c,
	0x65, 0x20, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x12, 0x29, 0x0a, 0x06, 0x69, 0x73, 0x5f, 0x64, 0x69, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08,
	0x42, 0x11, 0xba, 0x47, 0x0e, 0x92, 0x02, 0x0b, 0x66, 0x69, 0x6c, 0x65, 0x20, 0x69, 0x73, 0x20,
	0x64, 0x69, 0x72, 0x52, 0x06, 0x69, 0x73, 0x5f, 0x64, 0x69, 0x72, 0x1a, 0x39, 0x0a, 0x0b, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x65, 0x0a, 0x0d, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x68, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x42, 0x11, 0xba, 0x47, 0x0e, 0x92,
	0x02, 0x0b, 0x66, 0x69, 0x6c, 0x65, 0x20, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x09, 0x69,
	0x73, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x23, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x0f, 0xba, 0x47, 0x0c, 0x92, 0x02, 0x09, 0x66, 0x69,
	0x6c, 0x65, 0x20, 0x64, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xfd, 0x01,
	0x0a, 0x0e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x33, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x42, 0x19, 0xba, 0x47, 0x16, 0x92, 0x02, 0x13, 0x66, 0x69, 0x6c, 0x65, 0x20, 0x75, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x20, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x07, 0x73, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x23, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x0f, 0xba, 0x47, 0x0c, 0x92, 0x02, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x20,
	0x68, 0x61, 0x73, 0x68, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x23, 0x0a, 0x04, 0x70, 0x61,
	0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0f, 0xba, 0x47, 0x0c, 0x92, 0x02, 0x09,
	0x66, 0x69, 0x6c, 0x65, 0x20, 0x70, 0x61, 0x74, 0x68, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12,
	0x23, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x0f, 0xba,
	0x47, 0x0c, 0x92, 0x02, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x20, 0x73, 0x69, 0x7a, 0x65, 0x52, 0x04,
	0x73, 0x69, 0x7a, 0x65, 0x12, 0x47, 0x0a, 0x0b, 0x66, 0x61, 0x69, 0x6c, 0x5f, 0x72, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x25, 0xba, 0x47, 0x22, 0x92, 0x02,
	0x1f, 0x66, 0x69, 0x6c, 0x65, 0x20, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x20, 0x72, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x20, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x20, 0x69, 0x6e, 0x66, 0x6f,
	0x52, 0x0b, 0x66, 0x61, 0x69, 0x6c, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x32, 0x5c, 0x0a,
	0x11, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x47, 0x0a, 0x06, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x1c, 0x2e, 0x66,
	0x69, 0x6c, 0x65, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x66, 0x69, 0x6c,
	0x65, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x42, 0xbc, 0x01, 0x0a, 0x11,
	0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x76,
	0x31, 0x42, 0x0f, 0x46, 0x69, 0x6c, 0x65, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6f, 0x72, 0x69, 0x67, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x72, 0x75, 0x6e, 0x74, 0x69,
	0x6d, 0x65, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x75, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x2f, 0x76, 0x31, 0x3b, 0x66, 0x69, 0x6c, 0x65, 0x75, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x76, 0x31, 0xf8, 0x01, 0x01, 0xa2, 0x02, 0x03, 0x46, 0x58, 0x58, 0xaa, 0x02, 0x0d,
	0x46, 0x69, 0x6c, 0x65, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0d,
	0x46, 0x69, 0x6c, 0x65, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x19,
	0x46, 0x69, 0x6c, 0x65, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0e, 0x46, 0x69, 0x6c, 0x65,
	0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_fileupload_v1_fileupload_proto_rawDescOnce sync.Once
	file_fileupload_v1_fileupload_proto_rawDescData = file_fileupload_v1_fileupload_proto_rawDesc
)

func file_fileupload_v1_fileupload_proto_rawDescGZIP() []byte {
	file_fileupload_v1_fileupload_proto_rawDescOnce.Do(func() {
		file_fileupload_v1_fileupload_proto_rawDescData = protoimpl.X.CompressGZIP(file_fileupload_v1_fileupload_proto_rawDescData)
	})
	return file_fileupload_v1_fileupload_proto_rawDescData
}

var file_fileupload_v1_fileupload_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_fileupload_v1_fileupload_proto_goTypes = []any{
	(*FileHeader)(nil),     // 0: fileupload.v1.FileHeader
	(*UploadRequest)(nil),  // 1: fileupload.v1.UploadRequest
	(*UploadResponse)(nil), // 2: fileupload.v1.UploadResponse
	nil,                    // 3: fileupload.v1.FileHeader.HeaderEntry
}
var file_fileupload_v1_fileupload_proto_depIdxs = []int32{
	3, // 0: fileupload.v1.FileHeader.header:type_name -> fileupload.v1.FileHeader.HeaderEntry
	1, // 1: fileupload.v1.FileUploadService.Upload:input_type -> fileupload.v1.UploadRequest
	2, // 2: fileupload.v1.FileUploadService.Upload:output_type -> fileupload.v1.UploadResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_fileupload_v1_fileupload_proto_init() }
func file_fileupload_v1_fileupload_proto_init() {
	if File_fileupload_v1_fileupload_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_fileupload_v1_fileupload_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_fileupload_v1_fileupload_proto_goTypes,
		DependencyIndexes: file_fileupload_v1_fileupload_proto_depIdxs,
		MessageInfos:      file_fileupload_v1_fileupload_proto_msgTypes,
	}.Build()
	File_fileupload_v1_fileupload_proto = out.File
	file_fileupload_v1_fileupload_proto_rawDesc = nil
	file_fileupload_v1_fileupload_proto_goTypes = nil
	file_fileupload_v1_fileupload_proto_depIdxs = nil
}