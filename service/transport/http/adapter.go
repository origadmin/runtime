package http

import (
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	httpv1 "github.com/origadmin/runtime/api/gen/go/config/transport/http/v1"
	"github.com/origadmin/runtime/context"
	runtimeerrors "github.com/origadmin/runtime/errors"
	serviceselector "github.com/origadmin/runtime/service/selector"
	servicetls "github.com/origadmin/runtime/service/tls"
)

const Module = "transport.http"

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

// getMiddlewares resolves and returns a slice of middlewares based on configuration.
// It checks for configured middleware names, retrieves them from the available map,
// and falls back to default middlewares if none are explicitly configured.
func getMiddlewares(
	configuredNames []string,
	availableMws map[string]middleware.Middleware,
	defaultMws []middleware.Middleware,
	mwType string, // "server" or "client" for error messages
) ([]middleware.Middleware, error) {
	if len(configuredNames) > 0 {
		if len(availableMws) == 0 {
			return nil, runtimeerrors.NewStructured(Module, "application container is required for %s middlewares but not found in options", mwType)
		}
		var mws []middleware.Middleware
		for _, name := range configuredNames {
			m, ok := availableMws[name]
			if !ok {
				return nil, runtimeerrors.NewStructured(Module, "%s middleware '%s' not found in options", mwType, name)
			}
			mws = append(mws, m)
		}
		return mws, nil
	}
	return defaultMws, nil
}

// initHttpServerOptions initialize the http server option
func initHttpServerOptions(httpConfig *httpv1.Server, serverOpts *ServerOptions) ([]transhttp.ServerOption, error) {
	// Prepare the Kratos HTTP server options.
	var kratosOpts []transhttp.ServerOption

	// Add CORS support, merging proto config with code-based options.
	if corsConfig := httpConfig.GetCors(); corsConfig != nil {
		corsHandler, err := NewCorsHandler(corsConfig, serverOpts.CorsOptions...)
		if err != nil {
			// Propagate the configuration error upwards.
			return nil, runtimeerrors.WrapStructured(err, Module, "failed to create CORS handler").WithCaller()
		}
		if corsHandler != nil {
			kratosOpts = append(kratosOpts, transhttp.Filter(corsHandler))
		}
	}

	if len(serverOpts.ServerMiddlewares) > 0 {
		// Configure middlewares.
		mws, err := getMiddlewares(httpConfig.GetMiddlewares(), serverOpts.ServerMiddlewares, DefaultServerMiddlewares(), "server")
		if err != nil {
			return nil, err
		}
		if len(mws) > 0 {
			kratosOpts = append(kratosOpts, transhttp.Middleware(mws...))
		}
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
	if httpConfig.TlsConfig.GetEnabled() {
		tlsCfg, err := servicetls.NewServerTLSConfig(httpConfig.TlsConfig)
		if err != nil {
			return nil, runtimeerrors.WrapStructured(err, Module, "failed to create server TLS config")
		}
		kratosOpts = append(kratosOpts, transhttp.TLSConfig(tlsCfg))
	}

	// Apply any external Kratos HTTP server options passed via functional options.
	// These are applied last, allowing them to override previous options if needed.
	if len(serverOpts.ServerOptions) > 0 {
		kratosOpts = append(kratosOpts, serverOpts.ServerOptions...)
	}

	return kratosOpts, nil
}

// initHttpClientOptions initialize http client options
func initHttpClientOptions(_ context.Context, httpConfig *httpv1.Client, clientOpts *ClientOptions) ([]transhttp.ClientOption, error) {
	// Prepare the Kratos HTTP client options.
	var kratosOpts []transhttp.ClientOption

	// Apply options from the protobuf configuration.
	if httpConfig.GetTimeout() != nil {
		kratosOpts = append(kratosOpts, transhttp.WithTimeout(httpConfig.GetTimeout().AsDuration()))
	}

	// Configure middlewares.
	mws, err := getMiddlewares(httpConfig.GetMiddlewares(), clientOpts.ClientMiddlewares, DefaultClientMiddlewares(), "client")
	if err != nil {
		return nil, err
	}
	if len(mws) > 0 {
		kratosOpts = append(kratosOpts, transhttp.WithMiddleware(mws...))
	}

	// Configure service discovery and endpoint.
	var discoveryClient registry.Discovery
	endpoint := httpConfig.GetEndpoint()

	// 1. Try to get discovery client by name from config
	if discoveryName := httpConfig.GetDiscoveryName(); discoveryName != "" {
		if d, ok := clientOpts.Discoveries[discoveryName]; ok {
			discoveryClient = d
		} else {
			return nil, runtimeerrors.NewStructured(Module, "discovery client '%s' not found in options", discoveryName)
		}
	} else {
		// 2. If no specific name, try to find a default or single discovery client
		if d, ok := clientOpts.Discoveries["default"]; ok {
			discoveryClient = d
		} else if len(clientOpts.Discoveries) == 1 {
			// If there's only one discovery client, use it as the default
			for _, d := range clientOpts.Discoveries { // Iterate once to get the single client
				discoveryClient = d
				break
			}
		}
	}

	// Validate endpoint and discovery client combination
	if strings.HasPrefix(endpoint, "discovery:///") && discoveryClient == nil {
		return nil, runtimeerrors.NewStructured(Module, "endpoint '%s' requires a discovery client, but none is configured", endpoint)
	}

	// Apply Kratos options
	if endpoint != "" {
		kratosOpts = append(kratosOpts, transhttp.WithEndpoint(endpoint))
	}
	if discoveryClient != nil {
		kratosOpts = append(kratosOpts, transhttp.WithDiscovery(discoveryClient))
	}

	// Configure node filters (selector).
	if selectorConfig := httpConfig.GetSelector(); selectorConfig != nil {
		nodeFilter, err := serviceselector.NewFilter(selectorConfig)
		if err != nil {
			return nil, runtimeerrors.WrapStructured(err, Module, "failed to create node filter")
		}
		if nodeFilter != nil {
			kratosOpts = append(kratosOpts, transhttp.WithNodeFilter(nodeFilter))
		}
	}

	// Configure TLS and apply it to the transport.
	if tlsConfig := httpConfig.GetTlsConfig(); tlsConfig != nil && tlsConfig.GetEnabled() {
		tlsCfg, err := servicetls.NewClientTLSConfig(tlsConfig)
		if err != nil {
			return nil, runtimeerrors.WrapStructured(err, Module, "failed to create client TLS config")
		}
		kratosOpts = append(kratosOpts, transhttp.WithTLSConfig(tlsCfg))
	}

	// Apply any external Kratos HTTP client options passed via functional options.
	if len(clientOpts.ClientOptions) > 0 {
		kratosOpts = append(kratosOpts, clientOpts.ClientOptions...)
	}

	return kratosOpts, nil
}
