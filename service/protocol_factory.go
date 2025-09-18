package service

import (
	"context"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
)

// ProtocolFactory 定义了创建特定协议服务实例的工厂标准。
type ProtocolFactory interface {
	NewServer(cfg *transportv1.Transport, opts ...Option) (interfaces.Server, error)
	NewClient(ctx context.Context, cfg *transportv1.Transport, opts ...Option) (interfaces.Client, error)
}
