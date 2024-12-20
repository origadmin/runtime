// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: config/v1/gateway.proto

package configv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Protocol int32

const (
	Protocol_UNSPECIFIED Protocol = 0
	Protocol_HTTP        Protocol = 1
	Protocol_GRPC        Protocol = 2
)

// Enum value maps for Protocol.
var (
	Protocol_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "HTTP",
		2: "GRPC",
	}
	Protocol_value = map[string]int32{
		"UNSPECIFIED": 0,
		"HTTP":        1,
		"GRPC":        2,
	}
)

func (x Protocol) Enum() *Protocol {
	p := new(Protocol)
	*p = x
	return p
}

func (x Protocol) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Protocol) Descriptor() protoreflect.EnumDescriptor {
	return file_config_v1_gateway_proto_enumTypes[0].Descriptor()
}

func (Protocol) Type() protoreflect.EnumType {
	return &file_config_v1_gateway_proto_enumTypes[0]
}

func (x Protocol) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Protocol.Descriptor instead.
func (Protocol) EnumDescriptor() ([]byte, []int) {
	return file_config_v1_gateway_proto_rawDescGZIP(), []int{0}
}

type Gateway struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	// Deprecated: Marked as deprecated in config/v1/gateway.proto.
	Hosts       []string        `protobuf:"bytes,3,rep,name=hosts,proto3" json:"hosts,omitempty"`
	Endpoints   []*Endpoint     `protobuf:"bytes,4,rep,name=endpoints,proto3" json:"endpoints,omitempty"`
	Middlewares []*Middleware   `protobuf:"bytes,5,rep,name=middlewares,proto3" json:"middlewares,omitempty"`
	TlsStore    map[string]*TLS `protobuf:"bytes,6,rep,name=tls_store,json=tlsStore,proto3" json:"tls_store,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Gateway) Reset() {
	*x = Gateway{}
	mi := &file_config_v1_gateway_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Gateway) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Gateway) ProtoMessage() {}

func (x *Gateway) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_gateway_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Gateway.ProtoReflect.Descriptor instead.
func (*Gateway) Descriptor() ([]byte, []int) {
	return file_config_v1_gateway_proto_rawDescGZIP(), []int{0}
}

func (x *Gateway) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Gateway) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

// Deprecated: Marked as deprecated in config/v1/gateway.proto.
func (x *Gateway) GetHosts() []string {
	if x != nil {
		return x.Hosts
	}
	return nil
}

func (x *Gateway) GetEndpoints() []*Endpoint {
	if x != nil {
		return x.Endpoints
	}
	return nil
}

func (x *Gateway) GetMiddlewares() []*Middleware {
	if x != nil {
		return x.Middlewares
	}
	return nil
}

func (x *Gateway) GetTlsStore() map[string]*TLS {
	if x != nil {
		return x.TlsStore
	}
	return nil
}

type TLS struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Insecure   bool   `protobuf:"varint,1,opt,name=insecure,proto3" json:"insecure,omitempty"`
	Cacert     string `protobuf:"bytes,2,opt,name=cacert,proto3" json:"cacert,omitempty"`
	Cert       string `protobuf:"bytes,3,opt,name=cert,proto3" json:"cert,omitempty"`
	Key        string `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty"`
	ServerName string `protobuf:"bytes,5,opt,name=server_name,json=serverName,proto3" json:"server_name,omitempty"`
}

func (x *TLS) Reset() {
	*x = TLS{}
	mi := &file_config_v1_gateway_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TLS) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TLS) ProtoMessage() {}

func (x *TLS) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_gateway_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TLS.ProtoReflect.Descriptor instead.
func (*TLS) Descriptor() ([]byte, []int) {
	return file_config_v1_gateway_proto_rawDescGZIP(), []int{1}
}

func (x *TLS) GetInsecure() bool {
	if x != nil {
		return x.Insecure
	}
	return false
}

