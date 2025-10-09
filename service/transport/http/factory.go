package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/service"
)

// httpProtocolFactory implements the service.ProtocolFactory for HTTP.
type httpProtocolFactory struct{}

// init registers this factory with the framework's protocol registry.
func init() {
	service.RegisterProtocol(service.ProtocolHTTP, &httpProtocolFactory{})
}

// NewServer creates a new HTTP server instance based on the provided configuration.
func (f *httpProtocolFactory) NewServer(cfg *transportv1.Server, opts ...options.Option) (interfaces.Server, error) {
	// 1. Extract the specific HTTP server config from the transport configuration.
	httpConfig := cfg.GetHttp()
	if httpConfig == nil {
		return nil, fmt.Errorf("HTTP server config is missing in transport container")
	}

	// 2. Get all HTTP server-specific and common service-level options.
	serverOpts := FromServerOptions(opts)

	// Get the container instance. It will be nil if not provided in options.
	var c interfaces.Container
	if serverOpts.ServiceOptions != nil {
		c = serverOpts.ServiceOptions.Container
	}

	// Check if middlewares are configured.
	hasMiddlewaresConfigured := len(httpConfig.GetMiddlewares()) > 0

	// If middlewares are configured but no container is provided, return an error.
	// This consolidates the nil check for the container.
	if hasMiddlewaresConfigured && c == nil {
		return nil, fmt.Errorf("application container is required for server middlewares but not found in options")
	}

	var kOpts []transhttp.ServerOption

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
		kOpts = append(kOpts, transhttp.Middleware(mws...))
	}

	// Apply other server options from protobuf config
	if httpConfig.Network != "" {
		kOpts = append(kOpts, transhttp.Network(httpConfig.Network))
	}
	if httpConfig.Addr != "" {
		kOpts = append(kOpts, transhttp.Address(httpConfig.Addr))
	}
	if httpConfig.Timeout != nil {
		kOpts = append(kOpts, transhttp.Timeout(httpConfig.Timeout.AsDuration()))
	}
	if httpConfig.ShutdownTimeout != nil {
		kOpts = append(kOpts, transhttp.ShutdownTimeout(httpConfig.ShutdownTimeout.AsDuration()))
	}
	// TODO: Add TLS configuration

	// Create the HTTP server instance
	srv := transhttp.NewServer(kOpts...)

	// Register the user's business logic services if a registrar is provided.
	if serverOpts.ServiceOptions != nil && serverOpts.ServiceOptions.Registrar != nil {
		if httpRegistrar, ok := serverOpts.ServiceOptions.Registrar.(service.HTTPRegistrar); ok {
			httpRegistrar.RegisterHTTP(srv)
		} else {
			return nil, fmt.Errorf("invalid registrar: expected service.HTTPRegistrar, got %T", serverOpts.ServiceOptions.Registrar)
		}
	}

	return srv, nil
}

// NewClient creates a new HTTP client instance based on the provided configuration.
func (f *httpProtocolFactory) NewClient(ctx context.Context, cfg *transportv1.Client, opts ...options.Option) (interfaces.Client, error) {
	// 1. Extract the specific HTTP client config from the transport configuration.
	httpConfig := cfg.GetHttp()
	if httpConfig == nil {
		return nil, fmt.Errorf("HTTP client config is missing in transport container")
	}

	// 2. Get all HTTP client-specific and common service-level options.
	clientOpts := FromClientOptions(opts)

	// Get the container instance. It will be nil if not provided in options.
	var c interfaces.Container
	if clientOpts.ServiceOptions != nil {
		c = clientOpts.ServiceOptions.Container
	}

	// Determine if container-dependent features are configured.
	// ClientSelectorFilter might need the container for discovery.
	containerDependentFeaturesEnabled := len(httpConfig.GetMiddlewares()) > 0 || clientOpts.ClientSelectorFilter != nil

	// If container-dependent features are enabled but no container is provided (c is nil),
	// return an error immediately. This consolidates the nil checks for the container.
	if containerDependentFeaturesEnabled && c == nil {
		return nil, fmt.Errorf("application container is required for client configuration but not found in options")
	}

	var clientKratosOpts []transhttp.ClientOption
	var mws []middleware.Middleware

	// Configure middlewares.
	if len(httpConfig.GetMiddlewares()) > 0 {
		// 'c' is guaranteed to be non-nil at this point if middlewares are configured.
		for _, name := range httpConfig.GetMiddlewares() {
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
		clientKratosOpts = append(clientKratosOpts, transhttp.WithMiddleware(mws...))
	}

	// Apply other client options from protobuf config
	if httpConfig.Timeout != nil {
		clientKratosOpts = append(clientKratosOpts, transhttp.WithTimeout(httpConfig.Timeout.AsDuration()))
	}

	// Determine target endpoint: prioritize endpoint from options over direct target from config
	target := httpConfig.Endpoint
	if clientOpts.ClientEndpoint != "" {
		target = clientOpts.ClientEndpoint
	}

	// Apply selector filter if provided via options
	if clientOpts.ClientSelectorFilter != nil {
		// 'c' is guaranteed to be non-nil at this point if a selector filter is configured.
		// TODO: Kratos HTTP client needs a way to integrate NodeFilter with discovery.
		// For now, we'll just ensure the container is available if a filter is set.
		// Assuming the selector filter might need discovery from the container.
	}

	// Create a new HTTP client with custom transport to handle dial timeout and TLS
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // Default dial timeout
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig:       nil, // TODO: Apply TLS config from httpConfig.Tls
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if httpConfig.DialTimeout != nil {
		transport.DialContext = (&net.Dialer{
			Timeout:   httpConfig.DialTimeout.AsDuration(),
			KeepAlive: 30 * time.Second,
		}).DialContext
	}

	// Create the Kratos HTTP client
	client, err := transhttp.NewClient(ctx, transhttp.WithEndpoint(target), transhttp.WithTransport(transport), clientKratosOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client to %s: %w", target, err)
	}

	return client, nil
}
