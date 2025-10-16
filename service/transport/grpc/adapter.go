package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	commonv1 "github.com/origadmin/runtime/api/gen/go/runtime/common/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces"
	serviceselector "github.com/origadmin/runtime/service/selector"
	servicetls "github.com/origadmin/runtime/service/tls"
)

const Module = "grpc.adapter"

// DefaultServerMiddlewares provides a default set of server-side middlewares for gRPC services.
// These are essential for ensuring basic stability and observability.
func DefaultServerMiddlewares() []middleware.Middleware {
	return []middleware.Middleware{
		// recovery middleware recovers from panics and converts them into errors.
		recovery.Recovery(),
	}
}

// DefaultClientMiddlewares provides a default set of client-side middlewares for gRPC services.
func DefaultClientMiddlewares() []middleware.Middleware {
	return []middleware.Middleware{
		recovery.Recovery(),
	}
}

// initGrpcServerOptions initialize the grpc server option
func initGrpcServerOptions(grpcConfig *transportv1.GrpcServerConfig, serverOpts *ServerOptions) ([]transgrpc.ServerOption, error) {
	// Prepare the Kratos gRPC server options.
	var kratosOpts []transgrpc.ServerOption

	// Get the container instance. It will be nil if not provided in options.
	var c interfaces.Container
	if serverOpts.ServiceOptions != nil {
		c = serverOpts.ServiceOptions.Container
	}

	// Check if middlewares are configured.
	hasMiddlewaresConfigured := len(grpcConfig.GetMiddlewares()) > 0

	// If middlewares are configured but no container is provided, return an error.
	// This consolidates the nil check for the container.
	if hasMiddlewaresConfigured && c == nil {
		return nil, runtimeerrors.NewStructured(Module, "application container is required for server middlewares but not found in options").WithReason(commonv1.ErrorReason_VALIDATION_ERROR).WithCaller()
	}

	// Apply options from the protobuf configuration.
	if grpcConfig.GetAddr() != "" {
		kratosOpts = append(kratosOpts, transgrpc.Address(grpcConfig.GetAddr()))
	}
	// Configure TLS for server
	if grpcConfig.GetTlsConfig().GetEnabled() {
		tlsCfg, err := servicetls.NewServerTLSConfig(grpcConfig.GetTlsConfig())
		if err != nil {
			return nil, runtimeerrors.WrapStructured(err, Module, "failed to create server TLS config").WithReason(commonv1.ErrorReason_INTERNAL_SERVER_ERROR).WithCaller()
		}
		kratosOpts = append(kratosOpts, transgrpc.TLSConfig(tlsCfg))
	}

	if grpcConfig.GetTimeout() != nil {
		kratosOpts = append(kratosOpts, transgrpc.Timeout(grpcConfig.GetTimeout().AsDuration()))
	}

	// Configure middlewares.
	var mws []middleware.Middleware
	if hasMiddlewaresConfigured {
		// 'c' is guaranteed to be non-nil at this point due to the early check above.
		for _, name := range grpcConfig.GetMiddlewares() {
			m, ok := c.ServerMiddleware(name)
			if !ok {
				return nil, runtimeerrors.NewStructured(Module, "server middleware '%s' not found in container", name).WithReason(commonv1.ErrorReason_NOT_FOUND).WithCaller()
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

	return kratosOpts, nil
}

// initGrpcClientOptions initialize grpc client options
func initGrpcClientOptions(ctx context.Context, grpcConfig *transportv1.GrpcClientConfig, clientOpts *ClientOptions) ([]transgrpc.ClientOption, error) {
	// Prepare the Kratos gRPC client options.
	var kratosOpts []transgrpc.ClientOption

	// Get the container instance. It will be nil if not provided in options.
	var c interfaces.Container
	if clientOpts.ServiceOptions != nil {
		c = clientOpts.ServiceOptions.Container
	}

	// Apply options from the protobuf configuration.
	if grpcConfig.GetTimeout() != nil {
		kratosOpts = append(kratosOpts, transgrpc.WithTimeout(grpcConfig.GetTimeout().AsDuration()))
	}

	// Configure middlewares.
	var mws []middleware.Middleware
	if len(grpcConfig.GetMiddlewares()) > 0 {
		// If middlewares are configured but no container is provided, return an error.
		if c == nil {
			return nil, runtimeerrors.NewStructured(Module, "application container is required for client middlewares but not found in options").WithReason(commonv1.ErrorReason_VALIDATION_ERROR).WithCaller()
		}
		for _, name := range grpcConfig.GetMiddlewares() {
			m, ok := c.ClientMiddleware(name)
			if !ok {
				return nil, runtimeerrors.NewStructured(Module, "client middleware '%s' not found in container", name).WithReason(commonv1.ErrorReason_NOT_FOUND).WithCaller()
			}
			mws = append(mws, m)
		}
	} else {
		// If no specific middlewares are configured, use default ones from adapter.go.
		mws = DefaultClientMiddlewares()
	}

	if len(mws) > 0 {
		kratosOpts = append(kratosOpts, transgrpc.WithMiddleware(mws...))
	}

	// Configure service discovery and endpoint.
	var discoveryClient registry.Discovery
	endpoint := grpcConfig.GetEndpoint()

	// Always apply the endpoint option.
	if endpoint != "" {
		kratosOpts = append(kratosOpts, transgrpc.WithEndpoint(endpoint))
	}

	// Determine the discovery client.
	if discoveryName := grpcConfig.GetDiscoveryName(); discoveryName != "" {
		// If a named discovery client is configured but no container is provided, return an error.
		if c == nil {
			return nil, runtimeerrors.NewStructured(Module, "application container is required for named discovery client but not found in options").WithReason(commonv1.ErrorReason_VALIDATION_ERROR).WithCaller()
		}
		if d, ok := c.Discovery(discoveryName); ok {
			discoveryClient = d
		} else {
			return nil, runtimeerrors.NewStructured(Module, "discovery client '%s' not found in container", discoveryName).WithReason(commonv1.ErrorReason_NOT_FOUND).WithCaller()
		}
	} else if c != nil {
		// If no specific discovery name, try to infer if only one is available from the container.
		// This block is only executed if 'c' is not nil.
		discoveries := c.Discoveries()
		if len(discoveries) == 1 {
			for _, d := range discoveries {
				discoveryClient = d
				break
			}
		} else if len(discoveries) > 1 {
			return nil, runtimeerrors.NewStructured(Module, "multiple discovery clients found in container, but no specific discovery client is configured for gRPC client").WithReason(commonv1.ErrorReason_VALIDATION_ERROR).WithCaller()
		}
	}

	// Apply discovery option if a client was found.
	if discoveryClient != nil {
		kratosOpts = append(kratosOpts, transgrpc.WithDiscovery(discoveryClient))
	}

	// Crucial check: If the endpoint implies discovery but no discovery client is configured.
	if strings.HasPrefix(endpoint, "discovery:///") && discoveryClient == nil {
		return nil, runtimeerrors.NewStructured(Module, "endpoint '%s' requires a discovery client, but none is configured", endpoint).WithReason(commonv1.ErrorReason_VALIDATION_ERROR).WithCaller()
	}

	// Configure node filters (selector).
	if selectorConfig := grpcConfig.GetSelector(); selectorConfig != nil {
		// Call the original, trusted NewFilter function from your app's selector package.
		nodeFilter, err := serviceselector.NewFilter(selectorConfig)
		if err != nil {
			return nil, runtimeerrors.WrapStructured(err, Module, "failed to create node filter").WithReason(commonv1.ErrorReason_INTERNAL_SERVER_ERROR).WithCaller()
		}

		if nodeFilter != nil {
			kratosOpts = append(kratosOpts, transgrpc.WithNodeFilter(nodeFilter))
		}
	}

	// Configure TLS.
	if tlsConfig := grpcConfig.GetTlsConfig(); tlsConfig != nil && tlsConfig.GetEnabled() {
		tlsCfg, err := servicetls.NewClientTLSConfig(tlsConfig)
		if err != nil {
			return nil, runtimeerrors.WrapStructured(err, Module, "failed to create client TLS config").WithReason(commonv1.ErrorReason_INTERNAL_SERVER_ERROR).WithCaller()
		}
		kratosOpts = append(kratosOpts, transgrpc.WithTLSConfig(tlsCfg))
	}

	// Apply any external native gRPC dial options passed via functional options.
	// These are applied last, allowing them to override previous options if needed.
	if len(clientOpts.GrpcDialOptions) > 0 {
		kratosOpts = append(kratosOpts, transgrpc.WithOptions(clientOpts.GrpcDialOptions...))
	}

	return kratosOpts, nil
}
