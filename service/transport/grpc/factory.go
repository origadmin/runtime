package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/service"
	serviceselector "github.com/origadmin/runtime/service/selector"
	servicetls "github.com/origadmin/runtime/service/tls"
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
	grpcServerOpts := FromServerOptions(opts)
	// Access common service options via grpcServerOpts.ServiceOptions
	serviceOpts := grpcServerOpts.ServiceOptions

	// 3. Prepare the Kratos gRPC server options.
	var kratosOpts []transgrpc.ServerOption

	// 4. Apply options from the protobuf configuration.
	if grpcConfig.GetAddr() != "" {
		kratosOpts = append(kratosOpts, transgrpc.Address(grpcConfig.GetAddr()))
	}
	// 5. Configure TLS for server
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

	// 6. Configure middlewares.
	if len(grpcConfig.GetMiddlewares()) > 0 {
		if serviceOpts.Container == nil {
			return nil, fmt.Errorf("application container is required for middleware but not found in options")
		}
		var mws []middleware.Middleware
		for _, name := range grpcConfig.GetMiddlewares() {
			m, ok := serviceOpts.Container.ServerMiddleware(name)
			if !ok {
				return nil, fmt.Errorf("middleware '%s' not found in container", name)
			}
			mws = append(mws, m)
		}
		kratosOpts = append(kratosOpts, transgrpc.Middleware(mws...))
	}

	// 7. Apply any external Kratos gRPC server options passed via functional options.
	if len(grpcServerOpts.GrpcServerOptions) > 0 {
		kratosOpts = append(kratosOpts, grpcServerOpts.GrpcServerOptions...)
	}

	// 8. Create the Kratos gRPC server instance.
	srv := transgrpc.NewServer(kratosOpts...)

	// 9. Register the user's business logic services if a registrar is provided.
	if serviceOpts.Registrar != nil {
		if grpcRegistrar, ok := serviceOpts.Registrar.(service.GRPCRegistrar); ok {
			grpcRegistrar.RegisterGRPC(srv)
		} else {
			return nil, fmt.Errorf("invalid registrar: expected service.GRPCRegistrar, got %T", serviceOpts.Registrar)
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
	grpcClientOpts := FromClientOptions(opts)
	// Access common service options via grpcClientOpts.ServiceOptions
	serviceOpts := grpcClientOpts.ServiceOptions

	// 3. Prepare the Kratos gRPC client options.
	var kratosOpts []transgrpc.ClientOption

	// 4. Apply options from the protobuf configuration.
	if grpcConfig.GetTimeout() != nil {
		kratosOpts = append(kratosOpts, transgrpc.WithTimeout(grpcConfig.GetTimeout().AsDuration()))
	}

	// 5. Configure middlewares.
	if len(grpcConfig.GetMiddlewares()) > 0 {
		if serviceOpts.Container == nil {
			return nil, fmt.Errorf("application container is required for middleware but not found in options")
		}
		var mws []middleware.Middleware
		for _, name := range grpcConfig.GetMiddlewares() {
			m, ok := serviceOpts.Container.ClientMiddleware(name)
			if !ok {
				return nil, fmt.Errorf("client middleware '%s' not found in container", name)
			}
			mws = append(mws, m)
		}
		kratosOpts = append(kratosOpts, transgrpc.WithMiddleware(mws...))
	}

	// 6. Configure service discovery and endpoint.
	var discoveryClient registry.Discovery
	endpoint := grpcConfig.GetEndpoint()

	// Always apply the endpoint option.
	if endpoint != "" {
		kratosOpts = append(kratosOpts, transgrpc.WithEndpoint(endpoint))
	}

	// Determine the discovery client.
	if discoveryName := grpcConfig.GetDiscoveryName(); discoveryName != "" {
		if serviceOpts.Container == nil {
			return nil, fmt.Errorf("application container is required for named discovery client but not found in options")
		}
		discoveries := serviceOpts.Container.Discoveries()
		if d, ok := discoveries[discoveryName]; ok {
			discoveryClient = d
		} else {
			return nil, fmt.Errorf("discovery client '%s' not found in container", discoveryName)
		}
	} else if serviceOpts.Container != nil {
		// If no specific discovery name, try to infer if only one is available.
		discoveries := serviceOpts.Container.Discoveries()
		if len(discoveries) == 1 {
			for _, d := range discoveries {
				discoveryClient = d
				break
			}
		} else if len(discoveries) > 1 {
			return nil, fmt.Errorf("multiple discovery clients found in container, but no specific discovery client is configured for gRPC client")
		}
	}

	// Apply discovery option if a client was found.
	if discoveryClient != nil {
		kratosOpts = append(kratosOpts, transgrpc.WithDiscovery(discoveryClient))
	}

	// Crucial check: If the endpoint implies discovery but no discovery client is configured.
	if strings.HasPrefix(endpoint, "discovery:///") && discoveryClient == nil {
		return nil, fmt.Errorf("endpoint '%s' requires a discovery client, but none is configured", endpoint)
	}

	// 7. Configure node filters (selector).
	if selectorConfig := grpcConfig.GetSelector(); selectorConfig != nil {
		// Call the original, trusted NewFilter function from your app's selector package.
		nodeFilter, err := serviceselector.NewFilter(selectorConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create node filter: %w", err)
		}

		if nodeFilter != nil {
			// Pass the single filter to Kratos.
			kratosOpts = append(kratosOpts, transgrpc.WithNodeFilter(nodeFilter))
		}
	}

	// 8. Configure TLS.
	var conn interfaces.Client
	var err error

	tlsConfig := grpcConfig.GetTlsConfig()
	useSecureDial := false
	if tlsConfig != nil && tlsConfig.GetEnabled() {
		tlsCfg, err := servicetls.NewClientTLSConfig(tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create client TLS config: %w", err)
		}
		kratosOpts = append(kratosOpts, transgrpc.WithTLSConfig(tlsCfg))
		useSecureDial = true
	}

	// Apply any external native gRPC dial options passed via functional options.
	// These are applied after TLS config to allow overriding or adding to it.
	if len(grpcClientOpts.GrpcDialOptions) > 0 {
		kratosOpts = append(kratosOpts, transgrpc.WithOptions(grpcClientOpts.GrpcDialOptions...))
	}

	// 9. Create the Kratos gRPC client connection.
	if useSecureDial {
		conn, err = transgrpc.Dial(ctx, kratosOpts...)
	} else {
		conn, err = transgrpc.DialInsecure(ctx, kratosOpts...)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	return conn, nil
}
