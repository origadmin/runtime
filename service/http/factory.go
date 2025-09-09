package http

import (
	"context"
	"net/url"
	"time"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goexts/generic/configure"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service"
	"github.com/origadmin/runtime/service/selector"
	"github.com/origadmin/runtime/service/tls"
)

const (
	defaultTimeout = 5 * time.Second
)

// httpProtocolFactory implements service.ProtocolFactory for HTTP.
type httpProtocolFactory struct{}

// NewClient creates a new HTTP client instance.
func (f *httpProtocolFactory) NewClient(ctx context.Context, cfg *configv1.Service, opts ...service.Option) (interfaces.Client, error) {
	ll := log.NewHelper(log.With(log.GetLogger(), "module", "service/http"))
	ll.Debugf("Creating new HTTP client with config: %+v", cfg)

	svcOpts := &service.Options{ContextOptions: interfaces.ContextOptions{Context: ctx}}
	configure.Apply(svcOpts, opts)

	var clientOpts []transhttp.ClientOption
	timeout := defaultTimeout

	if cfg.GetHttp() != nil {
		httpCfg := cfg.GetHttp()
		if httpCfg.GetTimeout() != 0 {
			timeout = time.Duration(httpCfg.GetTimeout() * 1e6)
		}

		if httpCfg.GetUseTls() {
			tlsConfig, err := tls.NewClientTLSConfig(httpCfg.GetTlsConfig())
			if err != nil {
				return nil, err
			}
			if tlsConfig != nil {
				clientOpts = append(clientOpts, transhttp.WithTLSConfig(tlsConfig))
			}
		}

		if httpCfg.GetEndpoint() != "" {
			clientOpts = append(clientOpts, transhttp.WithEndpoint(httpCfg.GetEndpoint()))
		}

		if selectorCfg := cfg.GetSelector(); selectorCfg != nil {
			filter, err := selector.NewFilter(selectorCfg)
			if err == nil {
				clientOpts = append(clientOpts, transhttp.WithNodeFilter(filter))
			} else {
				ll.Warnf("Failed to create selector filter: %v", err)
			}
		}
	}

	clientOpts = append(clientOpts, transhttp.WithTimeout(timeout))

	clientOptsFromContext := FromClientOptions(svcOpts)
	clientOpts = append(clientOpts, clientOptsFromContext...)

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

	svcOpts := &service.Options{ContextOptions: interfaces.ContextOptions{Context: context.Background()}}
	configure.Apply(svcOpts, opts)

	var kratosServerOptions []transhttp.ServerOption
	kratosServerOptions = append(kratosServerOptions, transhttp.Middleware(recovery.Recovery()))
	kratosServerOptions = append(kratosServerOptions, transhttp.ErrorEncoder(errors.NewErrorEncoder()))

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
			timeout = time.Duration(httpCfg.GetTimeout() * 1e6)
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

	serverOptsFromContext := FromServerOptions(svcOpts)

	httpOpts := append(kratosServerOptions, serverOptsFromContext...)

	return transhttp.NewServer(httpOpts...), nil
}

func init() {
	service.RegisterProtocol("http", &httpProtocolFactory{})
}
