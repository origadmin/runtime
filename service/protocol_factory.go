package service

import (
	"github.com/origadmin/framework/runtime/interfaces"
	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
)

// ProtocolFactory 定义了创建特定协议服务实例的工厂标准。
type ProtocolFactory interface {
	NewServer(cfg *configv1.Service, opts ...Option) (interfaces.Server, error)
	// NewClient(...) 待定
}