func (x *TLS) GetCacert() string {
	if x != nil {
		return x.Cacert
	}
	return ""
}

func (x *TLS) GetCert() string {
	if x != nil {
		return x.Cert
	}
	return ""
}

func (x *TLS) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *TLS) GetServerName() string {
	if x != nil {
		return x.ServerName
	}
	return ""
}

type PriorityConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string      `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version   string      `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Endpoints []*Endpoint `protobuf:"bytes,3,rep,name=endpoints,proto3" json:"endpoints,omitempty"`
}

func (x *PriorityConfig) Reset() {
	*x = PriorityConfig{}
	mi := &file_config_v1_gateway_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PriorityConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PriorityConfig) ProtoMessage() {}

func (x *PriorityConfig) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_gateway_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PriorityConfig.ProtoReflect.Descriptor instead.
func (*PriorityConfig) Descriptor() ([]byte, []int) {
	return file_config_v1_gateway_proto_rawDescGZIP(), []int{2}
}

func (x *PriorityConfig) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PriorityConfig) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *PriorityConfig) GetEndpoints() []*Endpoint {
	if x != nil {
		return x.Endpoints
	}
	return nil
}

type Endpoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path        string               `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	Method      string               `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	Description string               `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Protocol    Protocol             `protobuf:"varint,4,opt,name=protocol,proto3,enum=config.v1.Protocol" json:"protocol,omitempty"`
	Timeout     *durationpb.Duration `protobuf:"bytes,5,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Middlewares []*Middleware        `protobuf:"bytes,6,rep,name=middlewares,proto3" json:"middlewares,omitempty"`
	Backends    []*Backend           `protobuf:"bytes,7,rep,name=backends,proto3" json:"backends,omitempty"`
	Retry       *Retry               `protobuf:"bytes,8,opt,name=retry,proto3" json:"retry,omitempty"`
	Metadata    map[string]string    `protobuf:"bytes,9,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Host        string               `protobuf:"bytes,10,opt,name=host,proto3" json:"host,omitempty"`
}

func (x *Endpoint) Reset() {
	*x = Endpoint{}
	mi := &file_config_v1_gateway_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Endpoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Endpoint) ProtoMessage() {}

func (x *Endpoint) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_gateway_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Endpoint.ProtoReflect.Descriptor instead.
func (*Endpoint) Descriptor() ([]byte, []int) {
	return file_config_v1_gateway_proto_rawDescGZIP(), []int{3}
}

func (x *Endpoint) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *Endpoint) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *Endpoint) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Endpoint) GetProtocol() Protocol {
	if x != nil {
		return x.Protocol
	}
	return Protocol_UNSPECIFIED
}

func (x *Endpoint) GetTimeout() *durationpb.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

func (x *Endpoint) GetMiddlewares() []*Middleware {
	if x != nil {
		return x.Middlewares
	}
	return nil
}

func (x *Endpoint) GetBackends() []*Backend {
	if x != nil {
		return x.Backends
	}
	return nil
}

func (x *Endpoint) GetRetry() *Retry {
	if x != nil {
		return x.Retry
	}
	return nil
}

func (x *Endpoint) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Endpoint) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

type Middleware struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Options  *anypb.Any `protobuf:"bytes,2,opt,name=options,proto3" json:"options,omitempty"`
	Required bool       `protobuf:"varint,3,opt,name=required,proto3" json:"required,omitempty"`
}

func (x *Middleware) Reset() {
	*x = Middleware{}
	mi := &file_config_v1_gateway_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Middleware) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Middleware) ProtoMessage() {}

func (x *Middleware) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_gateway_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Middleware.ProtoReflect.Descriptor instead.
func (*Middleware) Descriptor() ([]byte, []int) {
	return file_config_v1_gateway_proto_rawDescGZIP(), []int{4}
}

func (x *Middleware) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Middleware) GetOptions() *anypb.Any {
	if x != nil {
		return x.Options
	}
	return nil
}

func (x *Middleware) GetRequired() bool {
	if x != nil {
		return x.Required
	}
	return false
}

type Backend struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// localhost
	// 127.0.0.1:8000
	// discovery:///service_name
	Target        string            `protobuf:"bytes,1,opt,name=target,proto3" json:"target,omitempty"`
	Weight        *int64            `protobuf:"varint,2,opt,name=weight,proto3,oneof" json:"weight,omitempty"`
	HealthCheck   *HealthCheck      `protobuf:"bytes,3,opt,name=health_check,json=healthCheck,proto3" json:"health_check,omitempty"`
	Tls           bool              `protobuf:"varint,4,opt,name=tls,proto3" json:"tls,omitempty"`
	TlsConfigName string            `protobuf:"bytes,5,opt,name=tls_config_name,json=tlsConfigName,proto3" json:"tls_config_name,omitempty"`
	Metadata      map[string]string `protobuf:"bytes,6,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Backend) Reset() {
	*x = Backend{}
	mi := &file_config_v1_gateway_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Backend) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Backend) ProtoMessage() {}

