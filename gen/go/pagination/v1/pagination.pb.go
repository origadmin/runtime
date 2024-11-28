// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: pagination/v1/pagination.proto

package paginationv1

import (
	_ "github.com/google/gnostic/openapiv3"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// PageRequest common request
type PageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// current page number
	Current *int32 `protobuf:"varint,1,opt,name=current,proto3,oneof" json:"current,omitempty"`
	// The number of lines per page
	PageSize *int32 `protobuf:"varint,2,opt,name=page_size,proto3,oneof" json:"page_size,omitempty"`
	// The page_token is the query parameter for set the page token.
	PageToken string `protobuf:"bytes,3,opt,name=page_token,proto3" json:"page_token,omitempty"`
	// The only_count is the query parameter for set only to query the total number
	OnlyCount bool `protobuf:"varint,4,opt,name=only_count,proto3" json:"only_count,omitempty"`
	// The no_paging is used to disable pagination.
	NoPaging *bool `protobuf:"varint,5,opt,name=no_paging,proto3,oneof" json:"no_paging,omitempty"`
	// sort condition
	OrderBy []string `protobuf:"bytes,6,rep,name=order_by,proto3" json:"order_by,omitempty"`
	// Field mask
	FieldMask *fieldmaskpb.FieldMask `protobuf:"bytes,7,opt,name=field_mask,proto3" json:"field_mask,omitempty"`
}

func (x *PageRequest) Reset() {
	*x = PageRequest{}
	mi := &file_pagination_v1_pagination_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PageRequest) ProtoMessage() {}

func (x *PageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pagination_v1_pagination_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PageRequest.ProtoReflect.Descriptor instead.
func (*PageRequest) Descriptor() ([]byte, []int) {
	return file_pagination_v1_pagination_proto_rawDescGZIP(), []int{0}
}

func (x *PageRequest) GetCurrent() int32 {
	if x != nil && x.Current != nil {
		return *x.Current
	}
	return 0
}

func (x *PageRequest) GetPageSize() int32 {
	if x != nil && x.PageSize != nil {
		return *x.PageSize
	}
	return 0
}

func (x *PageRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *PageRequest) GetOnlyCount() bool {
	if x != nil {
		return x.OnlyCount
	}
	return false
}

func (x *PageRequest) GetNoPaging() bool {
	if x != nil && x.NoPaging != nil {
		return *x.NoPaging
	}
	return false
}

func (x *PageRequest) GetOrderBy() []string {
	if x != nil {
		return x.OrderBy
	}
	return nil
}

func (x *PageRequest) GetFieldMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.FieldMask
	}
	return nil
}

// PageResponse general result
type PageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	// The total number of items in the list.
	Total int32 `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
	// The paging data
	Data *anypb.Any `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	// The current page number.
	Current *int32 `protobuf:"varint,4,opt,name=current,proto3,oneof" json:"current,omitempty"`
	// The maximum number of items to return.
	PageSize *int32 `protobuf:"varint,5,opt,name=page_size,proto3,oneof" json:"page_size,omitempty"`
	// Token to retrieve the next page of results, or empty if there are no
	// more results in the list.
	NextPageToken string `protobuf:"bytes,6,opt,name=next_page_token,proto3" json:"next_page_token,omitempty"`
	// Additional information about this response.
	// content to be added without destroying the current data format
	Extra *anypb.Any `protobuf:"bytes,7,opt,name=extra,proto3" json:"extra,omitempty"`
}

func (x *PageResponse) Reset() {
	*x = PageResponse{}
	mi := &file_pagination_v1_pagination_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PageResponse) ProtoMessage() {}

func (x *PageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pagination_v1_pagination_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PageResponse.ProtoReflect.Descriptor instead.
func (*PageResponse) Descriptor() ([]byte, []int) {
	return file_pagination_v1_pagination_proto_rawDescGZIP(), []int{1}
}

func (x *PageResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *PageResponse) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *PageResponse) GetData() *anypb.Any {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *PageResponse) GetCurrent() int32 {
	if x != nil && x.Current != nil {
		return *x.Current
	}
	return 0
}

func (x *PageResponse) GetPageSize() int32 {
	if x != nil && x.PageSize != nil {
		return *x.PageSize
	}
	return 0
}

