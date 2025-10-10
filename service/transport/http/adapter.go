package http

import (
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/interfaces"
	serviceselector "github.com/origadmin/runtime/service/selector"
	servicetls "github.com/origadmin/runtime/service/tls"
)

// DefaultServerMiddlewares provides a default set of server-side middlewares for HTTP services.
// These are essential for ensuring basic stability and observability.
func DefaultServerMiddlewares() []middleware.Middleware {
	return []middleware.Middleware{
		// recovery middleware recovers from panics and converts them into errors.
		recovery.Recovery(),
	}
}

// DefaultClientMiddlewares provides a default set of client-side middlewares for HTTP services.
func DefaultClientMiddlewares() []middleware.Middleware {
	return []middleware.Middleware{
		recovery.Recovery(),
	}
}

// initHttpServerOptions initialize the http server option
func initHttpServerOptions(httpConfig *transportv1.HttpServerConfig, serverOpts *ServerOptions) ([]transhttp.ServerOption, error) {
	// Prepare the Kratos HTTP server options.
	var kratosOpts []transhttp.ServerOption

	// Get the container instance. It will be nil if not provided in options.
	var c interfaces.Container
	if serverOpts.ServiceOptions != nil {
		c = serverOpts.ServiceOptions.Container
	}

	// Add CORS support, merging proto config with code-based options.
	if corsConfig := httpConfig.GetCors(); corsConfig != nil {
		corsHandler, err := NewCorsHandler(corsConfig, serverOpts.CorsOptions...)
		if err != nil {
			// Propagate the configuration error upwards.
			return nil, fmt.Errorf("failed to create CORS handler: %w", err)
		}
		if corsHandler != nil {
			kratosOpts = append(kratosOpts, transhttp.Filter(corsHandler))
		}
	}

	// Check if middlewares are configured.
	hasMiddlewaresConfigured := len(httpConfig.GetMiddlewares()) > 0

	// If middlewares are configured but no container is provided, return an error.
	// This consolidates the nil check for the container.
	if hasMiddlewaresConfigured && c == nil {
		return nil, fmt.Errorf("application container is required for server middlewares but not found in options")
	}

	// Configure middlewares.
	var mws []middleware.Middleware
	if hasMiddlewaresConfigured {
		// 'c' is guaranteed to be non-nil at this point due to the early check above.
		for _, name := range httpConfig.GetMiddlewares() {
			m, ok := c.ServerMiddleware(name)
			if !ok {
				return nil, fmt.Errorf("server middleware '%s' not found in container", name)
			}
			mws = append(mws, m)
		}
	} else {
		// If no specific middlewares are configured, use default ones from adapter.go.
		mws = DefaultServerMiddlewares()
	}

	if len(mws) > 0 {
		kratosOpts = append(kratosOpts, transhttp.Middleware(mws...))
	}

	// Apply other server options from protobuf config
	if httpConfig.Network != "" {
		kratosOpts = append(kratosOpts, transhttp.Network(httpConfig.Network))
	}
	if httpConfig.Addr != "" {
		kratosOpts = append(kratosOpts, transhttp.Address(httpConfig.Addr))
	}
	if httpConfig.Timeout != nil {
		kratosOpts = append(kratosOpts, transhttp.Timeout(httpConfig.Timeout.AsDuration()))
	}

	// Apply TLS configuration
	// Configure TLS for server
	if tlsConfig := httpConfig.GetTlsConfig(); tlsConfig != nil && tlsConfig.GetEnabled() {
		tlsCfg, err := servicetls.NewServerTLSConfig(tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create server TLS config: %w", err)
		}
		kratosOpts = append(kratosOpts, transhttp.TLSConfig(tlsCfg))
	}

	// Apply any external Kratos HTTP server options passed via functional options.
	// These are applied last, allowing them to override previous options if needed.
	if len(serverOpts.HttpServerOptions) > 0 {
		kratosOpts = append(kratosOpts, serverOpts.HttpServerOptions...)
	}

	return kratosOpts, nil
}

// initHttpClientOptions initialize http client options
func initHttpClientOptions(ctx context.Context, httpConfig *transportv1.HttpClientConfig, clientOpts *ClientOptions) ([]transhttp.ClientOption, error) {
	// Prepare the Kratos HTTP client options.
	var kratosOpts []transhttp.ClientOption

	// Get the container instance.
	var c interfaces.Container
	if clientOpts.ServiceOptions != nil {
		c = clientOpts.ServiceOptions.Container
	}

	// Centralized check for container dependency.
	if (len(httpConfig.GetMiddlewares()) > 0 || httpConfig.GetDiscoveryName() != "" || httpConfig.GetSelector() != nil) && c == nil {
		return nil, fmt.Errorf("application container is required for client configuration (middlewares, discovery, or selector) but not found in options")
	}

	// Apply options from the protobuf configuration.
	if httpConfig.GetTimeout() != nil {
		kratosOpts = append(kratosOpts, transhttp.WithTimeout(httpConfig.GetTimeout().AsDuration()))
	}

	// Configure middlewares.
	var mws []middleware.Middleware
	if len(httpConfig.GetMiddlewares()) > 0 {
		for _, name := range httpConfig.GetMiddlewares() {
			m, ok := c.ClientMiddleware(name)
			if !ok {
				return nil, fmt.Errorf("client middleware '%s' not found in container", name)
			}
			mws = append(mws, m)
		}
	} else {
		mws = DefaultClientMiddlewares()
	}
	if len(mws) > 0 {
		kratosOpts = append(kratosOpts, transhttp.WithMiddleware(mws...))
	}

	// Configure service discovery and endpoint.
	var discoveryClient registry.Discovery
	endpoint := httpConfig.GetEndpoint()
	if endpoint != "" {
		kratosOpts = append(kratosOpts, transhttp.WithEndpoint(endpoint))
	}

	if discoveryName := httpConfig.GetDiscoveryName(); discoveryName != "" {
		if d, ok := c.Discovery(discoveryName); ok {
			discoveryClient = d
		} else {
			return nil, fmt.Errorf("discovery client '%s' not found in container", discoveryName)
		}
	} else if c != nil {
		discoveries := c.Discoveries()
		if len(discoveries) == 1 {
			for _, d := range discoveries {
				discoveryClient = d
				break
			}
		}
	}

	if discoveryClient != nil {
		kratosOpts = append(kratosOpts, transhttp.WithDiscovery(discoveryClient))
	}

	if strings.HasPrefix(endpoint, "discovery:///") && discoveryClient == nil {
		return nil, fmt.Errorf("endpoint '%s' requires a discovery client, but none is configured", endpoint)
	}

	// Configure node filters (selector).
	if selectorConfig := httpConfig.GetSelector(); selectorConfig != nil {
		nodeFilter, err := serviceselector.NewFilter(selectorConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create node filter: %w", err)
		}
		if nodeFilter != nil {
			kratosOpts = append(kratosOpts, transhttp.WithNodeFilter(nodeFilter))
		}
	}

	// Configure TLS and apply it to the transport.
	if tlsConfig := httpConfig.GetTlsConfig(); tlsConfig != nil && tlsConfig.GetEnabled() {
		tlsCfg, err := servicetls.NewClientTLSConfig(tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create client TLS config: %w", err)
		}
		kratosOpts = append(kratosOpts, transhttp.WithTLSConfig(tlsCfg))
	}

	// Apply any external Kratos HTTP client options passed via functional options.
	if len(clientOpts.HttpClientOptions) > 0 {
		kratosOpts = append(kratosOpts, clientOpts.HttpClientOptions...)
	}

	return kratosOpts, nil
}