func (x *Backend) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_gateway_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Backend.ProtoReflect.Descriptor instead.
func (*Backend) Descriptor() ([]byte, []int) {
	return file_config_v1_gateway_proto_rawDescGZIP(), []int{5}
}

func (x *Backend) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

func (x *Backend) GetWeight() int64 {
	if x != nil && x.Weight != nil {
		return *x.Weight
	}
	return 0
}

func (x *Backend) GetHealthCheck() *HealthCheck {
	if x != nil {
		return x.HealthCheck
	}
	return nil
}

func (x *Backend) GetTls() bool {
	if x != nil {
		return x.Tls
	}
	return false
}

func (x *Backend) GetTlsConfigName() string {
	if x != nil {
		return x.TlsConfigName
	}
	return ""
}

func (x *Backend) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type HealthCheck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *HealthCheck) Reset() {
	*x = HealthCheck{}
	mi := &file_config_v1_gateway_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HealthCheck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthCheck) ProtoMessage() {}

func (x *HealthCheck) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_gateway_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthCheck.ProtoReflect.Descriptor instead.
func (*HealthCheck) Descriptor() ([]byte, []int) {
	return file_config_v1_gateway_proto_rawDescGZIP(), []int{6}
}

type Retry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// default attempts is 1
	Attempts      uint32               `protobuf:"varint,1,opt,name=attempts,proto3" json:"attempts,omitempty"`
	PerTryTimeout *durationpb.Duration `protobuf:"bytes,2,opt,name=per_try_timeout,json=perTryTimeout,proto3" json:"per_try_timeout,omitempty"`
	Conditions    []*Condition         `protobuf:"bytes,3,rep,name=conditions,proto3" json:"conditions,omitempty"`
	// primary,secondary
	Priorities []string `protobuf:"bytes,4,rep,name=priorities,proto3" json:"priorities,omitempty"`
}

func (x *Retry) Reset() {
	*x = Retry{}
	mi := &file_config_v1_gateway_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Retry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Retry) ProtoMessage() {}

func (x *Retry) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_gateway_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Retry.ProtoReflect.Descriptor instead.
func (*Retry) Descriptor() ([]byte, []int) {
	return file_config_v1_gateway_proto_rawDescGZIP(), []int{7}
}

func (x *Retry) GetAttempts() uint32 {
	if x != nil {
		return x.Attempts
	}
	return 0
}

func (x *Retry) GetPerTryTimeout() *durationpb.Duration {
	if x != nil {
		return x.PerTryTimeout
	}
	return nil
}

func (x *Retry) GetConditions() []*Condition {
	if x != nil {
		return x.Conditions
	}
	return nil
}

func (x *Retry) GetPriorities() []string {
	if x != nil {
		return x.Priorities
	}
	return nil
}

