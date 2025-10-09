package http

import (
	"net/url"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"

	corsv1 "github.com/origadmin/runtime/api/gen/go/middleware/v1/cors"
	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service/selector"
	"github.com/origadmin/runtime/service/tls"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// DefaultServerMiddlewares returns the default server middlewares
func DefaultServerMiddlewares() []middleware.Middleware {
	return []middleware.Middleware{
		recovery.Recovery(),
	}
}

// DefaultClientMiddlewares returns the default client middlewares
func DefaultClientMiddlewares() []middleware.Middleware {
	return []middleware.Middleware{
		recovery.Recovery(),
	}
}

// defaultCorsOptions returns default CORS options when none are provided
func defaultCorsOptions() []handlers.CORSOption {
	return []handlers.CORSOption{
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
		handlers.MaxAge(86400), // 24 hours
	}
}

// corsOptionsFromConfig converts a CORS config to CORS options
func corsOptionsFromConfig(cfg *corsv1.Cors) []handlers.CORSOption {
	var options []handlers.CORSOption

	// Use default options if no values are provided
	if len(cfg.AllowOrigins) == 0 {
		options = append(options, handlers.AllowedOrigins([]string{"*"}))
	} else {
		options = append(options, handlers.AllowedOrigins(cfg.AllowOrigins))
	}

	if len(cfg.AllowMethods) == 0 {
		options = append(options, handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"}))
	} else {
		options = append(options, handlers.AllowedMethods(cfg.AllowMethods))
	}

	if len(cfg.AllowHeaders) == 0 {
		options = append(options, handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}))
	} else {
		options = append(options, handlers.AllowedHeaders(cfg.AllowHeaders))
	}

	if len(cfg.ExposeHeaders) > 0 {
		options = append(options, handlers.ExposedHeaders(cfg.ExposeHeaders))
	}

	if cfg.AllowCredentials {
		options = append(options, handlers.AllowCredentials())
	}

	if cfg.MaxAge > 0 {
		options = append(options, handlers.MaxAge(int(cfg.MaxAge)))
	} else if cfg.MaxAge == 0 {
		// Default max age if not specified
		options = append(options, handlers.MaxAge(86400)) // 24 hours
	}

	return options
}

// NewCorsHandler creates a new CORS handler for Kratos HTTP servers
func NewCorsHandler(cfg *corsv1.Cors) transhttp.FilterFunc {
	if cfg == nil || !cfg.Enabled {
		return nil
	}

	// If CORS is enabled but no configuration is provided, use defaults
	var options []handlers.CORSOption
	if len(cfg.AllowOrigins) == 0 && len(cfg.AllowMethods) == 0 && len(cfg.AllowHeaders) == 0 {
		options = defaultCorsOptions()
	} else {
		options = corsOptionsFromConfig(cfg)
	}

	// Add credentials option if specified
	if cfg.AllowCredentials {
		options = append(options, handlers.AllowCredentials())
	}

	return handlers.CORS(options...)
}

// adaptServerConfig converts service configuration to Kratos HTTP server options.
func adaptServerConfig(cfg *transportv1.Server) ([]transhttp.ServerOption, error) {
	if cfg == nil {
		return nil, tkerrors.Errorf("server config is required for creation")
	}

	httpCfg := cfg.GetHttp()
	if httpCfg == nil {
		return nil, tkerrors.Errorf("http server config is required for creation")
	}

	ll := log.NewHelper(log.With(log.GetLogger(), "module", "service/http"))

	// Start with default middlewares and options
	var opts []transhttp.ServerOption

	// Add CORS configuration if needed
	if httpCfg.GetCors() != nil && httpCfg.GetCors().GetEnabled() {
		corsHandler := NewCorsHandler(httpCfg.GetCors())
		opts = append(opts, transhttp.Filter(corsHandler))
	}

	// Add TLS configuration if needed
	if httpCfg.GetTlsConfig() != nil && httpCfg.GetTlsConfig().GetEnabled() {
		tlsConfig, err := tls.NewServerTLSConfig(httpCfg.GetTlsConfig())
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for server creation")
		}
		opts = append(opts, transhttp.TLSConfig(tlsConfig))
	}

	// Add network and address configurations
	if httpCfg.Network != "" {
		opts = append(opts, transhttp.Network(httpCfg.Network))
	}

	if httpCfg.Addr != "" {
		opts = append(opts, transhttp.Address(httpCfg.Addr))
	}

	// Configure timeout
	timeout := 5 * time.Second
	if httpCfg.Timeout != nil {
		timeout = httpCfg.Timeout.AsDuration()
	}
	opts = append(opts, transhttp.Timeout(timeout))

	// Configure shutdown timeout
	//if httpCfg.ShutdownTimeout != nil {
	//	opts = append(opts, transhttp.ShutdownTimeout(httpCfg.ShutdownTimeout.AsDuration()))
	//}

	// Handle endpoint configuration
	ll.Debugw("msg", "HTTP", "address", httpCfg.Addr)
	if httpCfg.Addr != "" {
		parsedAddr, err := url.Parse("http://" + httpCfg.Addr)
		if err != nil {
			return nil, tkerrors.Wrapf(err, "failed to parse address for server creation")
		}
		opts = append(opts, transhttp.Endpoint(parsedAddr))
	}

	return opts, nil
}

// adaptClientConfig converts service configuration to Kratos HTTP client options.
func adaptClientConfig(cfg *transportv1.Client) ([]transhttp.ClientOption, error) {
	// 1. Validate the configuration
	if cfg == nil {
		return nil, tkerrors.Errorf("client config is required for creation")
	}

	httpCfg := cfg.GetHttp()
	if httpCfg == nil {
		return nil, tkerrors.Errorf("http client config is required for creation")
	}

	var opts []transhttp.ClientOption

	// 2. Processing endpoints
	if httpCfg.Endpoint != "" {
		opts = append(opts, transhttp.WithEndpoint(httpCfg.Endpoint))
	}

	// 3. Processing timeout
	if httpCfg.Timeout != nil {
		opts = append(opts, transhttp.WithTimeout(httpCfg.Timeout.AsDuration()))
	}

	// 4. Handle TLS
	if httpCfg.GetTlsConfig() != nil && httpCfg.GetTlsConfig().GetEnabled() {
		tlsConfig, err := tls.NewClientTLSConfig(httpCfg.GetTlsConfig())
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for client creation")
		}
		opts = append(opts, transhttp.WithTLSConfig(tlsConfig))
	}

	// 5. Process selectors
	if selectorCfg := httpCfg.GetSelector(); selectorCfg != nil {
		filter, err := selector.NewFilter(selectorCfg)
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid selector config for client creation")
		}
		opts = append(opts, transhttp.WithNodeFilter(filter))
	}

	return opts, nil
}