func (x *PageResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

func (x *PageResponse) GetExtra() *anypb.Any {
	if x != nil {
		return x.Extra
	}
	return nil
}

var File_pagination_v1_pagination_proto protoreflect.FileDescriptor

var file_pagination_v1_pagination_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f,
	0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0d, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x1a,
	0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x24, 0x67, 0x6e,
	0x6f, 0x73, 0x74, 0x69, 0x63, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x33,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xfc, 0x04, 0x0a, 0x0b, 0x50, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x45, 0x0a, 0x07, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x42, 0x26, 0xba, 0x47, 0x23, 0x8a, 0x02, 0x09, 0x09, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0xf0, 0x3f, 0x92, 0x02, 0x14, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x20, 0x70,
	0x61, 0x67, 0x65, 0x20, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x20, 0x48, 0x00, 0x52, 0x07, 0x63,
	0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x88, 0x01, 0x01, 0x12, 0x51, 0x0a, 0x09, 0x70, 0x61, 0x67,
	0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x42, 0x2e, 0xba, 0x47,
	0x2b, 0x8a, 0x02, 0x09, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e, 0x40, 0x92, 0x02, 0x1c,
	0x54, 0x68, 0x65, 0x20, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x20, 0x6f, 0x66, 0x20, 0x6c, 0x69,
	0x6e, 0x65, 0x73, 0x20, 0x70, 0x65, 0x72, 0x20, 0x70, 0x61, 0x67, 0x65, 0x48, 0x01, 0x52, 0x09,
	0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x88, 0x01, 0x01, 0x12, 0x32, 0x0a, 0x0a,
	0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x12, 0xba, 0x47, 0x0f, 0x92, 0x02, 0x0c, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x20, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x36, 0x0a, 0x0a, 0x6f, 0x6e, 0x6c, 0x79, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x08, 0x42, 0x16, 0xba, 0x47, 0x13, 0x92, 0x02, 0x10, 0x71, 0x75, 0x65, 0x72,
	0x79, 0x20, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x20, 0x6f, 0x6e, 0x6c, 0x79, 0x52, 0x0a, 0x6f, 0x6e,
	0x6c, 0x79, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x3b, 0x0a, 0x09, 0x6e, 0x6f, 0x5f, 0x70,
	0x61, 0x67, 0x69, 0x6e, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x42, 0x18, 0xba, 0x47, 0x15,
	0x92, 0x02, 0x12, 0x77, 0x68, 0x65, 0x74, 0x68, 0x65, 0x72, 0x20, 0x6e, 0x6f, 0x74, 0x20, 0x70,
	0x61, 0x67, 0x69, 0x6e, 0x67, 0x48, 0x02, 0x52, 0x09, 0x6e, 0x6f, 0x5f, 0x70, 0x61, 0x67, 0x69,
	0x6e, 0x67, 0x88, 0x01, 0x01, 0x12, 0x86, 0x01, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f,
	0x62, 0x79, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x42, 0x6a, 0xba, 0x47, 0x67, 0x3a, 0x13, 0x12,
	0x11, 0x7b, 0x22, 0x76, 0x61, 0x6c, 0x31, 0x22, 0x2c, 0x20, 0x22, 0x2d, 0x76, 0x61, 0x6c, 0x32,
	0x22, 0x7d, 0x92, 0x02, 0x4f, 0x73, 0x6f, 0x72, 0x74, 0x20, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x2c, 0x20, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x20, 0x6e, 0x61, 0x6d, 0x65, 0x20,
	0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x20, 0x62, 0x79, 0x20, 0x27, 0x61, 0x73, 0x63,
	0x27, 0x20, 0x28, 0x61, 0x73, 0x63, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x29, 0x20, 0x6f, 0x72,
	0x20, 0x27, 0x64, 0x65, 0x73, 0x63, 0x27, 0x20, 0x28, 0x64, 0x65, 0x73, 0x63, 0x65, 0x6e, 0x64,
	0x69, 0x6e, 0x67, 0x29, 0x52, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x62, 0x79, 0x12, 0x79,
	0x0a, 0x0a, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4d, 0x61, 0x73, 0x6b, 0x42, 0x3d,
	0xba, 0x47, 0x3a, 0x3a, 0x0d, 0x12, 0x0b, 0x69, 0x64, 0x2c, 0x6e, 0x61, 0x6d, 0x65, 0x2c, 0x61,
	0x67, 0x65, 0x92, 0x02, 0x28, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x20, 0x6d, 0x61, 0x73, 0x6b, 0x2c,
	0x20, 0x69, 0x66, 0x20, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2c, 0x20, 0x73, 0x65, 0x6c, 0x65, 0x63,
	0x74, 0x20, 0x61, 0x6c, 0x6c, 0x20, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x2e, 0x52, 0x0a, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73,
	0x69, 0x7a, 0x65, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6e, 0x6f, 0x5f, 0x70, 0x61, 0x67, 0x69, 0x6e,
	0x67, 0x22, 0xa5, 0x04, 0x0a, 0x0c, 0x50, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x27, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x42, 0x0d, 0xba, 0x47, 0x0a, 0x92, 0x02, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x28, 0x0a, 0x05, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x42, 0x12, 0xba, 0x47, 0x0f, 0x92,
	0x02, 0x0c, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x20, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x05,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x34, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x42, 0x0a, 0xba, 0x47, 0x07, 0x92, 0x02,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x38, 0x0a, 0x07, 0x63,
	0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x42, 0x19, 0xba, 0x47,
	0x16, 0x92, 0x02, 0x13, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x20, 0x70, 0x61, 0x67, 0x65,
	0x20, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x48, 0x00, 0x52, 0x07, 0x63, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x74, 0x88, 0x01, 0x01, 0x12, 0x4a, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69,
	0x7a, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x42, 0x27, 0xba, 0x47, 0x24, 0x92, 0x02, 0x21,
	0x6d, 0x61, 0x78, 0x69, 0x6d, 0x75, 0x6d, 0x20, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x20, 0x6f,
	0x66, 0x20, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x20, 0x74, 0x6f, 0x20, 0x72, 0x65, 0x74, 0x75, 0x72,
	0x6e, 0x48, 0x01, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x88, 0x01,
	0x01, 0x12, 0x8d, 0x01, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x63, 0xba, 0x47, 0x60,
	0x92, 0x02, 0x5d, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x20, 0x74, 0x6f, 0x20, 0x72, 0x65, 0x74, 0x72,
	0x69, 0x65, 0x76, 0x65, 0x20, 0x74, 0x68, 0x65, 0x20, 0x6e, 0x65, 0x78, 0x74, 0x20, 0x70, 0x61,
	0x67, 0x65, 0x20, 0x6f, 0x66, 0x20, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x2c, 0x20, 0x6f,
	0x72, 0x20, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x20, 0x69, 0x66, 0x20, 0x74, 0x68, 0x65, 0x72, 0x65,
	0x20, 0x61, 0x72, 0x65, 0x20, 0x6e, 0x6f, 0x20, 0x6d, 0x6f, 0x72, 0x65, 0x20, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x20, 0x69, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x6c, 0x69, 0x73, 0x74,
	0x52, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x5c, 0x0a, 0x05, 0x65, 0x78, 0x74, 0x72, 0x61, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x42, 0x30, 0xba, 0x47, 0x2d, 0x92, 0x02, 0x2a, 0x61, 0x64,
	0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x20, 0x69, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x20, 0x61, 0x62, 0x6f, 0x75, 0x74, 0x20, 0x74, 0x68, 0x69, 0x73, 0x20,
	0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x05, 0x65, 0x78, 0x74, 0x72, 0x61, 0x42,
	0x0a, 0x0a, 0x08, 0x5f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x42, 0x0c, 0x0a, 0x0a, 0x5f,
	0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x42, 0xbc, 0x01, 0x0a, 0x11, 0x63, 0x6f,
	0x6d, 0x2e, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x42,
	0x0f, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x3e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f,
	0x72, 0x69, 0x67, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x76, 0x31, 0xf8, 0x01, 0x01, 0xa2, 0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02, 0x0d, 0x50, 0x61,
	0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0d, 0x50, 0x61,
	0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x19, 0x50, 0x61,
	0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pagination_v1_pagination_proto_rawDescOnce sync.Once
	file_pagination_v1_pagination_proto_rawDescData = file_pagination_v1_pagination_proto_rawDesc
)