type Condition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Condition:
	//
	//	*Condition_ByStatusCode
	//	*Condition_ByHeader
	Condition isCondition_Condition `protobuf_oneof:"condition"`
}

func (x *Condition) Reset() {
	*x = Condition{}
	mi := &file_config_v1_gateway_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Condition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Condition) ProtoMessage() {}

func (x *Condition) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_gateway_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Condition.ProtoReflect.Descriptor instead.
func (*Condition) Descriptor() ([]byte, []int) {
	return file_config_v1_gateway_proto_rawDescGZIP(), []int{8}
}

func (m *Condition) GetCondition() isCondition_Condition {
	if m != nil {
		return m.Condition
	}
	return nil
}

func (x *Condition) GetByStatusCode() string {
	if x, ok := x.GetCondition().(*Condition_ByStatusCode); ok {
		return x.ByStatusCode
	}
	return ""
}

func (x *Condition) GetByHeader() *ConditionHeader {
	if x, ok := x.GetCondition().(*Condition_ByHeader); ok {
		return x.ByHeader
	}
	return nil
}

type isCondition_Condition interface {
	isCondition_Condition()
}

type Condition_ByStatusCode struct {
	// "500-599", "429"
	ByStatusCode string `protobuf:"bytes,1,opt,name=by_status_code,json=byStatusCode,proto3,oneof"`
}

type Condition_ByHeader struct {
	// {"name": "grpc-status", "value": "14"}
	ByHeader *ConditionHeader `protobuf:"bytes,2,opt,name=by_header,json=byHeader,proto3,oneof"`
}

func (*Condition_ByStatusCode) isCondition_Condition() {}

func (*Condition_ByHeader) isCondition_Condition() {}

type ConditionHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ConditionHeader) Reset() {
	*x = ConditionHeader{}
	mi := &file_config_v1_gateway_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConditionHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConditionHeader) ProtoMessage() {}

func (x *ConditionHeader) ProtoReflect() protoreflect.Message {
	mi := &file_config_v1_gateway_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConditionHeader.ProtoReflect.Descriptor instead.
func (*ConditionHeader) Descriptor() ([]byte, []int) {
	return file_config_v1_gateway_proto_rawDescGZIP(), []int{8, 0}
}

func (x *ConditionHeader) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ConditionHeader) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_config_v1_gateway_proto protoreflect.FileDescriptor

