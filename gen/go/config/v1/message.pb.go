// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: config/v1/message.proto

package configv1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

// Message
type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// name is for register multiple message service
	Name     string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Mqtt     *Message_MQTT     `protobuf:"bytes,3,opt,name=mqtt,proto3" json:"mqtt,omitempty"`
	Kafka    *Message_Kafka    `protobuf:"bytes,4,opt,name=kafka,proto3" json:"kafka,omitempty"`
	Rabbitmq *Message_RabbitMQ `protobuf:"bytes,5,opt,name=rabbitmq,proto3" json:"rabbitmq,omitempty"`
	Activemq *Message_ActiveMQ `protobuf:"bytes,6,opt,name=activemq,proto3" json:"activemq,omitempty"`
	Nats     *Message_NATS     `protobuf:"bytes,7,opt,name=nats,proto3" json:"nats,omitempty"`
	Nsq      *Message_NSQ      `protobuf:"bytes,8,opt,name=nsq,proto3" json:"nsq,omitempty"`
	Pulsar   *Message_Pulsar   `protobuf:"bytes,9,opt,name=pulsar,proto3" json:"pulsar,omitempty"`
	Redis    *Message_Redis    `protobuf:"bytes,10,opt,name=redis,proto3" json:"redis,omitempty"`
	Rocketmq *Message_RocketMQ `protobuf:"bytes,11,opt,name=rocketmq,proto3" json:"rocketmq,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	mi := &file_config_v1_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_message_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_config_v1_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Message) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Message) GetMqtt() *Message_MQTT {
	if x != nil {
		return x.Mqtt
	}
	return nil
}

func (x *Message) GetKafka() *Message_Kafka {
	if x != nil {
		return x.Kafka
	}
	return nil
}

func (x *Message) GetRabbitmq() *Message_RabbitMQ {
	if x != nil {
		return x.Rabbitmq
	}
	return nil
}

func (x *Message) GetActivemq() *Message_ActiveMQ {
	if x != nil {
		return x.Activemq
	}
	return nil
}

func (x *Message) GetNats() *Message_NATS {
	if x != nil {
		return x.Nats
	}
	return nil
}

func (x *Message) GetNsq() *Message_NSQ {
	if x != nil {
		return x.Nsq
	}
	return nil
}

func (x *Message) GetPulsar() *Message_Pulsar {
	if x != nil {
		return x.Pulsar
	}
	return nil
}

func (x *Message) GetRedis() *Message_Redis {
	if x != nil {
		return x.Redis
	}
	return nil
}

func (x *Message) GetRocketmq() *Message_RocketMQ {
	if x != nil {
		return x.Rocketmq
	}
	return nil
}

// MQTT
type Message_MQTT struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Endpoint string `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Codec    string `protobuf:"bytes,2,opt,name=codec,proto3" json:"codec,omitempty"`
}

func (x *Message_MQTT) Reset() {
	*x = Message_MQTT{}
	mi := &file_config_v1_message_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message_MQTT) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message_MQTT) ProtoMessage() {}

func (x *Message_MQTT) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_message_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message_MQTT.ProtoReflect.Descriptor instead.
func (*Message_MQTT) Descriptor() ([]byte, []int) {
	return file_config_v1_message_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Message_MQTT) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Message_MQTT) GetCodec() string {
	if x != nil {
		return x.Codec
	}
	return ""
}

// Kafka
type Message_Kafka struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Endpoint string `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Codec    string `protobuf:"bytes,2,opt,name=codec,proto3" json:"codec,omitempty"`
}

func (x *Message_Kafka) Reset() {
	*x = Message_Kafka{}
	mi := &file_config_v1_message_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message_Kafka) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message_Kafka) ProtoMessage() {}

func (x *Message_Kafka) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_message_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message_Kafka.ProtoReflect.Descriptor instead.
func (*Message_Kafka) Descriptor() ([]byte, []int) {
	return file_config_v1_message_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Message_Kafka) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Message_Kafka) GetCodec() string {
	if x != nil {
		return x.Codec
	}
	return ""
}

// RabbitMQ
type Message_RabbitMQ struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Endpoint string `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Codec    string `protobuf:"bytes,2,opt,name=codec,proto3" json:"codec,omitempty"`
}

