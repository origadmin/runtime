package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/service"
	mw "github.com/origadmin/runtime/middleware"
)

// httpProtocolFactory implements the service.ProtocolFactory for HTTP.
type httpProtocolFactory struct{}

// init registers this factory with the framework's protocol registry.
func init() {
	service.RegisterProtocol("http", &httpProtocolFactory{})
}

// NewServer creates a new HTTP server instance.
// It conforms to the updated ProtocolFactory interface.
func (f *httpProtocolFactory) NewServer(cfg *transportv1.Server, opts ...service.Option) (interfaces.Server, error) {
	// 1. Extract the specific HTTP server config from the container.
	httpConfig := cfg.GetHttp()
	if httpConfig == nil {
		return nil, fmt.Errorf("HTTP server config is missing in transport container")
	}

	// 2. Process options to extract registrar.
	var sOpts service.Options
	sOpts.Apply(opts...)

	httpRegistrar, ok := sOpts.Value().registrar.(service.HTTPRegistrar)
	if !ok && sOpts.Value().registrar != nil {
		return nil, fmt.Errorf("invalid registrar: expected service.HTTPRegistrar, got %T", sOpts.Value().registrar)
	}

	// --- All creation logic below uses the extracted, concrete 'httpConfig' ---

	var kOpts []khttp.ServerOption
	var mws []middleware.Middleware

	// Build middleware chain
	for _, name := range httpConfig.Middlewares {
		m, ok := mw.Get(name)
		if !ok {
			return nil, fmt.Errorf("middleware '%s' not found in registry", name)
		}
		mws = append(mws, m)
	}
	if len(mws) > 0 {
		kOpts = append(kOpts, khttp.Middleware(mws...))
	}

	// Apply other server options
	if httpConfig.Network != "" {
		kOpts = append(kOpts, khttp.Network(httpConfig.Network))
	}
	if httpConfig.Addr != "" {
		kOpts = append(kOpts, khttp.Address(httpConfig.Addr))
	}
	if httpConfig.Timeout != nil {
		kOpts = append(kOpts, khttp.Timeout(httpConfig.Timeout.AsDuration()))
	}
	if httpConfig.ShutdownTimeout != nil {
		kOpts = append(kOpts, khttp.ShutdownTimeout(httpConfig.ShutdownTimeout.AsDuration()))
	}
	// TODO: Add TLS configuration

	// Create the HTTP server instance
	srv := khttp.NewServer(kOpts...)

	// Register business logic
	if httpRegistrar != nil {
		httpRegistrar.RegisterHTTP(context.Background(), srv)
	}

	return srv, nil
}

// NewClient creates a new HTTP client instance.
// It conforms to the updated ProtocolFactory interface.
func (f *httpProtocolFactory) NewClient(ctx context.Context, cfg *transportv1.Client, opts ...service.Option) (interfaces.Client, error) {
	// 1. Extract the specific HTTP client config from the container.
	httpConfig := cfg.GetHttp()
	if httpConfig == nil {
		return nil, fmt.Errorf("HTTP client config is missing in transport container")
	}

	// 2. Process options to extract client-specific settings (endpoint, selector filter).
	var sOpts service.Options
	sOpts.Apply(opts...)

	// --- Client creation logic below uses the extracted, concrete 'httpConfig' and 'sOpts' ---

	var clientOpts []khttp.ClientOption
	var mws []middleware.Middleware

	// Build client interceptors (middlewares)
	for _, name := range httpConfig.Middlewares {
		m, ok := mw.Get(name)
		if !ok {
			return nil, fmt.Errorf("client middleware '%s' not found in registry", name)
		}
		mws = append(mws, m)
	}
	if len(mws) > 0 {
		// Kratos HTTP client middleware is applied differently than gRPC
		// It's usually wrapped around the http.RoundTripper or directly in the client creation
		// For simplicity, let's assume a way to apply them here.
		// This part might need more specific Kratos HTTP client middleware integration.
		// For now, we'll just pass them as an option if Kratos supports it directly.
		// clientOpts = append(clientOpts, khttp.WithMiddleware(mws...)) // Example, check Kratos API
	}

	// Apply other client options
	if httpConfig.Timeout != nil {
		clientOpts = append(clientOpts, khttp.WithTimeout(httpConfig.Timeout.AsDuration()))
	}

	// Determine target endpoint: prioritize endpoint from options (discovery) over direct target
	target := httpConfig.Target
	if sOpts.Value().clientEndpoint != "" {
		target = sOpts.Value().clientEndpoint
	}

	// Apply selector filter if provided via options
	if sOpts.Value().clientSelectorFilter != nil {
		// Similar to gRPC, Kratos HTTP client needs a way to integrate NodeFilter with discovery.
		// This might involve a custom Kratos client option or a custom resolver.
		// For example:
		// clientOpts = append(clientOpts, khttp.WithNodeFilter(sOpts.Value().clientSelectorFilter)) // Hypothetical Kratos option
		// Or, if using a custom Kratos selector builder:
		// selectorBuilder := selector.NewBuilderWithFilter(sOpts.Value().clientSelectorFilter)
		// clientOpts = append(clientOpts, khttp.WithDiscovery(discovery.NewDiscovery(target)), khttp.WithSelector(selectorBuilder))
	}

	// Create a new HTTP client with custom transport to handle dial timeout and TLS
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   service.DefaultTimeout, // Default dial timeout
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: nil, // TODO: Apply TLS config from httpConfig.Tls
		ForceAttemptHTTP2: true,
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if httpConfig.DialTimeout != nil {
		transport.DialContext = (&net.Dialer{
			Timeout:   httpConfig.DialTimeout.AsDuration(),
			KeepAlive: 30 * time.Second,
		}).DialContext
	}

	// Create the Kratos HTTP client
	client, err := khttp.NewClient(ctx, khttp.WithEndpoint(target), khttp.WithTransport(transport), clientOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client to %s: %w", target, err)
	}

	// Return the client (which implements interfaces.Client if type aliased correctly)
	return client, nil
}