var file_config_v1_gateway_proto_rawDesc = []byte{
	0x0a, 0x17, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x61, 0x74, 0x65,
	0x77, 0x61, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x76, 0x31, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xc9, 0x02, 0x0a, 0x07, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x05, 0x68, 0x6f, 0x73,
	0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x42, 0x02, 0x18, 0x01, 0x52, 0x05, 0x68, 0x6f,
	0x73, 0x74, 0x73, 0x12, 0x31, 0x0a, 0x09, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73,
	0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x09, 0x65, 0x6e, 0x64,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x37, 0x0a, 0x0b, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65,
	0x77, 0x61, 0x72, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61,
	0x72, 0x65, 0x52, 0x0b, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x73, 0x12,
	0x3d, 0x0a, 0x09, 0x74, 0x6c, 0x73, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x18, 0x06, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x20, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x47,
	0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x54, 0x6c, 0x73, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x74, 0x6c, 0x73, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x1a, 0x4b,
	0x0a, 0x0d, 0x54, 0x6c, 0x73, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x24, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0e, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x4c, 0x53,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x80, 0x01, 0x0a, 0x03,
	0x54, 0x4c, 0x53, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x6e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x63, 0x61, 0x63, 0x65, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x63, 0x61, 0x63, 0x65, 0x72, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x65, 0x72, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x65, 0x72, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1f, 0x0a,
	0x0b, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x71,
	0x0a, 0x0e, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x31,
	0x0a, 0x09, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e,
	0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x09, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x73, 0x22, 0xdf, 0x03, 0x0a, 0x08, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61,
	0x74, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2f, 0x0a, 0x08,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13,
	0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x33, 0x0a,
	0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x12, 0x37, 0x0a, 0x0b, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65,
	0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x52, 0x0b,
	0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x73, 0x12, 0x2e, 0x0a, 0x08, 0x62,
	0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e,
	0x64, 0x52, 0x08, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x73, 0x12, 0x26, 0x0a, 0x05, 0x72,
	0x65, 0x74, 0x72, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x74, 0x72, 0x79, 0x52, 0x05, 0x72, 0x65,
	0x74, 0x72, 0x79, 0x12, 0x3d, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76,
	0x31, 0x2e, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x22, 0x6c, 0x0a, 0x0a, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2e, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x07, 0x6f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65,
	0x64, 0x22, 0xb9, 0x02, 0x0a, 0x07, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x1b, 0x0a, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x88,
	0x01, 0x01, 0x12, 0x39, 0x0a, 0x0c, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x5f, 0x63, 0x68, 0x65,
	0x63, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x52, 0x0b, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x10, 0x0a,
	0x03, 0x74, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x74, 0x6c, 0x73, 0x12,
	0x26, 0x0a, 0x0f, 0x74, 0x6c, 0x73, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x74, 0x6c, 0x73, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x3c, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x0d, 0x0a,
	0x0b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x22, 0xbc, 0x01, 0x0a,
	0x05, 0x52, 0x65, 0x74, 0x72, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x61, 0x74, 0x74, 0x65, 0x6d, 0x70,
	0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x61, 0x74, 0x74, 0x65, 0x6d, 0x70,
	0x74, 0x73, 0x12, 0x41, 0x0a, 0x0f, 0x70, 0x65, 0x72, 0x5f, 0x74, 0x72, 0x79, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x70, 0x65, 0x72, 0x54, 0x72, 0x79, 0x54, 0x69,
	0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x34, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x0a, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x70,
	0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x0a, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x22, 0xb0, 0x01, 0x0a, 0x09,
	0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x26, 0x0a, 0x0e, 0x62, 0x79, 0x5f,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x0c, 0x62, 0x79, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64,
	0x65, 0x12, 0x3a, 0x0a, 0x09, 0x62, 0x79, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x68, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x48, 0x00, 0x52, 0x08, 0x62, 0x79, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x1a, 0x32, 0x0a,
	0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x42, 0x0b, 0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2a, 0x2f,
	0x0a, 0x08, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4e,
	0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x48,
	0x54, 0x54, 0x50, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x47, 0x52, 0x50, 0x43, 0x10, 0x02, 0x42,
	0x9d, 0x01, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76,
	0x31, 0x42, 0x0c, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x72,
	0x69, 0x67, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2f,
	0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31,
	0x3b, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x76, 0x31, 0xf8, 0x01, 0x01, 0xa2, 0x02, 0x03, 0x43,
	0x58, 0x58, 0xaa, 0x02, 0x09, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x56, 0x31, 0xca, 0x02,
	0x09, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x15, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0xea, 0x02, 0x0a, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x3a, 0x3a, 0x56, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_config_v1_gateway_proto_rawDescOnce sync.Once
	file_config_v1_gateway_proto_rawDescData = file_config_v1_gateway_proto_rawDesc
)

func file_config_v1_gateway_proto_rawDescGZIP() []byte {
	file_config_v1_gateway_proto_rawDescOnce.Do(func() {
		file_config_v1_gateway_proto_rawDescData = protoimpl.X.CompressGZIP(file_config_v1_gateway_proto_rawDescData)
	})
	return file_config_v1_gateway_proto_rawDescData
}