func (x *Message_RabbitMQ) Reset() {
	*x = Message_RabbitMQ{}
	mi := &file_config_v1_message_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message_RabbitMQ) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message_RabbitMQ) ProtoMessage() {}

func (x *Message_RabbitMQ) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_message_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message_RabbitMQ.ProtoReflect.Descriptor instead.
func (*Message_RabbitMQ) Descriptor() ([]byte, []int) {
	return file_config_v1_message_proto_rawDescGZIP(), []int{0, 2}
}

func (x *Message_RabbitMQ) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Message_RabbitMQ) GetCodec() string {
	if x != nil {
		return x.Codec
	}
	return ""
}

type Message_ActiveMQ struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Endpoint string `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Codec    string `protobuf:"bytes,2,opt,name=codec,proto3" json:"codec,omitempty"`
}

func (x *Message_ActiveMQ) Reset() {
	*x = Message_ActiveMQ{}
	mi := &file_config_v1_message_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message_ActiveMQ) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message_ActiveMQ) ProtoMessage() {}

func (x *Message_ActiveMQ) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_message_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message_ActiveMQ.ProtoReflect.Descriptor instead.
func (*Message_ActiveMQ) Descriptor() ([]byte, []int) {
	return file_config_v1_message_proto_rawDescGZIP(), []int{0, 3}
}

func (x *Message_ActiveMQ) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Message_ActiveMQ) GetCodec() string {
	if x != nil {
		return x.Codec
	}
	return ""
}

type Message_NATS struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Endpoint string `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Codec    string `protobuf:"bytes,2,opt,name=codec,proto3" json:"codec,omitempty"`
}

func (x *Message_NATS) Reset() {
	*x = Message_NATS{}
	mi := &file_config_v1_message_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message_NATS) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message_NATS) ProtoMessage() {}

func (x *Message_NATS) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_message_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message_NATS.ProtoReflect.Descriptor instead.
func (*Message_NATS) Descriptor() ([]byte, []int) {
	return file_config_v1_message_proto_rawDescGZIP(), []int{0, 4}
}

func (x *Message_NATS) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Message_NATS) GetCodec() string {
	if x != nil {
		return x.Codec
	}
	return ""
}

type Message_NSQ struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Endpoint string `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Codec    string `protobuf:"bytes,2,opt,name=codec,proto3" json:"codec,omitempty"`
}

func (x *Message_NSQ) Reset() {
	*x = Message_NSQ{}
	mi := &file_config_v1_message_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message_NSQ) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message_NSQ) ProtoMessage() {}

func (x *Message_NSQ) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_message_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message_NSQ.ProtoReflect.Descriptor instead.
func (*Message_NSQ) Descriptor() ([]byte, []int) {
	return file_config_v1_message_proto_rawDescGZIP(), []int{0, 5}
}

func (x *Message_NSQ) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Message_NSQ) GetCodec() string {
	if x != nil {
		return x.Codec
	}
	return ""
}

type Message_Pulsar struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Endpoint string `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Codec    string `protobuf:"bytes,2,opt,name=codec,proto3" json:"codec,omitempty"`
}

func (x *Message_Pulsar) Reset() {
	*x = Message_Pulsar{}
	mi := &file_config_v1_message_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message_Pulsar) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message_Pulsar) ProtoMessage() {}

func (x *Message_Pulsar) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_message_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message_Pulsar.ProtoReflect.Descriptor instead.
func (*Message_Pulsar) Descriptor() ([]byte, []int) {
	return file_config_v1_message_proto_rawDescGZIP(), []int{0, 6}
}

func (x *Message_Pulsar) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Message_Pulsar) GetCodec() string {
	if x != nil {
		return x.Codec
	}
	return ""
}