func file_pagination_v1_pagination_proto_rawDescGZIP() []byte {
	file_pagination_v1_pagination_proto_rawDescOnce.Do(func() {
		file_pagination_v1_pagination_proto_rawDescData = protoimpl.X.CompressGZIP(file_pagination_v1_pagination_proto_rawDescData)
	})
	return file_pagination_v1_pagination_proto_rawDescData
}

var file_pagination_v1_pagination_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pagination_v1_pagination_proto_goTypes = []any{
	(*PageRequest)(nil),           // 0: pagination.v1.PageRequest
	(*PageResponse)(nil),          // 1: pagination.v1.PageResponse
	(*fieldmaskpb.FieldMask)(nil), // 2: google.protobuf.FieldMask
	(*anypb.Any)(nil),             // 3: google.protobuf.Any
}
var file_pagination_v1_pagination_proto_depIdxs = []int32{
	2, // 0: pagination.v1.PageRequest.field_mask:type_name -> google.protobuf.FieldMask
	3, // 1: pagination.v1.PageResponse.data:type_name -> google.protobuf.Any
	3, // 2: pagination.v1.PageResponse.extra:type_name -> google.protobuf.Any
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_pagination_v1_pagination_proto_init() }
func file_pagination_v1_pagination_proto_init() {
	if File_pagination_v1_pagination_proto != nil {
		return
	}
	file_pagination_v1_pagination_proto_msgTypes[0].OneofWrappers = []any{}
	file_pagination_v1_pagination_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pagination_v1_pagination_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pagination_v1_pagination_proto_goTypes,
		DependencyIndexes: file_pagination_v1_pagination_proto_depIdxs,
		MessageInfos:      file_pagination_v1_pagination_proto_msgTypes,
	}.Build()
	File_pagination_v1_pagination_proto = out.File
	file_pagination_v1_pagination_proto_rawDesc = nil
	file_pagination_v1_pagination_proto_goTypes = nil
	file_pagination_v1_pagination_proto_depIdxs = nil
}
