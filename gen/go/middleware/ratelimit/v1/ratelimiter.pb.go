// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: middleware/ratelimit/v1/ratelimiter.proto

package ratelimitv1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
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

// Rate limiter
type RateLimiter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Enabled bool `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	// rate limiter name, supported: bbr, memory, redis.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The number of seconds in a rate limit window
	Period int32 `protobuf:"varint,3,opt,name=period,proto3" json:"period,omitempty"`
	// The number of requests allowed in a window of time
	XRatelimitLimit int32 `protobuf:"varint,5,opt,name=x_ratelimit_limit,proto3" json:"x_ratelimit_limit,omitempty"`
	// The number of requests that can still be made in the current window of time
	XRatelimitRemaining int32 `protobuf:"varint,6,opt,name=x_ratelimit_remaining,proto3" json:"x_ratelimit_remaining,omitempty"`
	// The number of seconds until the current rate limit window completely resets
	XRatelimitReset int32 `protobuf:"varint,7,opt,name=x_ratelimit_reset,proto3" json:"x_ratelimit_reset,omitempty"`
	// When rate limited, the number of seconds to wait before another request will be accepted
	RetryAfter int32               `protobuf:"varint,8,opt,name=retry_after,proto3" json:"retry_after,omitempty"`
	Memory     *RateLimiter_Memory `protobuf:"bytes,101,opt,name=memory,proto3" json:"memory,omitempty"`
	Redis      *RateLimiter_Redis  `protobuf:"bytes,102,opt,name=redis,proto3" json:"redis,omitempty"`
}

func (x *RateLimiter) Reset() {
	*x = RateLimiter{}
	mi := &file_middleware_ratelimit_v1_ratelimiter_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RateLimiter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimiter) ProtoMessage() {}

func (x *RateLimiter) ProtoReflect() protoreflect.Message {
	mi := &file_middleware_ratelimit_v1_ratelimiter_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateLimiter.ProtoReflect.Descriptor instead.
func (*RateLimiter) Descriptor() ([]byte, []int) {
	return file_middleware_ratelimit_v1_ratelimiter_proto_rawDescGZIP(), []int{0}
}

func (x *RateLimiter) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *RateLimiter) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RateLimiter) GetPeriod() int32 {
	if x != nil {
		return x.Period
	}
	return 0
}

func (x *RateLimiter) GetXRatelimitLimit() int32 {
	if x != nil {
		return x.XRatelimitLimit
	}
	return 0
}

func (x *RateLimiter) GetXRatelimitRemaining() int32 {
	if x != nil {
		return x.XRatelimitRemaining
	}
	return 0
}

func (x *RateLimiter) GetXRatelimitReset() int32 {
	if x != nil {
		return x.XRatelimitReset
	}
	return 0
}

func (x *RateLimiter) GetRetryAfter() int32 {
	if x != nil {
		return x.RetryAfter
	}
	return 0
}

func (x *RateLimiter) GetMemory() *RateLimiter_Memory {
	if x != nil {
		return x.Memory
	}
	return nil
}

func (x *RateLimiter) GetRedis() *RateLimiter_Redis {
	if x != nil {
		return x.Redis
	}
	return nil
}

type RateLimiter_Redis struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr     string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Username string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	Db       int32  `protobuf:"varint,4,opt,name=db,proto3" json:"db,omitempty"`
}

func (x *RateLimiter_Redis) Reset() {
	*x = RateLimiter_Redis{}
	mi := &file_middleware_ratelimit_v1_ratelimiter_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RateLimiter_Redis) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimiter_Redis) ProtoMessage() {}

func (x *RateLimiter_Redis) ProtoReflect() protoreflect.Message {
	mi := &file_middleware_ratelimit_v1_ratelimiter_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateLimiter_Redis.ProtoReflect.Descriptor instead.
func (*RateLimiter_Redis) Descriptor() ([]byte, []int) {
	return file_middleware_ratelimit_v1_ratelimiter_proto_rawDescGZIP(), []int{0, 0}
}

func (x *RateLimiter_Redis) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *RateLimiter_Redis) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *RateLimiter_Redis) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *RateLimiter_Redis) GetDb() int32 {
	if x != nil {
		return x.Db
	}
	return 0
}

type RateLimiter_Memory struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Expiration      *durationpb.Duration `protobuf:"bytes,1,opt,name=expiration,proto3" json:"expiration,omitempty"`
	CleanupInterval *durationpb.Duration `protobuf:"bytes,2,opt,name=cleanup_interval,proto3" json:"cleanup_interval,omitempty"`
}

func (x *RateLimiter_Memory) Reset() {
	*x = RateLimiter_Memory{}
	mi := &file_middleware_ratelimit_v1_ratelimiter_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RateLimiter_Memory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimiter_Memory) ProtoMessage() {}

func (x *RateLimiter_Memory) ProtoReflect() protoreflect.Message {
	mi := &file_middleware_ratelimit_v1_ratelimiter_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateLimiter_Memory.ProtoReflect.Descriptor instead.
func (*RateLimiter_Memory) Descriptor() ([]byte, []int) {
	return file_middleware_ratelimit_v1_ratelimiter_proto_rawDescGZIP(), []int{0, 1}
}

func (x *RateLimiter_Memory) GetExpiration() *durationpb.Duration {
	if x != nil {
		return x.Expiration
	}
	return nil
}

func (x *RateLimiter_Memory) GetCleanupInterval() *durationpb.Duration {
	if x != nil {
		return x.CleanupInterval
	}
	return nil
}

var File_middleware_ratelimit_v1_ratelimiter_proto protoreflect.FileDescriptor

var file_middleware_ratelimit_v1_ratelimiter_proto_rawDesc = []byte{
	0x0a, 0x29, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2f, 0x72, 0x61, 0x74,
	0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x61, 0x74, 0x65, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x6d, 0x69, 0x64,
	0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2e, 0x72, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9b,
	0x05, 0x0a, 0x0b, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x12, 0x18,
	0x0a, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x2d, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x19, 0xfa, 0x42, 0x16, 0x72, 0x14, 0x52, 0x03, 0x62,
	0x62, 0x72, 0x52, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x52, 0x05, 0x72, 0x65, 0x64, 0x69,
	0x73, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x65, 0x72, 0x69, 0x6f,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x12,
	0x2c, 0x0a, 0x11, 0x78, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x5f, 0x6c,
	0x69, 0x6d, 0x69, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x78, 0x5f, 0x72, 0x61,
	0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x34, 0x0a,
	0x15, 0x78, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x5f, 0x72, 0x65, 0x6d,
	0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x15, 0x78, 0x5f,
	0x72, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x5f, 0x72, 0x65, 0x6d, 0x61, 0x69, 0x6e,
	0x69, 0x6e, 0x67, 0x12, 0x2c, 0x0a, 0x11, 0x78, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d,
	0x69, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x65, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11,
	0x78, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x65,
	0x74, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x65, 0x74, 0x72, 0x79, 0x5f, 0x61, 0x66, 0x74, 0x65, 0x72,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x72, 0x65, 0x74, 0x72, 0x79, 0x5f, 0x61, 0x66,
	0x74, 0x65, 0x72, 0x12, 0x43, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x18, 0x65, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65,
	0x2e, 0x72, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61,
	0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79,
	0x52, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x12, 0x40, 0x0a, 0x05, 0x72, 0x65, 0x64, 0x69,
	0x73, 0x18, 0x66, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65,
	0x77, 0x61, 0x72, 0x65, 0x2e, 0x72, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x2e, 0x52, 0x65,
	0x64, 0x69, 0x73, 0x52, 0x05, 0x72, 0x65, 0x64, 0x69, 0x73, 0x1a, 0x63, 0x0a, 0x05, 0x52, 0x65,
	0x64, 0x69, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12,
	0x0e, 0x0a, 0x02, 0x64, 0x62, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x64, 0x62, 0x1a,
	0x8a, 0x01, 0x0a, 0x06, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x12, 0x39, 0x0a, 0x0a, 0x65, 0x78,
	0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x45, 0x0a, 0x10, 0x63, 0x6c, 0x65, 0x61, 0x6e, 0x75, 0x70,
	0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x10, 0x63, 0x6c, 0x65, 0x61,
	0x6e, 0x75, 0x70, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x42, 0xf9, 0x01, 0x0a,
	0x1b, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2e,
	0x72, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x42, 0x10, 0x52, 0x61,
	0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x72, 0x69,
	0x67, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2f, 0x67,
	0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65,
	0x2f, 0x72, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x2f, 0x76, 0x31, 0x3b, 0x72, 0x61,
	0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x76, 0x31, 0xf8, 0x01, 0x01, 0xa2, 0x02, 0x03, 0x4d,
	0x52, 0x58, 0xaa, 0x02, 0x17, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2e,
	0x52, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x17, 0x4d,
	0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x5c, 0x52, 0x61, 0x74, 0x65, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x23, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77,
	0x61, 0x72, 0x65, 0x5c, 0x52, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x5c, 0x56, 0x31,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x19, 0x4d,
	0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x3a, 0x3a, 0x52, 0x61, 0x74, 0x65, 0x6c,
	0x69, 0x6d, 0x69, 0x74, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_middleware_ratelimit_v1_ratelimiter_proto_rawDescOnce sync.Once
	file_middleware_ratelimit_v1_ratelimiter_proto_rawDescData = file_middleware_ratelimit_v1_ratelimiter_proto_rawDesc
)

func file_middleware_ratelimit_v1_ratelimiter_proto_rawDescGZIP() []byte {
	file_middleware_ratelimit_v1_ratelimiter_proto_rawDescOnce.Do(func() {
		file_middleware_ratelimit_v1_ratelimiter_proto_rawDescData = protoimpl.X.CompressGZIP(file_middleware_ratelimit_v1_ratelimiter_proto_rawDescData)
	})
	return file_middleware_ratelimit_v1_ratelimiter_proto_rawDescData
}

var file_middleware_ratelimit_v1_ratelimiter_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_middleware_ratelimit_v1_ratelimiter_proto_goTypes = []any{
	(*RateLimiter)(nil),         // 0: middleware.ratelimit.v1.RateLimiter
	(*RateLimiter_Redis)(nil),   // 1: middleware.ratelimit.v1.RateLimiter.Redis
	(*RateLimiter_Memory)(nil),  // 2: middleware.ratelimit.v1.RateLimiter.Memory
	(*durationpb.Duration)(nil), // 3: google.protobuf.Duration
}
var file_middleware_ratelimit_v1_ratelimiter_proto_depIdxs = []int32{
	2, // 0: middleware.ratelimit.v1.RateLimiter.memory:type_name -> middleware.ratelimit.v1.RateLimiter.Memory
	1, // 1: middleware.ratelimit.v1.RateLimiter.redis:type_name -> middleware.ratelimit.v1.RateLimiter.Redis
	3, // 2: middleware.ratelimit.v1.RateLimiter.Memory.expiration:type_name -> google.protobuf.Duration
	3, // 3: middleware.ratelimit.v1.RateLimiter.Memory.cleanup_interval:type_name -> google.protobuf.Duration
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_middleware_ratelimit_v1_ratelimiter_proto_init() }
func file_middleware_ratelimit_v1_ratelimiter_proto_init() {
	if File_middleware_ratelimit_v1_ratelimiter_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_middleware_ratelimit_v1_ratelimiter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_middleware_ratelimit_v1_ratelimiter_proto_goTypes,
		DependencyIndexes: file_middleware_ratelimit_v1_ratelimiter_proto_depIdxs,
		MessageInfos:      file_middleware_ratelimit_v1_ratelimiter_proto_msgTypes,
	}.Build()
	File_middleware_ratelimit_v1_ratelimiter_proto = out.File
	file_middleware_ratelimit_v1_ratelimiter_proto_rawDesc = nil
	file_middleware_ratelimit_v1_ratelimiter_proto_goTypes = nil
	file_middleware_ratelimit_v1_ratelimiter_proto_depIdxs = nil
}