type Message_Redis struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Endpoint string `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Codec    string `protobuf:"bytes,2,opt,name=codec,proto3" json:"codec,omitempty"`
}

func (x *Message_Redis) Reset() {
	*x = Message_Redis{}
	mi := &file_config_v1_message_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message_Redis) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message_Redis) ProtoMessage() {}

func (x *Message_Redis) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_message_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message_Redis.ProtoReflect.Descriptor instead.
func (*Message_Redis) Descriptor() ([]byte, []int) {
	return file_config_v1_message_proto_rawDescGZIP(), []int{0, 7}
}

func (x *Message_Redis) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Message_Redis) GetCodec() string {
	if x != nil {
		return x.Codec
	}
	return ""
}

type Message_RocketMQ struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Endpoint         string   `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Codec            string   `protobuf:"bytes,2,opt,name=codec,proto3" json:"codec,omitempty"`
	EnableTrace      bool     `protobuf:"varint,3,opt,name=enable_trace,proto3" json:"enable_trace,omitempty"`
	NameServers      []string `protobuf:"bytes,4,rep,name=name_servers,proto3" json:"name_servers,omitempty"`
	NameServerDomain string   `protobuf:"bytes,5,opt,name=name_server_domain,proto3" json:"name_server_domain,omitempty"`
	AccessKey        string   `protobuf:"bytes,6,opt,name=access_key,proto3" json:"access_key,omitempty"`
	SecretKey        string   `protobuf:"bytes,7,opt,name=secret_key,proto3" json:"secret_key,omitempty"`
	SecurityToken    string   `protobuf:"bytes,8,opt,name=security_token,proto3" json:"security_token,omitempty"`
	Namespace        string   `protobuf:"bytes,9,opt,name=namespace,proto3" json:"namespace,omitempty"`
	InstanceName     string   `protobuf:"bytes,10,opt,name=instance_name,proto3" json:"instance_name,omitempty"`
	GroupName        string   `protobuf:"bytes,11,opt,name=group_name,proto3" json:"group_name,omitempty"`
}

func (x *Message_RocketMQ) Reset() {
	*x = Message_RocketMQ{}
	mi := &file_config_v1_message_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message_RocketMQ) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message_RocketMQ) ProtoMessage() {}

func (x *Message_RocketMQ) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_message_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message_RocketMQ.ProtoReflect.Descriptor instead.
func (*Message_RocketMQ) Descriptor() ([]byte, []int) {
	return file_config_v1_message_proto_rawDescGZIP(), []int{0, 8}
}

func (x *Message_RocketMQ) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Message_RocketMQ) GetCodec() string {
	if x != nil {
		return x.Codec
	}
	return ""
}

func (x *Message_RocketMQ) GetEnableTrace() bool {
	if x != nil {
		return x.EnableTrace
	}
	return false
}

func (x *Message_RocketMQ) GetNameServers() []string {
	if x != nil {
		return x.NameServers
	}
	return nil
}

func (x *Message_RocketMQ) GetNameServerDomain() string {
	if x != nil {
		return x.NameServerDomain
	}
	return ""
}

func (x *Message_RocketMQ) GetAccessKey() string {
	if x != nil {
		return x.AccessKey
	}
	return ""
}

func (x *Message_RocketMQ) GetSecretKey() string {
	if x != nil {
		return x.SecretKey
	}
	return ""
}

func (x *Message_RocketMQ) GetSecurityToken() string {
	if x != nil {
		return x.SecurityToken
	}
	return ""
}

func (x *Message_RocketMQ) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *Message_RocketMQ) GetInstanceName() string {
	if x != nil {
		return x.InstanceName
	}
	return ""
}

func (x *Message_RocketMQ) GetGroupName() string {
	if x != nil {
		return x.GroupName
	}
	return ""
}

var File_config_v1_message_proto protoreflect.FileDescriptor

