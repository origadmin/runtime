package http

import (
	"context"
	"fmt"

	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
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
		return nil, service.ErrMissingServerConfig
	}

	// 2. Get all HTTP server-specific and common service-level options.
	serverOpts := FromServerOptions(opts)

	// Call the concrete server creation function.
	srv, err := NewHTTPServer(httpConfig, serverOpts)
	if err != nil {
		return nil, err
	}

	// Register pprof handlers if enabled
	if httpConfig.GetEnablePprof() {
		registerPprof(srv)
	}

	ctx := context.Background()
	// Register the user's business logic services if a registrar is provided.
	if serverOpts.ServiceOptions != nil && serverOpts.ServiceOptions.Registrar != nil {
		if httpRegistrar, ok := serverOpts.ServiceOptions.Registrar.(service.HTTPRegistrar); ok {
			httpRegistrar.RegisterHTTP(ctx, srv)
		} else {
			serverOpts.ServiceOptions.Registrar.Register(ctx, srv)
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

	// Call the concrete client creation function.
	client, err := NewHTTPClient(ctx, httpConfig, clientOpts)
	if err != nil {
		return nil, err
	}

	return client, nil
}
