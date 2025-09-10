package grpc

import (
	"time"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/context" // Use project's context

	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/service"
)

const (
	defaultTimeout = 5 * time.Second
)

// grpcProtocolFactory implements service.ProtocolFactory for gRPC.
type grpcProtocolFactory struct{}

// NewClient creates a new gRPC client instance by delegating to the direct implementation.
func (f *grpcProtocolFactory) NewClient(ctx context.Context, cfg *configv1.Service, opts ...service.Option) (interfaces.Client, error) {
	return NewClient(ctx, cfg, opts...)
}

// NewServer creates a new gRPC server instance by delegating to the direct implementation.
func (f *grpcProtocolFactory) NewServer(cfg *configv1.Service, opts ...service.Option) (interfaces.Server, error) {
	return NewServer(cfg, opts...)
}

func init() {
	service.RegisterProtocol("grpc", &grpcProtocolFactory{})
}
