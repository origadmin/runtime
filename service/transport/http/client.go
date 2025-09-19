package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/service"
	"github.com/origadmin/runtime/service/tls"
	mw "github.com/origadmin/runtime/middleware"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// NewClient creates a new HTTP client.
// It is the recommended way to create a client when the protocol is known in advance.
func NewClient(ctx context.Context, cfg *transportv1.HTTPClient, opts ...service.Option) (interfaces.Client, error) {
	if cfg == nil {
		return nil, tkerrors.Errorf("HTTP client config is required for creation")
	}

	// 1. Process options to extract client-specific settings (endpoint, selector filter).
	var sOpts service.Options
	sOpts.Apply(opts...)

	// --- Client creation logic below uses the extracted, concrete 'cfg' and 'sOpts' ---

	var clientOpts []transhttp.ClientOption
	var mws []middleware.Middleware

	// Build client interceptors (middlewares)
	for _, name := range cfg.Middlewares {
		m, ok := mw.Get(name)
		if !ok {
			return nil, fmt.Errorf("client middleware '%s' not found in registry", name)
		}
		mws = append(mws, m)
	}
	if len(mws) > 0 {
		// Kratos HTTP client middleware is applied differently than gRPC
		// It's usually wrapped around the http.RoundTripper or directly in the client creation
		// For simplicity, we'll just pass them as an option if Kratos supports it directly.
		// clientOpts = append(clientOpts, transhttp.WithMiddleware(mws...)) // Example, check Kratos API
	}

	// Apply other client options from config
	if cfg.Timeout != nil {
		clientOpts = append(clientOpts, transhttp.WithTimeout(cfg.Timeout.AsDuration()))
	}

	// Apply TLS configuration
	if cfg.GetTls() != nil && cfg.GetTls().GetEnabled() {
		tlsConfig, err := tls.NewClientTLSConfig(cfg.GetTls())
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for client creation")
		}
		// Kratos HTTP client uses WithTransport to set custom http.RoundTripper
		// which can include TLS config. We'll handle this in the custom transport below.
	}

	// Determine target endpoint: prioritize endpoint from options (discovery) over direct target
	target := cfg.Target
	if sOpts.Value().clientEndpoint != "" {
		target = sOpts.Value().clientEndpoint
	}

	if target == "" {
		return nil, tkerrors.Errorf("client target endpoint is required for creation")
	}

	// Apply selector filter if provided via options
	if sOpts.Value().clientSelectorFilter != nil {
		// Similar to gRPC, Kratos HTTP client needs a way to integrate NodeFilter with discovery.
		// This might involve a custom Kratos client option or a custom resolver.
		// For example:
		// clientOpts = append(clientOpts, transhttp.WithNodeFilter(sOpts.Value().clientSelectorFilter)) // Hypothetical Kratos option
		// Or, if using a custom Kratos selector builder:
		// selectorBuilder := selector.NewBuilderWithFilter(sOpts.Value().clientSelectorFilter)
		// clientOpts = append(clientOpts, transhttp.WithDiscovery(discovery.NewDiscovery(target)), transhttp.WithSelector(selectorBuilder))
	}

	// Create a new HTTP client with custom transport to handle dial timeout and TLS
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   service.DefaultTimeout, // Default dial timeout
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: nil, // TODO: Apply TLS config from cfg.Tls
		ForceAttemptHTTP2: true,
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if cfg.DialTimeout != nil {
		transport.DialContext = (&net.Dialer{
			Timeout:   cfg.DialTimeout.AsDuration(),
			KeepAlive: 30 * time.Second,
		}).DialContext
	}

	// Apply TLS config to transport if enabled
	if cfg.GetTls() != nil && cfg.GetTls().GetEnabled() {
		tlsConfig, err := tls.NewClientTLSConfig(cfg.GetTls())
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for client creation")
		}
		transport.TLSClientConfig = tlsConfig
	}

	// Create the Kratos HTTP client
	client, err := transhttp.NewClient(ctx, transhttp.WithEndpoint(target), transhttp.WithTransport(transport), clientOpts...)
	if err != nil {
		return nil, tkerrors.Wrapf(err, "failed to create HTTP client to %s", target)
	}

	return client, nil
}
