package service

import (
	"context"

	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
)

// ProtocolFactory defines the factory standard for creating a specific protocol service instance。
type ProtocolFactory interface {
	NewServer(cfg *transportv1.Server, opts ...options.Option) (interfaces.Server, error)
	NewClient(ctx context.Context, cfg *transportv1.Client, opts ...options.Option) (interfaces.Client, error)
}
