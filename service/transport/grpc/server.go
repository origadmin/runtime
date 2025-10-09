package grpc

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	servicetls "github.com/origadmin/runtime/service/tls"
)

// NewGRPCServer creates a new concrete gRPC server instance based on the provided configuration.
// It returns *transgrpc.Server, not the generic interfaces.Server.
func NewGRPCServer(grpcConfig *transportv1.GrpcServerConfig, serverOpts *ServerOptions) (*transgrpc.Server, error) {
	// Prepare the Kratos gRPC server options.
	var kratosOpts []transgrpc.ServerOption

	// Get Container instance if available.
	var c interfaces.Container
	if serverOpts.ServiceOptions != nil {
		c = serverOpts.ServiceOptions.Container
	}

	// Return error if middlewares are configured but no container is provided.
	if len(grpcConfig.GetMiddlewares()) > 0 && c == nil {
		return nil, fmt.Errorf("application container is required for middleware but not found in options")
	}

	// Apply options from the protobuf configuration.
	if grpcConfig.GetAddr() != "" {
		kratosOpts = append(kratosOpts, transgrpc.Address(grpcConfig.GetAddr()))
	}
	// Configure TLS for server
	if tlsConfig := grpcConfig.GetTlsConfig(); tlsConfig != nil && tlsConfig.GetEnabled() {
		tlsCfg, err := servicetls.NewServerTLSConfig(tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create server TLS config: %w", err)
		}
		kratosOpts = append(kratosOpts, transgrpc.TLSConfig(tlsCfg))
	}

	if grpcConfig.GetTimeout() != nil {
		kratosOpts = append(kratosOpts, transgrpc.Timeout(grpcConfig.GetTimeout().AsDuration()))
	}

	// Configure middlewares.
	var mws []middleware.Middleware
	if len(grpcConfig.GetMiddlewares()) > 0 {
		// 'c' is guaranteed to be non-nil at this point.
		for _, name := range grpcConfig.GetMiddlewares() {
			m, ok := c.ServerMiddleware(name)
			if !ok {
				return nil, fmt.Errorf("middleware '%s' not found in container", name)
			}
			mws = append(mws, m)
		}
	} else {
		// If no specific middlewares are configured, use default ones from adapter.go.
		mws = DefaultServerMiddlewares()
	}

	if len(mws) > 0 {
		kratosOpts = append(kratosOpts, transgrpc.Middleware(mws...))
	}

	// Apply any external Kratos gRPC server options passed via functional options.
	// These are applied last, allowing them to override previous options if needed.
	if len(serverOpts.GrpcServerOptions) > 0 {
		kratosOpts = append(kratosOpts, serverOpts.GrpcServerOptions...)
	}

	// Create the Kratos gRPC server instance.
	srv := transgrpc.NewServer(kratosOpts...)

	return srv, nil
}
