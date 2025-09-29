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
	if httpCfg.GetNetwork() != "" {
		opts = append(opts, transhttp.Network(httpCfg.GetNetwork()))
	}

	if httpCfg.GetAddr() != "" {
		opts = append(opts, transhttp.Address(httpCfg.GetAddr()))
	}

	// Configure timeout
	timeout := 5 * time.Second
	if httpCfg.GetTimeout() != 0 {
		timeout = time.Duration(httpCfg.GetTimeout() * 1e6)
	}
	opts = append(opts, transhttp.Timeout(timeout))

	// Handle endpoint configuration
	ll.Debugw("msg", "HTTP", "endpoint", httpCfg.GetEndpoint())
	if httpCfg.GetEndpoint() != "" {
		parsedEndpoint, err := url.Parse(httpCfg.GetEndpoint())
		if err != nil {
			return nil, tkerrors.Wrapf(err, "failed to parse endpoint for server creation")
		}
		opts = append(opts, transhttp.Endpoint(parsedEndpoint))
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
	if target := httpCfg.GetTarget(); target != "" {
		opts = append(opts, transhttp.WithEndpoint(target))
	}

	// 3. Processing timeout
	if timeout := httpCfg.GetTimeout(); timeout > 0 {
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