var file_config_v1_message_proto_rawDesc = []byte{
	0x0a, 0x17, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x76, 0x31, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa3, 0x0b,
	0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x64, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x50, 0xfa, 0x42, 0x4d, 0x72, 0x4b, 0x52, 0x04,
	0x6e, 0x6f, 0x6e, 0x65, 0x52, 0x04, 0x6d, 0x71, 0x74, 0x74, 0x52, 0x05, 0x6b, 0x61, 0x66, 0x6b,
	0x61, 0x52, 0x08, 0x72, 0x61, 0x62, 0x62, 0x69, 0x74, 0x6d, 0x71, 0x52, 0x08, 0x61, 0x63, 0x74,
	0x69, 0x76, 0x65, 0x6d, 0x71, 0x52, 0x04, 0x6e, 0x61, 0x74, 0x73, 0x52, 0x03, 0x6e, 0x73, 0x71,
	0x52, 0x06, 0x70, 0x75, 0x6c, 0x73, 0x61, 0x72, 0x52, 0x05, 0x72, 0x65, 0x64, 0x69, 0x73, 0x52,
	0x08, 0x72, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x6d, 0x71, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x6d, 0x71, 0x74, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x4d, 0x51, 0x54, 0x54, 0x52, 0x04, 0x6d, 0x71, 0x74, 0x74,
	0x12, 0x2e, 0x0a, 0x05, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x18, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x2e, 0x4b, 0x61, 0x66, 0x6b, 0x61, 0x52, 0x05, 0x6b, 0x61, 0x66, 0x6b, 0x61,
	0x12, 0x37, 0x0a, 0x08, 0x72, 0x61, 0x62, 0x62, 0x69, 0x74, 0x6d, 0x71, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x52, 0x61, 0x62, 0x62, 0x69, 0x74, 0x4d, 0x51, 0x52,
	0x08, 0x72, 0x61, 0x62, 0x62, 0x69, 0x74, 0x6d, 0x71, 0x12, 0x37, 0x0a, 0x08, 0x61, 0x63, 0x74,
	0x69, 0x76, 0x65, 0x6d, 0x71, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x4d, 0x51, 0x52, 0x08, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65,
	0x6d, 0x71, 0x12, 0x2b, 0x0a, 0x04, 0x6e, 0x61, 0x74, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2e, 0x4e, 0x41, 0x54, 0x53, 0x52, 0x04, 0x6e, 0x61, 0x74, 0x73, 0x12,
	0x28, 0x0a, 0x03, 0x6e, 0x73, 0x71, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x4e, 0x53, 0x51, 0x52, 0x03, 0x6e, 0x73, 0x71, 0x12, 0x31, 0x0a, 0x06, 0x70, 0x75, 0x6c,
	0x73, 0x61, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x50, 0x75,
	0x6c, 0x73, 0x61, 0x72, 0x52, 0x06, 0x70, 0x75, 0x6c, 0x73, 0x61, 0x72, 0x12, 0x2e, 0x0a, 0x05,
	0x72, 0x65, 0x64, 0x69, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x52, 0x65, 0x64, 0x69, 0x73, 0x52, 0x05, 0x72, 0x65, 0x64, 0x69, 0x73, 0x12, 0x37, 0x0a, 0x08,
	0x72, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x6d, 0x71, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b,
	0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x52, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x4d, 0x51, 0x52, 0x08, 0x72, 0x6f, 0x63,
	0x6b, 0x65, 0x74, 0x6d, 0x71, 0x1a, 0x38, 0x0a, 0x04, 0x4d, 0x51, 0x54, 0x54, 0x12, 0x1a, 0x0a,
	0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x64,
	0x65, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x1a,
	0x39, 0x0a, 0x05, 0x4b, 0x61, 0x66, 0x6b, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x1a, 0x3c, 0x0a, 0x08, 0x52, 0x61,
	0x62, 0x62, 0x69, 0x74, 0x4d, 0x51, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x1a, 0x3c, 0x0a, 0x08, 0x41, 0x63, 0x74, 0x69,
	0x76, 0x65, 0x4d, 0x51, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x1a, 0x38, 0x0a, 0x04, 0x4e, 0x41, 0x54, 0x53, 0x12, 0x1a,
	0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f,
	0x64, 0x65, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6f, 0x64, 0x65, 0x63,
	0x1a, 0x37, 0x0a, 0x03, 0x4e, 0x53, 0x51, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x1a, 0x3a, 0x0a, 0x06, 0x50, 0x75, 0x6c,
	0x73, 0x61, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x63, 0x6f, 0x64, 0x65, 0x63, 0x1a, 0x39, 0x0a, 0x05, 0x52, 0x65, 0x64, 0x69, 0x73, 0x12, 0x1a,
	0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f,
	0x64, 0x65, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6f, 0x64, 0x65, 0x63,
	0x1a, 0x80, 0x03, 0x0a, 0x08, 0x52, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x4d, 0x51, 0x12, 0x1a, 0x0a,
	0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x64,
	0x65, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x12,
	0x22, 0x0a, 0x0c, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x74, 0x72,
	0x61, 0x63, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x6e, 0x61, 0x6d, 0x65, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x73, 0x12, 0x2e, 0x0a, 0x12, 0x6e, 0x61, 0x6d, 0x65, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x12, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x5f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x5f, 0x6b, 0x65, 0x79, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x65, 0x63, 0x72, 0x65,
	0x74, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x65, 0x63,
	0x72, 0x65, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x12, 0x26, 0x0a, 0x0e, 0x73, 0x65, 0x63, 0x75, 0x72,
	0x69, 0x74, 0x79, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12,
	0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x24, 0x0a,
	0x0d, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x42, 0x9d, 0x01, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6f, 0x72, 0x69, 0x67, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x72, 0x75, 0x6e, 0x74,
	0x69, 0x6d, 0x65, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x76, 0x31, 0xf8, 0x01, 0x01,
	0xa2, 0x02, 0x03, 0x43, 0x58, 0x58, 0xaa, 0x02, 0x09, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x56, 0x31, 0xca, 0x02, 0x09, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5c, 0x56, 0x31, 0xe2, 0x02,
	0x15, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0a, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x3a,
	0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_config_v1_message_proto_rawDescOnce sync.Once
	file_config_v1_message_proto_rawDescData = file_config_v1_message_proto_rawDesc
)

