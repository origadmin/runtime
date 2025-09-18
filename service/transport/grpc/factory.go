package grpc

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	kgprc "github.com/go-kratos/kratos/v2/transport/grpc"
	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/service"
	mw "github.com/origadmin/runtime/middleware"
)

// grpcProtocolFactory implements the service.ProtocolFactory for gRPC.
type grpcProtocolFactory struct{}

// init registers this factory with the framework's protocol registry.
func init() {
	service.RegisterProtocol("grpc", &grpcProtocolFactory{})
}

// NewServer creates a new gRPC server instance.
// It conforms to the updated ProtocolFactory interface.
func (f *grpcProtocolFactory) NewServer(cfg *transportv1.Transport, opts ...service.Option) (interfaces.Server, error) {
	// 1. Extract the specific gRPC config from the container.
	grpcConfig := cfg.GetGrpc()
	if grpcConfig == nil {
		return nil, fmt.Errorf("gRPC config is missing in transport container")
	}

	// 2. Process options to extract registrar.
	var sOpts service.Options
	sOpts.Apply(opts...)

	grpcRegistrar, ok := sOpts.Value().registrar.(service.GRPCRegistrar)
	if !ok && sOpts.Value().registrar != nil {
		return nil, fmt.Errorf("invalid registrar: expected service.GRPCRegistrar, got %T", sOpts.Value().registrar)
	}

	// --- All creation logic below uses the extracted, concrete 'grpcConfig' ---

	var kOpts []kgprc.ServerOption
	var mws []middleware.Middleware

	// Build middleware chain
	for _, name := range grpcConfig.Middlewares {
		m, ok := mw.Get(name)
		if !ok {
			return nil, fmt.Errorf("middleware '%s' not found in registry", name)
		}
		mws = append(mws, m)
	}
	if len(mws) > 0 {
		kOpts = append(kOpts, kgprc.Middleware(mws...))
	}

	// Apply other server options
	if grpcConfig.Network != "" {
		kOpts = append(kOpts, kgprc.Network(grpcConfig.Network))
	}
	if grpcConfig.Addr != "" {
		kOpts = append(kOpts, kgprc.Address(grpcConfig.Addr))
	}
	if grpcConfig.Timeout != nil {
		kOpts = append(kOpts, kgprc.Timeout(grpcConfig.Timeout.AsDuration()))
	}

	// Create the gRPC server instance
	srv := kgprc.NewServer(kOpts...)

	// Register business logic
	if grpcRegistrar != nil {
		grpcRegistrar.RegisterGRPC(context.Background(), srv)
	}

	return srv, nil
}

// NewClient creates a new gRPC client instance.
// This is a placeholder implementation for now.
func (f *grpcProtocolFactory) NewClient(ctx context.Context, cfg *transportv1.Transport, opts ...service.Option) (interfaces.Client, error) {
	// TODO: Implement gRPC client creation logic
	return nil, fmt.Errorf("gRPC client creation not yet implemented")
}
