package grpc

import (
	"context"
	"fmt"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/service"
)

// grpcProtocolFactory implements the service.ProtocolFactory for the gRPC protocol.
type grpcProtocolFactory struct{}

// init registers this factory with the framework's central protocol registry.
func init() {
	service.RegisterProtocol("grpc", &grpcProtocolFactory{})
}

// NewServer creates a new gRPC server instance based on the provided configuration.
func (f *grpcProtocolFactory) NewServer(cfg *transportv1.Server, opts ...options.Option) (interfaces.Server, error) {
	// 1. Extract the gRPC-specific configuration from the protobuf config.
	grpcConfig := cfg.GetGrpc()
	if grpcConfig == nil {
		return nil, fmt.Errorf("gRPC server config is missing in transport container")
	}

	// 2. Get all gRPC server-specific and common service-level options.
	// This uses the FromServerOptions pattern to correctly apply functional options.
	serverOpts := FromServerOptions(opts)

	// Call the concrete server creation function.
	srv, err := NewGRPCServer(grpcConfig, serverOpts)
	if err != nil {
		return nil, err
	}

	// 9. Register the user's business logic services if a registrar is provided.
	if serverOpts.ServiceOptions != nil && serverOpts.ServiceOptions.Registrar != nil {
		if grpcRegistrar, ok := serverOpts.ServiceOptions.Registrar.(service.GRPCRegistrar); ok {
			grpcRegistrar.RegisterGRPC(srv)
		} else {
			return nil, fmt.Errorf("invalid registrar: expected service.GRPCRegistrar, got %T", serverOpts.ServiceOptions.Registrar)
		}
	}

	return srv, nil
}

// NewClient creates a new gRPC client instance based on the provided configuration.
func (f *grpcProtocolFactory) NewClient(ctx context.Context, cfg *transportv1.Client, opts ...options.Option) (interfaces.Client, error) {
	// 1. Extract the gRPC-specific configuration from the protobuf config.
	grpcConfig := cfg.GetGrpc()
	if grpcConfig == nil {
		return nil, fmt.Errorf("gRPC client config is missing in transport container")
	}

	// 2. Get all gRPC client-specific and common service-level options.
	// This uses the FromClientOptions pattern to correctly apply functional options.
	clientOpts := FromClientOptions(opts)

	// Call the concrete client creation function.
	conn, err := NewGRPCClient(ctx, grpcConfig, clientOpts)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
