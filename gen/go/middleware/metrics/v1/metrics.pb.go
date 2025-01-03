// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: middleware/metrics/v1/metrics.proto

package metricsv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Type of indicator (e.g. counter, timer, histogram, etc.)
type UserMetric_MetricType int32

const (
	UserMetric_METRIC_TYPE_UNSPECIFIED UserMetric_MetricType = 0
	UserMetric_METRIC_TYPE_COUNTER     UserMetric_MetricType = 1
	UserMetric_METRIC_TYPE_GAUGE       UserMetric_MetricType = 2
	UserMetric_METRIC_TYPE_HISTOGRAM   UserMetric_MetricType = 3
	UserMetric_METRIC_TYPE_SUMMARY     UserMetric_MetricType = 4
)

// Enum value maps for UserMetric_MetricType.
var (
	UserMetric_MetricType_name = map[int32]string{
		0: "METRIC_TYPE_UNSPECIFIED",
		1: "METRIC_TYPE_COUNTER",
		2: "METRIC_TYPE_GAUGE",
		3: "METRIC_TYPE_HISTOGRAM",
		4: "METRIC_TYPE_SUMMARY",
	}
	UserMetric_MetricType_value = map[string]int32{
		"METRIC_TYPE_UNSPECIFIED": 0,
		"METRIC_TYPE_COUNTER":     1,
		"METRIC_TYPE_GAUGE":       2,
		"METRIC_TYPE_HISTOGRAM":   3,
		"METRIC_TYPE_SUMMARY":     4,
	}
)

func (x UserMetric_MetricType) Enum() *UserMetric_MetricType {
	p := new(UserMetric_MetricType)
	*p = x
	return p
}

func (x UserMetric_MetricType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserMetric_MetricType) Descriptor() protoreflect.EnumDescriptor {
	return file_middleware_metrics_v1_metrics_proto_enumTypes[0].Descriptor()
}

func (UserMetric_MetricType) Type() protoreflect.EnumType {
	return &file_middleware_metrics_v1_metrics_proto_enumTypes[0]
}

func (x UserMetric_MetricType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserMetric_MetricType.Descriptor instead.
func (UserMetric_MetricType) EnumDescriptor() ([]byte, []int) {
	return file_middleware_metrics_v1_metrics_proto_rawDescGZIP(), []int{0, 0}
}

type UserMetric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Timestamp: indicates the time of indicator data
	Timestamp int64 `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Indicator name
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Indicator value
	Value float64 `protobuf:"fixed64,3,opt,name=value,proto3" json:"value,omitempty"`
	// Indicator label for classification or filtering
	Labels map[string]string `protobuf:"bytes,4,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Indicator unit
	Unit string                `protobuf:"bytes,5,opt,name=unit,proto3" json:"unit,omitempty"`
	Type UserMetric_MetricType `protobuf:"varint,6,opt,name=type,proto3,enum=middleware.metrics.v1.UserMetric_MetricType" json:"type,omitempty"`
	// Description of indicators
	Description string `protobuf:"bytes,7,opt,name=description,proto3" json:"description,omitempty"`
	// Indicator context information
	Context string `protobuf:"bytes,8,opt,name=context,proto3" json:"context,omitempty"`
	// Additional information for metrics that can be used to store arbitrary metadata
	Metadata map[string]string `protobuf:"bytes,9,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *UserMetric) Reset() {
	*x = UserMetric{}
	mi := &file_middleware_metrics_v1_metrics_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserMetric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserMetric) ProtoMessage() {}

func (x *UserMetric) ProtoReflect() protoreflect.Message {
	mi := &file_middleware_metrics_v1_metrics_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserMetric.ProtoReflect.Descriptor instead.
func (*UserMetric) Descriptor() ([]byte, []int) {
	return file_middleware_metrics_v1_metrics_proto_rawDescGZIP(), []int{0}
}

func (x *UserMetric) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *UserMetric) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UserMetric) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *UserMetric) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *UserMetric) GetUnit() string {
	if x != nil {
		return x.Unit
	}
	return ""
}

func (x *UserMetric) GetType() UserMetric_MetricType {
	if x != nil {
		return x.Type
	}
	return UserMetric_METRIC_TYPE_UNSPECIFIED
}

func (x *UserMetric) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *UserMetric) GetContext() string {
	if x != nil {
		return x.Context
	}
	return ""
}

func (x *UserMetric) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

// Metrics
type Metrics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Enabled bool `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	// Add a list of supported metrics for enabling or disabling specific metrics
	SupportedMetrics []string `protobuf:"bytes,5,rep,name=supported_metrics,proto3" json:"supported_metrics,omitempty"`
	// Repeated field for user-defined metrics
	UserMetrics []*UserMetric `protobuf:"bytes,6,rep,name=user_metrics,proto3" json:"user_metrics,omitempty"`
}

func (x *Metrics) Reset() {
	*x = Metrics{}
	mi := &file_middleware_metrics_v1_metrics_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Metrics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metrics) ProtoMessage() {}

func (x *Metrics) ProtoReflect() protoreflect.Message {
	mi := &file_middleware_metrics_v1_metrics_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metrics.ProtoReflect.Descriptor instead.
func (*Metrics) Descriptor() ([]byte, []int) {
	return file_middleware_metrics_v1_metrics_proto_rawDescGZIP(), []int{1}
}

func (x *Metrics) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *Metrics) GetSupportedMetrics() []string {
	if x != nil {
		return x.SupportedMetrics
	}
	return nil
}