func file_config_v1_message_proto_rawDescGZIP() []byte {
	file_config_v1_message_proto_rawDescOnce.Do(func() {
		file_config_v1_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_config_v1_message_proto_rawDescData)
	})
	return file_config_v1_message_proto_rawDescData
}

var file_config_v1_message_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_config_v1_message_proto_goTypes = []any{
	(*Message)(nil),          // 0: config.v1.Message
	(*Message_MQTT)(nil),     // 1: config.v1.Message.MQTT
	(*Message_Kafka)(nil),    // 2: config.v1.Message.Kafka
	(*Message_RabbitMQ)(nil), // 3: config.v1.Message.RabbitMQ
	(*Message_ActiveMQ)(nil), // 4: config.v1.Message.ActiveMQ
	(*Message_NATS)(nil),     // 5: config.v1.Message.NATS
	(*Message_NSQ)(nil),      // 6: config.v1.Message.NSQ
	(*Message_Pulsar)(nil),   // 7: config.v1.Message.Pulsar
	(*Message_Redis)(nil),    // 8: config.v1.Message.Redis
	(*Message_RocketMQ)(nil), // 9: config.v1.Message.RocketMQ
}
var file_config_v1_message_proto_depIdxs = []int32{
	1, // 0: config.v1.Message.mqtt:type_name -> config.v1.Message.MQTT
	2, // 1: config.v1.Message.kafka:type_name -> config.v1.Message.Kafka
	3, // 2: config.v1.Message.rabbitmq:type_name -> config.v1.Message.RabbitMQ
	4, // 3: config.v1.Message.activemq:type_name -> config.v1.Message.ActiveMQ
	5, // 4: config.v1.Message.nats:type_name -> config.v1.Message.NATS
	6, // 5: config.v1.Message.nsq:type_name -> config.v1.Message.NSQ
	7, // 6: config.v1.Message.pulsar:type_name -> config.v1.Message.Pulsar
	8, // 7: config.v1.Message.redis:type_name -> config.v1.Message.Redis
	9, // 8: config.v1.Message.rocketmq:type_name -> config.v1.Message.RocketMQ
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_config_v1_message_proto_init() }
func file_config_v1_message_proto_init() {
	if File_config_v1_message_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_config_v1_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_v1_message_proto_goTypes,
		DependencyIndexes: file_config_v1_message_proto_depIdxs,
		MessageInfos:      file_config_v1_message_proto_msgTypes,
	}.Build()
	File_config_v1_message_proto = out.File
	file_config_v1_message_proto_rawDesc = nil
	file_config_v1_message_proto_goTypes = nil
	file_config_v1_message_proto_depIdxs = nil
}
