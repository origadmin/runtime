package grpc

import (
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	grpcv1 "github.com/origadmin/runtime/api/gen/go/config/transport/grpc/v1"
	runtimeerrors "github.com/origadmin/runtime/errors"
	serviceselector "github.com/origadmin/runtime/service/selector"
	servicetls "github.com/origadmin/runtime/service/tls"
)

const Module = "transport.grpc"

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

// getMiddlewares resolves and returns a slice of middlewares based on configuration.
// It checks for configured middleware names, retrieves them from the available map,
// and falls back to default middlewares if none are explicitly configured.
func getMiddlewares(
	configuredNames []string,
	availableMws map[string]middleware.Middleware,
	defaultMws []middleware.Middleware,
	mwType string, // "server" or "client" for error messages
) ([]middleware.Middleware, error) {
	if len(configuredNames) > 0 {
		if len(availableMws) == 0 {
			return nil, runtimeerrors.NewStructured(Module, "application container is required for %s middlewares but not found in options", mwType)
		}
		var mws []middleware.Middleware
		for _, name := range configuredNames {
			m, ok := availableMws[name]
			if !ok {
				return nil, runtimeerrors.NewStructured(Module, "%s middleware '%s' not found in options", mwType, name)
			}
			mws = append(mws, m)
		}
		return mws, nil
	}
	return defaultMws, nil
}

// initGrpcServerOptions initialize the grpc server option
func initGrpcServerOptions(grpcConfig *grpcv1.Server, serverOpts *ServerOptions) ([]transgrpc.ServerOption, error) {
	// Prepare the Kratos gRPC server options.
	var kratosOpts []transgrpc.ServerOption

	// Apply options from the protobuf configuration.
	if grpcConfig.GetAddr() != "" {
		kratosOpts = append(kratosOpts, transgrpc.Address(grpcConfig.GetAddr()))
	}
	// Configure TLS for server
	if grpcConfig.GetTlsConfig().GetEnabled() {
		tlsCfg, err := servicetls.NewServerTLSConfig(grpcConfig.GetTlsConfig())
		if err != nil {
			return nil, runtimeerrors.WrapStructured(err, Module, "failed to create server TLS config")
		}
		kratosOpts = append(kratosOpts, transgrpc.TLSConfig(tlsCfg))
	}

	if grpcConfig.GetTimeout() != nil {
		kratosOpts = append(kratosOpts, transgrpc.Timeout(grpcConfig.GetTimeout().AsDuration()))
	}

	// Configure middlewares.
	mws, err := getMiddlewares(grpcConfig.GetMiddlewares(), serverOpts.ServerMiddlewares, DefaultServerMiddlewares(), "server")
	if err != nil {
		return nil, err
	}
	if len(mws) > 0 {
		kratosOpts = append(kratosOpts, transgrpc.Middleware(mws...))
	}

	// Apply any external Kratos gRPC server options passed via functional options.
	// These are applied last, allowing them to override previous options if needed.
	if len(serverOpts.ServerOptions) > 0 {
		kratosOpts = append(kratosOpts, serverOpts.ServerOptions...)
	}

	return kratosOpts, nil
}

// initGrpcClientOptions initialize grpc client options
func initGrpcClientOptions(grpcConfig *grpcv1.Client, clientOpts *ClientOptions) ([]transgrpc.ClientOption, error) {
	// Prepare the Kratos gRPC client options.
	var kratosOpts []transgrpc.ClientOption

	// Apply options from the protobuf configuration.
	if grpcConfig.GetTimeout() != nil {
		kratosOpts = append(kratosOpts, transgrpc.WithTimeout(grpcConfig.GetTimeout().AsDuration()))
	}

	// Configure middlewares.
	mws, err := getMiddlewares(grpcConfig.GetMiddlewares(), clientOpts.ClientMiddlewares, DefaultClientMiddlewares(), "client")
	if err != nil {
		return nil, err
	}
	if len(mws) > 0 {
		kratosOpts = append(kratosOpts, transgrpc.WithMiddleware(mws...))
	}

	// Configure service discovery and endpoint.
	var discoveryClient registry.Discovery
	endpoint := grpcConfig.GetEndpoint()

	// 1. Try to get discovery client by name from config
	if discoveryName := grpcConfig.GetDiscoveryName(); discoveryName != "" {
		if d, ok := clientOpts.Discoveries[discoveryName]; ok {
			discoveryClient = d
		} else {
			return nil, runtimeerrors.NewStructured(Module, "discovery client '%s' not found in options", discoveryName)
		}
	} else {
		// 2. If no specific name, try to find a default or single discovery client
		if d, ok := clientOpts.Discoveries["default"]; ok {
			discoveryClient = d
		} else if len(clientOpts.Discoveries) == 1 {
			// If there's only one discovery client, use it as the default
			for _, d := range clientOpts.Discoveries { // Iterate once to get the single client
				discoveryClient = d
				break
			}
		}
	}

	// Validate endpoint and discovery client combination
	if strings.HasPrefix(endpoint, "discovery:///") && discoveryClient == nil {
		return nil, runtimeerrors.NewStructured(Module, "endpoint '%s' requires a discovery client, but none is configured", endpoint)
	}

	// Apply Kratos options
	if endpoint != "" {
		kratosOpts = append(kratosOpts, transgrpc.WithEndpoint(endpoint))
	}
	if discoveryClient != nil {
		kratosOpts = append(kratosOpts, transgrpc.WithDiscovery(discoveryClient))
	}

	// Configure node filters (selector).
	if selectorConfig := grpcConfig.GetSelector(); selectorConfig != nil {
		nodeFilter, err := serviceselector.NewFilter(selectorConfig)
		if err != nil {
			return nil, runtimeerrors.WrapStructured(err, Module, "failed to create node filter")
		}
		if nodeFilter != nil {
			kratosOpts = append(kratosOpts, transgrpc.WithNodeFilter(nodeFilter))
		}
	}

	// Configure TLS.
	if tlsConfig := grpcConfig.GetTlsConfig(); tlsConfig != nil && tlsConfig.GetEnabled() {
		tlsCfg, err := servicetls.NewClientTLSConfig(tlsConfig)
		if err != nil {
			return nil, runtimeerrors.WrapStructured(err, Module, "failed to create client TLS config")
		}
		kratosOpts = append(kratosOpts, transgrpc.WithTLSConfig(tlsCfg))
	}

	// Apply any external native gRPC dial options passed via functional options.
	// These are applied last, allowing them to override previous options if needed.
	if len(clientOpts.DialOptions) > 0 {
		kratosOpts = append(kratosOpts, transgrpc.WithOptions(clientOpts.DialOptions...))
	}

	return kratosOpts, nil
}
