package http

import (
	"net/url"
	"time"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service/selector"
	"github.com/origadmin/runtime/service/tls"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// DefaultServerMiddlewares returns the default server middlewares and options
func DefaultServerMiddlewares() []transhttp.ServerOption {
	return []transhttp.ServerOption{
		transhttp.Middleware(recovery.Recovery()),
		transhttp.ErrorEncoder(errors.NewErrorEncoder()),
	}
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
	opts := DefaultServerMiddlewares()

	// Add TLS configuration if needed
	if httpCfg.GetTls() != nil && httpCfg.GetTls().GetEnabled() {
		tlsConfig, err := tls.NewServerTLSConfig(httpCfg.GetTls())
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
	if httpCfg.ShutdownTimeout != nil {
		opts = append(opts, transhttp.ShutdownTimeout(httpCfg.ShutdownTimeout.AsDuration()))
	}

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
		opts = append(opts, transhttp.WithTimeout(time.Duration(timeout)*time.Millisecond))
	}

	// 4. Handle TLS
	if httpCfg.GetTls() != nil && httpCfg.GetTls().GetEnabled() {
		tlsConfig, err := tls.NewClientTLSConfig(httpCfg.GetTls())
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for client creation")
		}
		opts = append(opts, transhttp.WithTLSConfig(tlsConfig))
	}

	// 5. Process selectors
	if selectorCfg := cfg.GetSelector(); selectorCfg != nil {
		filter, err := selector.NewFilter(selectorCfg)
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid selector config for client creation")
		}
		opts = append(opts, transhttp.WithNodeFilter(filter))
	}

	return opts, nil
}
