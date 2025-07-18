// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: config/v1/task.proto

package configv1

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

// Task config
type Task struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Type          string                 `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Asynq         *Task_Asynq            `protobuf:"bytes,3,opt,name=asynq,proto3" json:"asynq,omitempty"`
	Machinery     *Task_Machinery        `protobuf:"bytes,4,opt,name=machinery,proto3" json:"machinery,omitempty"`
	Cron          *Task_Cron             `protobuf:"bytes,5,opt,name=cron,proto3" json:"cron,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
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
	state protoimpl.MessageState `protogen:"open.v1"`
	// endpoint is peer network address
	Endpoint string `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	// login password
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	// database index
	Db int32 `protobuf:"varint,3,opt,name=db,proto3" json:"db,omitempty"`
	// timezone location
	Location      string `protobuf:"bytes,4,opt,name=location,proto3" json:"location,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
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
	state protoimpl.MessageState `protogen:"open.v1"`
	// brokers address, which can be specified as Redis, AMQP, or AWS SQS according to the actual storage medium used
	Brokers []string `protobuf:"bytes,1,rep,name=brokers,proto3" json:"brokers,omitempty"`
	// backends configures the media for storing results. The value can be Redis, memcached, or mongodb as required
	Backends      []string `protobuf:"bytes,2,rep,name=backends,proto3" json:"backends,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
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
	state protoimpl.MessageState `protogen:"open.v1"`
	// addr is peer network address
	Addr          string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
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

const file_config_v1_task_proto_rawDesc = "" +
	"\n" +
	"\x14config/v1/task.proto\x12\tconfig.v1\x1a\x17validate/validate.proto\"\xaf\x03\n" +
	"\x04Task\x127\n" +
	"\x04type\x18\x01 \x01(\tB#\xfaB r\x1eR\x04noneR\x05asynqR\tmachineryR\x04cronR\x04type\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12+\n" +
	"\x05asynq\x18\x03 \x01(\v2\x15.config.v1.Task.AsynqR\x05asynq\x127\n" +
	"\tmachinery\x18\x04 \x01(\v2\x19.config.v1.Task.MachineryR\tmachinery\x12(\n" +
	"\x04cron\x18\x05 \x01(\v2\x14.config.v1.Task.CronR\x04cron\x1ak\n" +
	"\x05Asynq\x12\x1a\n" +
	"\bendpoint\x18\x01 \x01(\tR\bendpoint\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\x12\x0e\n" +
	"\x02db\x18\x03 \x01(\x05R\x02db\x12\x1a\n" +
	"\blocation\x18\x04 \x01(\tR\blocation\x1aA\n" +
	"\tMachinery\x12\x18\n" +
	"\abrokers\x18\x01 \x03(\tR\abrokers\x12\x1a\n" +
	"\bbackends\x18\x02 \x03(\tR\bbackends\x1a\x1a\n" +
	"\x04Cron\x12\x12\n" +
	"\x04addr\x18\x01 \x01(\tR\x04addrB\x9e\x01\n" +
	"\rcom.config.v1B\tTaskProtoP\x01Z:github.com/origadmin/runtime/api/gen/go/config/v1;configv1\xf8\x01\x01\xa2\x02\x03CXX\xaa\x02\tConfig.V1\xca\x02\tConfig\\V1\xe2\x02\x15Config\\V1\\GPBMetadata\xea\x02\n" +
	"Config::V1b\x06proto3"

var (
	file_config_v1_task_proto_rawDescOnce sync.Once
	file_config_v1_task_proto_rawDescData []byte
)

func file_config_v1_task_proto_rawDescGZIP() []byte {
	file_config_v1_task_proto_rawDescOnce.Do(func() {
		file_config_v1_task_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_config_v1_task_proto_rawDesc), len(file_config_v1_task_proto_rawDesc)))
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
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_config_v1_task_proto_rawDesc), len(file_config_v1_task_proto_rawDesc)),
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
	file_config_v1_task_proto_goTypes = nil
	file_config_v1_task_proto_depIdxs = nil
}
