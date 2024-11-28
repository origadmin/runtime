// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: config/v1/task.proto

package configv1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
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

// Task config
type Task struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      string          `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Name      string          `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Asynq     *Task_Asynq     `protobuf:"bytes,3,opt,name=asynq,proto3" json:"asynq,omitempty"`
	Machinery *Task_Machinery `protobuf:"bytes,4,opt,name=machinery,proto3" json:"machinery,omitempty"`
	Cron      *Task_Cron      `protobuf:"bytes,5,opt,name=cron,proto3" json:"cron,omitempty"`
}

func (x *Task) Reset() {
	*x = Task{}
	mi := &file_config_v1_task_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_task_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task.ProtoReflect.Descriptor instead.
func (*Task) Descriptor() ([]byte, []int) {
	return file_config_v1_task_proto_rawDescGZIP(), []int{0}
}

func (x *Task) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Task) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Task) GetAsynq() *Task_Asynq {
	if x != nil {
		return x.Asynq
	}
	return nil
}

func (x *Task) GetMachinery() *Task_Machinery {
	if x != nil {
		return x.Machinery
	}
	return nil
}

func (x *Task) GetCron() *Task_Cron {
	if x != nil {
		return x.Cron
	}
	return nil
}

// Asynq config
type Task_Asynq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// endpoint is peer network address
	Endpoint string `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	// login password
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	// database index
	Db int32 `protobuf:"varint,3,opt,name=db,proto3" json:"db,omitempty"`
	// timezone location
	Location string `protobuf:"bytes,4,opt,name=location,proto3" json:"location,omitempty"`
}

func (x *Task_Asynq) Reset() {
	*x = Task_Asynq{}
	mi := &file_config_v1_task_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Task_Asynq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task_Asynq) ProtoMessage() {}

func (x *Task_Asynq) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_task_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task_Asynq.ProtoReflect.Descriptor instead.
func (*Task_Asynq) Descriptor() ([]byte, []int) {
	return file_config_v1_task_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Task_Asynq) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Task_Asynq) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Task_Asynq) GetDb() int32 {
	if x != nil {
		return x.Db
	}
	return 0
}

func (x *Task_Asynq) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

// Machinery config
type Task_Machinery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// brokers address, which can be specified as Redis, AMQP, or AWS SQS according to the actual storage medium used
	Brokers []string `protobuf:"bytes,1,rep,name=brokers,proto3" json:"brokers,omitempty"`
	// backends configures the media for storing results. The value can be Redis, memcached, or mongodb as required
	Backends []string `protobuf:"bytes,2,rep,name=backends,proto3" json:"backends,omitempty"`
}

func (x *Task_Machinery) Reset() {
	*x = Task_Machinery{}
	mi := &file_config_v1_task_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Task_Machinery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task_Machinery) ProtoMessage() {}

func (x *Task_Machinery) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_task_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task_Machinery.ProtoReflect.Descriptor instead.
func (*Task_Machinery) Descriptor() ([]byte, []int) {
	return file_config_v1_task_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Task_Machinery) GetBrokers() []string {
	if x != nil {
		return x.Brokers
	}
	return nil
}

func (x *Task_Machinery) GetBackends() []string {
	if x != nil {
		return x.Backends
	}
	return nil
}

// Cron config
type Task_Cron struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// addr is peer network address
	Addr string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
}

func (x *Task_Cron) Reset() {
	*x = Task_Cron{}
	mi := &file_config_v1_task_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Task_Cron) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task_Cron) ProtoMessage() {}

func (x *Task_Cron) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_task_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task_Cron.ProtoReflect.Descriptor instead.
func (*Task_Cron) Descriptor() ([]byte, []int) {
	return file_config_v1_task_proto_rawDescGZIP(), []int{0, 2}
}

func (x *Task_Cron) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

var File_config_v1_task_proto protoreflect.FileDescriptor

var file_config_v1_task_proto_rawDesc = []byte{
	0x0a, 0x14, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x61, 0x73, 0x6b,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76,
	0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xaf,
	0x03, 0x0a, 0x04, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x37, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x23, 0xba, 0x48, 0x20, 0x72, 0x1e, 0x52, 0x04, 0x6e, 0x6f,
	0x6e, 0x65, 0x52, 0x05, 0x61, 0x73, 0x79, 0x6e, 0x71, 0x52, 0x09, 0x6d, 0x61, 0x63, 0x68, 0x69,
	0x6e, 0x65, 0x72, 0x79, 0x52, 0x04, 0x63, 0x72, 0x6f, 0x6e, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2b, 0x0a, 0x05, 0x61, 0x73, 0x79, 0x6e, 0x71, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e,
	0x54, 0x61, 0x73, 0x6b, 0x2e, 0x41, 0x73, 0x79, 0x6e, 0x71, 0x52, 0x05, 0x61, 0x73, 0x79, 0x6e,
	0x71, 0x12, 0x37, 0x0a, 0x09, 0x6d, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x72, 0x79, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31,
	0x2e, 0x54, 0x61, 0x73, 0x6b, 0x2e, 0x4d, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x72, 0x79, 0x52,
	0x09, 0x6d, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x72, 0x79, 0x12, 0x28, 0x0a, 0x04, 0x63, 0x72,
	0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x2e, 0x43, 0x72, 0x6f, 0x6e, 0x52, 0x04,
	0x63, 0x72, 0x6f, 0x6e, 0x1a, 0x6b, 0x0a, 0x05, 0x41, 0x73, 0x79, 0x6e, 0x71, 0x12, 0x1a, 0x0a,
	0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x64, 0x62, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x02, 0x64, 0x62, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x1a, 0x41, 0x0a, 0x09, 0x4d, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x72, 0x79, 0x12, 0x18,
	0x0a, 0x07, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x07, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x61, 0x63, 0x6b,
	0x65, 0x6e, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x62, 0x61, 0x63, 0x6b,
	0x65, 0x6e, 0x64, 0x73, 0x1a, 0x1a, 0x0a, 0x04, 0x43, 0x72, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04,
	0x61, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72,
	0x42, 0xa3, 0x01, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x76, 0x31, 0x42, 0x09, 0x54, 0x61, 0x73, 0x6b, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x72, 0x69, 0x67,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x74, 0x6f, 0x6f, 0x6c, 0x6b, 0x69, 0x74, 0x73, 0x2f, 0x72,
	0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x76, 0x31,
	0xf8, 0x01, 0x01, 0xa2, 0x02, 0x03, 0x43, 0x58, 0x58, 0xaa, 0x02, 0x09, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x09, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5c, 0x56,
	0x31, 0xe2, 0x02, 0x15, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0a, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_config_v1_task_proto_rawDescOnce sync.Once
	file_config_v1_task_proto_rawDescData = file_config_v1_task_proto_rawDesc
)

func file_config_v1_task_proto_rawDescGZIP() []byte {
	file_config_v1_task_proto_rawDescOnce.Do(func() {
		file_config_v1_task_proto_rawDescData = protoimpl.X.CompressGZIP(file_config_v1_task_proto_rawDescData)
	})
	return file_config_v1_task_proto_rawDescData
}

var file_config_v1_task_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_config_v1_task_proto_goTypes = []any{
	(*Task)(nil),           // 0: config.v1.Task
	(*Task_Asynq)(nil),     // 1: config.v1.Task.Asynq
	(*Task_Machinery)(nil), // 2: config.v1.Task.Machinery
	(*Task_Cron)(nil),      // 3: config.v1.Task.Cron
}
var file_config_v1_task_proto_depIdxs = []int32{
	1, // 0: config.v1.Task.asynq:type_name -> config.v1.Task.Asynq
	2, // 1: config.v1.Task.machinery:type_name -> config.v1.Task.Machinery
	3, // 2: config.v1.Task.cron:type_name -> config.v1.Task.Cron
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_config_v1_task_proto_init() }
func file_config_v1_task_proto_init() {
	if File_config_v1_task_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_config_v1_task_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_v1_task_proto_goTypes,
		DependencyIndexes: file_config_v1_task_proto_depIdxs,
		MessageInfos:      file_config_v1_task_proto_msgTypes,
	}.Build()
	File_config_v1_task_proto = out.File
	file_config_v1_task_proto_rawDesc = nil
	file_config_v1_task_proto_goTypes = nil
	file_config_v1_task_proto_depIdxs = nil
}