package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"

	"github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	serviceselector "github.com/origadmin/runtime/service/selector"
	servicetls "github.com/origadmin/runtime/service/tls"
)

// NewGRPCClient creates a new concrete gRPC client connection based on the provided configuration.
// It returns *transgrpc.ClientConn.
func NewGRPCClient(ctx context.Context, grpcConfig *transportv1.GrpcClientConfig, clientOpts *ClientOptions) (*grpc.ClientConn, error) {
	// Prepare the Kratos gRPC client options.
	var kratosOpts []transgrpc.ClientOption

	// Check in advance and get Container instance
	var c interfaces.Container
	if clientOpts.ServiceOptions != nil {
		c = clientOpts.ServiceOptions.Container
	}

	// Return error if middlewares are configured but no container is provided.
	if len(grpcConfig.GetMiddlewares()) > 0 && c == nil {
		return nil, fmt.Errorf("application container is required for middleware but not found in options")
	}

	// Return error if a named discovery client is configured but no container is provided.
	if discoveryName := grpcConfig.GetDiscoveryName(); discoveryName != "" && c == nil {
		return nil, fmt.Errorf("application container is required for named discovery client but not found in options")
	}

	// Apply options from the protobuf configuration.
	if grpcConfig.GetTimeout() != nil {
		kratosOpts = append(kratosOpts, transgrpc.WithTimeout(grpcConfig.GetTimeout().AsDuration()))
	}

	// Configure middlewares.
	var mws []middleware.Middleware
	if len(grpcConfig.GetMiddlewares()) > 0 {
		// 'c' is guaranteed to be non-nil at this point.
		for _, name := range grpcConfig.GetMiddlewares() {
			m, ok := c.ClientMiddleware(name)
			if !ok {
				return nil, fmt.Errorf("client middleware '%s' not found in container", name)
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
		// 'c' is guaranteed to be non-nil at this point.
		if d, ok := c.Discovery(discoveryName); ok {
			discoveryClient = d
		} else {
			return nil, fmt.Errorf("discovery client '%s' not found in container", discoveryName)
		}
	} else if c != nil {
		// If no specific discovery name, try to infer if only one is available.
		discoveries := c.Discoveries()
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

	// Configure node filters (selector).
	if selectorConfig := grpcConfig.GetSelector(); selectorConfig != nil {
		// Call the original, trusted NewFilter function from your app's selector package.
		nodeFilter, err := serviceselector.NewFilter(selectorConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create node filter: %w", err)
		}

		if nodeFilter != nil {
			kratosOpts = append(kratosOpts, transgrpc.WithNodeFilter(nodeFilter))
		}
	}

	// Configure TLS.
	useSecureDial := false
	if tlsConfig := grpcConfig.GetTlsConfig(); tlsConfig != nil && tlsConfig.GetEnabled() {
		tlsCfg, err := servicetls.NewClientTLSConfig(tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create client TLS config: %w", err)
		}
		kratosOpts = append(kratosOpts, transgrpc.WithTLSConfig(tlsCfg))
		useSecureDial = true
	}

	// Apply any external native gRPC dial options passed via functional options.
	// These are applied last, allowing them to override previous options if needed.
	if len(clientOpts.GrpcDialOptions) > 0 {
		kratosOpts = append(kratosOpts, transgrpc.WithOptions(clientOpts.GrpcDialOptions...))
	}

	// Create the Kratos gRPC client connection.
	var conn *grpc.ClientConn
	var err error
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