var file_config_v1_gateway_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_config_v1_gateway_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_config_v1_gateway_proto_goTypes = []any{
	(Protocol)(0),               // 0: config.v1.Protocol
	(*Gateway)(nil),             // 1: config.v1.Gateway
	(*TLS)(nil),                 // 2: config.v1.TLS
	(*PriorityConfig)(nil),      // 3: config.v1.PriorityConfig
	(*Endpoint)(nil),            // 4: config.v1.Endpoint
	(*Middleware)(nil),          // 5: config.v1.Middleware
	(*Backend)(nil),             // 6: config.v1.Backend
	(*HealthCheck)(nil),         // 7: config.v1.HealthCheck
	(*Retry)(nil),               // 8: config.v1.Retry
	(*Condition)(nil),           // 9: config.v1.Condition
	nil,                         // 10: config.v1.Gateway.TlsStoreEntry
	nil,                         // 11: config.v1.Endpoint.MetadataEntry
	nil,                         // 12: config.v1.Backend.MetadataEntry
	(*ConditionHeader)(nil),     // 13: config.v1.Condition.header
	(*durationpb.Duration)(nil), // 14: google.protobuf.Duration
	(*anypb.Any)(nil),           // 15: google.protobuf.Any
}
var file_config_v1_gateway_proto_depIdxs = []int32{
	4,  // 0: config.v1.Gateway.endpoints:type_name -> config.v1.Endpoint
	5,  // 1: config.v1.Gateway.middlewares:type_name -> config.v1.Middleware
	10, // 2: config.v1.Gateway.tls_store:type_name -> config.v1.Gateway.TlsStoreEntry
	4,  // 3: config.v1.PriorityConfig.endpoints:type_name -> config.v1.Endpoint
	0,  // 4: config.v1.Endpoint.protocol:type_name -> config.v1.Protocol
	14, // 5: config.v1.Endpoint.timeout:type_name -> google.protobuf.Duration
	5,  // 6: config.v1.Endpoint.middlewares:type_name -> config.v1.Middleware
	6,  // 7: config.v1.Endpoint.backends:type_name -> config.v1.Backend
	8,  // 8: config.v1.Endpoint.retry:type_name -> config.v1.Retry
	11, // 9: config.v1.Endpoint.metadata:type_name -> config.v1.Endpoint.MetadataEntry
	15, // 10: config.v1.Middleware.options:type_name -> google.protobuf.Any
	7,  // 11: config.v1.Backend.health_check:type_name -> config.v1.HealthCheck
	12, // 12: config.v1.Backend.metadata:type_name -> config.v1.Backend.MetadataEntry
	14, // 13: config.v1.Retry.per_try_timeout:type_name -> google.protobuf.Duration
	9,  // 14: config.v1.Retry.conditions:type_name -> config.v1.Condition
	13, // 15: config.v1.Condition.by_header:type_name -> config.v1.Condition.header
	2,  // 16: config.v1.Gateway.TlsStoreEntry.value:type_name -> config.v1.TLS
	17, // [17:17] is the sub-list for method output_type
	17, // [17:17] is the sub-list for method input_type
	17, // [17:17] is the sub-list for extension type_name
	17, // [17:17] is the sub-list for extension extendee
	0,  // [0:17] is the sub-list for field type_name
}

func init() { file_config_v1_gateway_proto_init() }
func file_config_v1_gateway_proto_init() {
	if File_config_v1_gateway_proto != nil {
		return
	}
	file_config_v1_gateway_proto_msgTypes[5].OneofWrappers = []any{}
	file_config_v1_gateway_proto_msgTypes[8].OneofWrappers = []any{
		(*Condition_ByStatusCode)(nil),
		(*Condition_ByHeader)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_config_v1_gateway_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_v1_gateway_proto_goTypes,
		DependencyIndexes: file_config_v1_gateway_proto_depIdxs,
		EnumInfos:         file_config_v1_gateway_proto_enumTypes,
		MessageInfos:      file_config_v1_gateway_proto_msgTypes,
	}.Build()
	File_config_v1_gateway_proto = out.File
	file_config_v1_gateway_proto_rawDesc = nil
	file_config_v1_gateway_proto_goTypes = nil
	file_config_v1_gateway_proto_depIdxs = nil
}