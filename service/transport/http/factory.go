package http

import (
	"context"
	"fmt"

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
func (f *httpProtocolFactory) NewServer(cfg *transportv1.Transport, opts ...service.Option) (interfaces.Server, error) {
	// 1. Extract the specific HTTP config from the container.
	httpConfig := cfg.GetHttp()
	if httpConfig == nil {
		return nil, fmt.Errorf("HTTP config is missing in transport container")
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

	// Create the HTTP server instance
	srv := khttp.NewServer(kOpts...)

	// Register business logic
	if httpRegistrar != nil {
		httpRegistrar.RegisterHTTP(context.Background(), srv)
	}

	return srv, nil
}

// NewClient creates a new HTTP client instance.
// This is a placeholder implementation for now.
func (f *httpProtocolFactory) NewClient(ctx context.Context, cfg *transportv1.Transport, opts ...service.Option) (interfaces.Client, error) {
	// TODO: Implement HTTP client creation logic
	return nil, fmt.Errorf("HTTP client creation not yet implemented")
}
