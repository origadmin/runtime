package http

import (
	"context"
	"net/url"
	"time"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goexts/generic/configure"

	"github.com/origadmin/framework/runtime/interfaces"
	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service"
	serviceSelector "github.com/origadmin/runtime/service/selector"
	"github.com/origadmin/runtime/service/tls"
)

const (
	defaultTimeout = 5 * time.Second
	Scheme         = "http"
	hostName       = "HOST"
)

// httpProtocolFactory implements service.ProtocolFactory for HTTP.
type httpProtocolFactory struct{}

// NewClient creates a new HTTP client instance.
func (f *httpProtocolFactory) NewClient(ctx context.Context, cfg *configv1.Service, opts ...service.Option) (interface{}, error) {
	ll := log.NewHelper(log.With(log.GetLogger(), "module", "service/http"))
	ll.Debugf("Creating new HTTP client with config: %+v", cfg)

	// Create a new service.Options instance and apply the incoming options.
	svcOpts := &service.Options{ContextOptions: interfaces.ContextOptions{Context: ctx}}
	configure.Apply(svcOpts, opts)

	// Initialize client options
	var clientOpts []transhttp.ClientOption
	timeout := defaultTimeout

	// Apply configuration from configv1.Service
	if cfg.GetHttp() != nil {
		httpCfg := cfg.GetHttp()
		if httpCfg.GetTimeout() != 0 {
			timeout = time.Duration(httpCfg.GetTimeout() * 1e6) // Convert milliseconds to nanoseconds
		}

		// Configure TLS if needed
		if httpCfg.GetUseTls() {
			tlsConfig, err := tls.NewClientTLSConfig(httpCfg.GetTlsConfig())
			if err != nil {
				return nil, err
			}
			if tlsConfig != nil {
				clientOpts = append(clientOpts, transhttp.WithTLSConfig(tlsConfig))
			}
		}

		// Set endpoint if provided
		if httpCfg.GetEndpoint() != "" {
			clientOpts = append(clientOpts, transhttp.WithEndpoint(httpCfg.GetEndpoint()))
		}

		// Handle service discovery and selector
		if selectorCfg := cfg.GetSelector(); selectorCfg != nil {
			filter, err := serviceSelector.NewFilter(selectorCfg)
			if err == nil {
				clientOpts = append(clientOpts, transhttp.WithNodeFilter(filter))
			} else {
				ll.Warnf("Failed to create selector filter: %v", err)
			}
		}
	}

	// Apply timeout
	clientOpts = append(clientOpts, transhttp.WithTimeout(timeout))

	// Extract HTTP client specific options from the service.Options' Context.
	// These are the options added via http.WithClientOption
	clientOptsFromContext := FromClientOptions(svcOpts)
	clientOpts = append(clientOpts, clientOptsFromContext...)

	// Create the client with the merged options
	client, err := transhttp.NewClient(ctx, clientOpts...)
	if err != nil {
		return nil, errors.Newf(500, "INTERNAL_SERVER_ERROR", "create http client failed: %v", err)
	}

	return client, nil
}

// NewServer creates a new HTTP server instance.
func (f *httpProtocolFactory) NewServer(cfg *configv1.Service, opts ...service.Option) (interfaces.Server, error) {
	ll := log.NewHelper(log.With(log.GetLogger(), "module", "service/http"))
	ll.Debugf("Creating new HTTP server instance with config: %+v", cfg)

	// Create a new service.Options instance and apply the incoming options.
	// Initialize its ContextOptions.Context with a background context.
	svcOpts := &service.Options{ContextOptions: interfaces.ContextOptions{Context: context.Background()}}
	configure.Apply(svcOpts, opts)

	// Initialize Kratos HTTP server options
	var kratosServerOptions []transhttp.ServerOption

	// Add default recovery middleware and error encoder
	kratosServerOptions = append(kratosServerOptions, transhttp.Middleware(recovery.Recovery()))
	kratosServerOptions = append(kratosServerOptions, transhttp.ErrorEncoder(errors.NewErrorEncoder()))

	// Apply configuration from configv1.Service
	if cfg.GetHttp() != nil {
		httpCfg := cfg.GetHttp()

		if httpCfg.GetUseTls() {
			tlsConfig, err := tls.NewServerTLSConfig(httpCfg.GetTlsConfig())
			if err != nil {
				return nil, err
			}
			if tlsConfig != nil {
				kratosServerOptions = append(kratosServerOptions, transhttp.TLSConfig(tlsConfig))
			}
		}
		if httpCfg.GetNetwork() != "" {
			kratosServerOptions = append(kratosServerOptions, transhttp.Network(httpCfg.GetNetwork()))
		}
		if httpCfg.GetAddr() != "" {
			kratosServerOptions = append(kratosServerOptions, transhttp.Address(httpCfg.GetAddr()))
		}
		timeout := defaultTimeout
		if httpCfg.GetTimeout() != 0 {
			timeout = time.Duration(httpCfg.GetTimeout() * 1e6) // Convert milliseconds to nanoseconds
		}
		kratosServerOptions = append(kratosServerOptions, transhttp.Timeout(timeout))

		ll.Debugw("msg", "HTTP", "endpoint", httpCfg.GetEndpoint())
		if httpCfg.GetEndpoint() != "" {
			parsedEndpoint, err := url.Parse(httpCfg.GetEndpoint())
			if err == nil {
				kratosServerOptions = append(kratosServerOptions, transhttp.Endpoint(parsedEndpoint))
			} else {
				ll.Errorf("Failed to parse endpoint: %v", err)
			}
		}
	}

	// Extract HTTP specific options from the service.Options' Context.
	// These are the options added via http.WithServerOption
	serverOptsFromContext := FromServerOptions(svcOpts)

	// Combine all HTTP options
	httpOpts := append(kratosServerOptions, serverOptsFromContext...)

	return transhttp.NewServer(httpOpts...), nil
}

// init registers the HTTP protocol factory with the global service registry.
func init() {
	// Register the HTTP protocol factory with the service module.
	service.RegisterProtocol("http", &httpProtocolFactory{})
}