func (x *Metrics) GetUserMetrics() []*UserMetric {
	if x != nil {
		return x.UserMetrics
	}
	return nil
}

var File_middleware_metrics_v1_metrics_proto protoreflect.FileDescriptor

var file_middleware_metrics_v1_metrics_proto_rawDesc = []byte{
	0x0a, 0x23, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2f, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72,
	0x65, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x82, 0x05,
	0x0a, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x1c, 0x0a, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x45, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72,
	0x65, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x75,
	0x6e, 0x69, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x12,
	0x40, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2c, 0x2e,
	0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x12, 0x4b, 0x0a,
	0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x2f, 0x2e, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2e, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61,
	0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x8d, 0x01, 0x0a, 0x0a, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x1b, 0x0a, 0x17, 0x4d, 0x45, 0x54, 0x52, 0x49, 0x43, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x17,
	0x0a, 0x13, 0x4d, 0x45, 0x54, 0x52, 0x49, 0x43, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x4f,
	0x55, 0x4e, 0x54, 0x45, 0x52, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x4d, 0x45, 0x54, 0x52, 0x49,
	0x43, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x47, 0x41, 0x55, 0x47, 0x45, 0x10, 0x02, 0x12, 0x19,
	0x0a, 0x15, 0x4d, 0x45, 0x54, 0x52, 0x49, 0x43, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x48, 0x49,
	0x53, 0x54, 0x4f, 0x47, 0x52, 0x41, 0x4d, 0x10, 0x03, 0x12, 0x17, 0x0a, 0x13, 0x4d, 0x45, 0x54,
	0x52, 0x49, 0x43, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x55, 0x4d, 0x4d, 0x41, 0x52, 0x59,
	0x10, 0x04, 0x22, 0x98, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x18,
	0x0a, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x2c, 0x0a, 0x11, 0x73, 0x75, 0x70, 0x70,
	0x6f, 0x72, 0x74, 0x65, 0x64, 0x5f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x05, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x11, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x5f, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x45, 0x0a, 0x0c, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x6d,
	0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52,
	0x0c, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x42, 0xe7, 0x01,
	0x0a, 0x19, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65,
	0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x43, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x72, 0x69, 0x67, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2f, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f,
	0x2f, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2f, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x76, 0x31,
	0xf8, 0x01, 0x01, 0xa2, 0x02, 0x03, 0x4d, 0x4d, 0x58, 0xaa, 0x02, 0x15, 0x4d, 0x69, 0x64, 0x64,
	0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x56,
	0x31, 0xca, 0x02, 0x15, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x5c, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x21, 0x4d, 0x69, 0x64, 0x64,
	0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x5c, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x5c, 0x56,
	0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x17,
	0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x3a, 0x3a, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_middleware_metrics_v1_metrics_proto_rawDescOnce sync.Once
	file_middleware_metrics_v1_metrics_proto_rawDescData = file_middleware_metrics_v1_metrics_proto_rawDesc
)

func file_middleware_metrics_v1_metrics_proto_rawDescGZIP() []byte {
	file_middleware_metrics_v1_metrics_proto_rawDescOnce.Do(func() {
		file_middleware_metrics_v1_metrics_proto_rawDescData = protoimpl.X.CompressGZIP(file_middleware_metrics_v1_metrics_proto_rawDescData)
	})
	return file_middleware_metrics_v1_metrics_proto_rawDescData
}

var file_middleware_metrics_v1_metrics_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_middleware_metrics_v1_metrics_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_middleware_metrics_v1_metrics_proto_goTypes = []any{
	(UserMetric_MetricType)(0), // 0: middleware.metrics.v1.UserMetric.MetricType
	(*UserMetric)(nil),         // 1: middleware.metrics.v1.UserMetric
	(*Metrics)(nil),            // 2: middleware.metrics.v1.Metrics
	nil,                        // 3: middleware.metrics.v1.UserMetric.LabelsEntry
	nil,                        // 4: middleware.metrics.v1.UserMetric.MetadataEntry
}
var file_middleware_metrics_v1_metrics_proto_depIdxs = []int32{
	3, // 0: middleware.metrics.v1.UserMetric.labels:type_name -> middleware.metrics.v1.UserMetric.LabelsEntry
	0, // 1: middleware.metrics.v1.UserMetric.type:type_name -> middleware.metrics.v1.UserMetric.MetricType
	4, // 2: middleware.metrics.v1.UserMetric.metadata:type_name -> middleware.metrics.v1.UserMetric.MetadataEntry
	1, // 3: middleware.metrics.v1.Metrics.user_metrics:type_name -> middleware.metrics.v1.UserMetric
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_middleware_metrics_v1_metrics_proto_init() }
func file_middleware_metrics_v1_metrics_proto_init() {
	if File_middleware_metrics_v1_metrics_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_middleware_metrics_v1_metrics_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_middleware_metrics_v1_metrics_proto_goTypes,
		DependencyIndexes: file_middleware_metrics_v1_metrics_proto_depIdxs,
		EnumInfos:         file_middleware_metrics_v1_metrics_proto_enumTypes,
		MessageInfos:      file_middleware_metrics_v1_metrics_proto_msgTypes,
	}.Build()
	File_middleware_metrics_v1_metrics_proto = out.File
	file_middleware_metrics_v1_metrics_proto_rawDesc = nil
	file_middleware_metrics_v1_metrics_proto_goTypes = nil
	file_middleware_metrics_v1_metrics_proto_depIdxs = nil
}
