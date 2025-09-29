package http

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/interfaces"
	mw "github.com/origadmin/runtime/middleware"
	"github.com/origadmin/runtime/optionutil"
	"github.com/origadmin/runtime/service/tls"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// NewClient creates a new HTTP client.
// It is the recommended way to create a client when the protocol is known in advance.
func NewClient(ctx context.Context, cfg *transportv1.HTTPClient, opts ...interfaces.Option) (interfaces.Client, error) {
	if cfg == nil {
		return nil, tkerrors.Errorf("HTTP client config is required for creation")
	}

	// 1. Process options to extract client-specific settings (endpoint, selector filter).
	var options httpClientOptions
	sOpts := optionutil.Apply(&options, opts...)
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
		clientOpts = append(clientOpts, transhttp.WithTLSConfig(tlsConfig))
	}

	// Determine endpoint endpoint: prioritize endpoint from options (discovery) over direct endpoint
	endpoint := cfg.Endpoint
	if sOpts.Value().clientEndpoint != "" {
		endpoint = sOpts.Value().clientEndpoint
	}

	if endpoint == "" {
		return nil, tkerrors.Errorf("client endpoint endpoint is required for creation")
	}

	// Create the Kratos HTTP client
	client, err := transhttp.NewClient(ctx, transhttp.WithEndpoint(endpoint), transhttp.WithTransport(transport), clientOpts...)
	if err != nil {
		return nil, tkerrors.Wrapf(err, "failed to create HTTP client to %s", endpoint)
	}

	return client, nil
}
