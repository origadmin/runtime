package grpc

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/goexts/generic/configure"

	"github.com/origadmin/framework/runtime/interfaces"
	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service"
	"github.com/origadmin/runtime/service/tls"
)

// grpcProtocolFactory implements service.ProtocolFactory for gRPC.
type grpcProtocolFactory struct{}

const defaultTimeout = 5 * time.Second

// NewClient creates a new gRPC client instance.
func (f *grpcProtocolFactory) NewClient(ctx context.Context, cfg *configv1.Service, opts ...service.Option) (interfaces.Client, error) {
	ll := log.NewHelper(log.With(log.GetLogger(), "module", "service/grpc"))
	ll.Debugf("Creating new gRPC client with config: %+v", cfg)

	// Create a new service.Options instance and apply the incoming options.
	svcOpts := &service.Options{ContextOptions: interfaces.ContextOptions{Context: ctx}}
	configure.Apply(svcOpts, opts)

	// Initialize client options
	var clientOpts []grpc.ClientOption
	timeout := defaultTimeout

	// Apply configuration from configv1.Service
	if cfg.GetGrpc() != nil {
		grpcCfg := cfg.GetGrpc()
		if grpcCfg.GetTimeout() != 0 {
			timeout = time.Duration(grpcCfg.GetTimeout() * 1e6) // Convert milliseconds to nanoseconds
		}

		// Configure TLS if needed
		if grpcCfg.GetUseTls() {
			tlsConfig, err := tls.NewClientTLSConfig(grpcCfg.GetTlsConfig())
			if err != nil {
				return nil, err
			}
			if tlsConfig != nil {
				clientOpts = append(clientOpts, grpc.WithTLSConfig(tlsConfig))
			}
		}

		// Set endpoint if provided
		if grpcCfg.GetEndpoint() != "" {
			clientOpts = append(clientOpts, grpc.WithEndpoint(grpcCfg.GetEndpoint()))
		}
	}

	// Apply timeout
	clientOpts = append(clientOpts, grpc.WithTimeout(timeout))

	// Extract gRPC client specific options from the service.Options' Context.
	// These are the options added via grpc.WithClientOption
	clientOptsFromContext := FromClientOptions(svcOpts)
	clientOpts = append(clientOpts, clientOptsFromContext...)

	// Create the client with the merged options
	client, err := grpc.DialInsecure(ctx, clientOpts...)
	if err != nil {
		return nil, errors.Newf(500, "INTERNAL_SERVER_ERROR", "create grpc client failed: %v", err)
	}

	return client, nil
}

// NewServer creates a new gRPC server instance.
func (f *grpcProtocolFactory) NewServer(cfg *configv1.Service, opts ...service.Option) (interfaces.Server, error) {
	ll := log.NewHelper(log.With(log.GetLogger(), "module", "service/grpc"))
	ll.Debugf("Creating new gRPC server instance with config: %+v", cfg)

	// Create a new service.Options instance and apply the incoming options.
	svcOpts := &service.Options{ContextOptions: interfaces.ContextOptions{Context: context.Background()}}
	configure.Apply(svcOpts, opts)

	// Initialize Kratos gRPC server options
	var serverOpts []grpc.ServerOption

	// Add default recovery middleware and error encoder
	serverOpts = append(serverOpts, grpc.Middleware(recovery.Recovery()))

	// Apply configuration from configv1.Service
	if cfg.GetGrpc() != nil {
		grpcCfg := cfg.GetGrpc()

		if grpcCfg.GetUseTls() {
			tlsConfig, err := tls.NewServerTLSConfig(grpcCfg.GetTlsConfig())
			if err != nil {
				return nil, err
			}
			if tlsConfig != nil {
				serverOpts = append(serverOpts, grpc.TLSConfig(tlsConfig))
			}
		}

		if grpcCfg.GetNetwork() != "" {
			serverOpts = append(serverOpts, grpc.Network(grpcCfg.GetNetwork()))
		}

		if grpcCfg.GetAddr() != "" {
			serverOpts = append(serverOpts, grpc.Address(grpcCfg.GetAddr()))
		}

		timeout := defaultTimeout
		if grpcCfg.GetTimeout() != 0 {
			timeout = time.Duration(grpcCfg.GetTimeout() * 1e6) // Convert milliseconds to nanoseconds
		}
		serverOpts = append(serverOpts, grpc.Timeout(timeout))

		ll.Debugw("gRPC server configured", "endpoint", grpcCfg.GetEndpoint())
	}

	// Extract gRPC specific options from the service.Options' Context.
	serverOptsFromContext := FromServerOptions(svcOpts)

	// Combine all gRPC options
	grpcOpts := append(serverOpts, serverOptsFromContext...)

	return grpc.NewServer(grpcOpts...), nil
}

// init registers the gRPC protocol factory with the global service registry.
func init() {
	// Register the gRPC protocol factory with the service module.
	service.RegisterProtocol("grpc", &grpcProtocolFactory{})
}
