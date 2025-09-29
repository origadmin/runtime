package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/optionutil"
	"github.com/origadmin/runtime/service"
)

// httpProtocolFactory implements the service.ProtocolFactory for HTTP.
type httpProtocolFactory struct{}

// init registers this factory with the framework's protocol registry.
func init() {
	service.RegisterProtocol("http", &httpProtocolFactory{})
}

// NewServer creates a new HTTP server instance.
func (f *httpProtocolFactory) NewServer(cfg *transportv1.Server, opts ...interfaces.Option) (interfaces.Server, error) {
	// 1. Extract the specific HTTP server config from the container.
	httpConfig := cfg.GetHttp()
	if httpConfig == nil {
		return nil, fmt.Errorf("HTTP server config is missing in transport container")
	}

	// 2. Apply options to a ServiceOptions struct to get a configured context.
	initialServiceCfg := &service.ServiceOptions{}
	configuredContext := optionutil.Apply(initialServiceCfg, opts...)

	// 3. Retrieve the fully configured ServiceOptions from the context.
	configuredServiceCfg, ok := optionutil.ConfigFromContext[*service.ServiceOptions](configuredContext)
	if !ok {
		return nil, fmt.Errorf("failed to retrieve configured service options from context")
	}

	// 4. Get the registrar from the configured service options.
	httpRegistrar, ok := configuredServiceCfg.Registrar.(service.HTTPRegistrar)
	if !ok && configuredServiceCfg.Registrar != nil {
		return nil, fmt.Errorf("invalid registrar: expected service.HTTPRegistrar, got %T", configuredServiceCfg.Registrar)
	}

	// 5. Get the middleware provider.
	if configuredServiceCfg.MiddlewareProvider == nil {
		return nil, fmt.Errorf("middleware provider not found in options")
	}

	var kOpts []khttp.ServerOption
	var mws []middleware.Middleware

	// Build middleware chain using the provider
	for _, name := range httpConfig.Middlewares {
		m, ok := configuredServiceCfg.MiddlewareProvider.GetMiddleware(name)
		if !ok {
			return nil, fmt.Errorf("middleware '%s' not found via provider", name)
		}
		mws = append(mws, m)
	}
	if len(mws) > 0 {
		kOpts = append(kOpts, khttp.Middleware(mws...))
	}

	// Apply other server options from protobuf config
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
		httpRegistrar.RegisterHTTP(srv)
	}

	return srv, nil
}

// NewClient creates a new HTTP client instance.
func (f *httpProtocolFactory) NewClient(ctx context.Context, cfg *transportv1.Client, opts ...interfaces.Option) (interfaces.Client, error) {
	// 1. Extract the specific HTTP client config from the container.
	httpConfig := cfg.GetHttp()
	if httpConfig == nil {
		return nil, fmt.Errorf("HTTP client config is missing in transport container")
	}

	// 2. Apply options to get the configured context and service options.
	initialServiceCfg := &service.ServiceOptions{}
	configuredContext := optionutil.Apply(initialServiceCfg, opts...)
	configuredServiceCfg, ok := optionutil.ConfigFromContext[*service.ServiceOptions](configuredContext)
	if !ok {
		return nil, fmt.Errorf("failed to retrieve configured service options from context")
	}

	// 3. Get the middleware provider.
	if configuredServiceCfg.MiddlewareProvider == nil {
		return nil, fmt.Errorf("middleware provider not found in options")
	}

	var clientOpts []khttp.ClientOption
	var mws []middleware.Middleware

	// Build client interceptors (middlewares) using the provider
	for _, name := range httpConfig.Middlewares {
		m, ok := configuredServiceCfg.MiddlewareProvider.GetMiddleware(name)
		if !ok {
			return nil, fmt.Errorf("client middleware '%s' not found via provider", name)
		}
		mws = append(mws, m)
	}
	if len(mws) > 0 {
		clientOpts = append(clientOpts, khttp.WithMiddleware(mws...))
	}

	// Apply other client options from protobuf config
	if httpConfig.Timeout != nil {
		clientOpts = append(clientOpts, khttp.WithTimeout(httpConfig.Timeout.AsDuration()))
	}

	// Determine target endpoint: prioritize endpoint from options over direct target from config
	target := httpConfig.Target
	if configuredServiceCfg.ClientEndpoint != "" {
		target = configuredServiceCfg.ClientEndpoint
	}

	// Apply selector filter if provided via options
	if configuredServiceCfg.ClientSelectorFilter != nil {
		// TODO: Kratos HTTP client needs a way to integrate NodeFilter with discovery.
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
	client, err := khttp.NewClient(ctx, khttp.WithEndpoint(target), khttp.WithTransport(transport), clientOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client to %s: %w", target, err)
	}

	return client, nil
}
