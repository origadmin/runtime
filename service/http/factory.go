package http

import (
	"net/url"
	"time"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goexts/generic/configure"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/context"
	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service"
	"github.com/origadmin/runtime/service/tls"
	tkerrors "github.com/origadmin/toolkits/errors"
)

const (
	defaultTimeout = 5 * time.Second
)

// httpProtocolFactory implements service.ProtocolFactory for HTTP.
type httpProtocolFactory struct{}

// NewClient creates a new HTTP client instance by delegating to the direct implementation.
func (f *httpProtocolFactory) NewClient(ctx context.Context, cfg *configv1.Service, opts ...service.Option) (interfaces.Client, error) {
	return NewClient(ctx, cfg, opts...)
}

// NewServer creates a new HTTP server instance.
func (f *httpProtocolFactory) NewServer(cfg *configv1.Service, opts ...service.Option) (interfaces.Server, error) {
	ll := log.NewHelper(log.With(log.GetLogger(), "module", "service/http"))
	ll.Debugf("Creating new HTTP server instance with config: %+v", cfg)

	svcOpts := &service.Options{ContextOptions: interfaces.ContextOptions{Context: context.Background()}}
	configure.Apply(svcOpts, opts)

	var kratosServerOptions []transhttp.ServerOption
	kratosServerOptions = append(kratosServerOptions, transhttp.Middleware(recovery.Recovery()))
	kratosServerOptions = append(kratosServerOptions, transhttp.ErrorEncoder(runtimeerrors.NewErrorEncoder())) // This is correct, as it's setting up the encoder for *external* errors

	if cfg.GetHttp() != nil {
		httpCfg := cfg.GetHttp()

		if httpCfg.GetUseTls() {
			tlsConfig, err := tls.NewServerTLSConfig(httpCfg.GetTlsConfig())
			if err != nil {
				// This error occurs during server creation, it's an internal error for this function
				return nil, tkerrors.Wrapf(err, "invalid TLS config for server creation") // 修正为 Wrapf
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
			if err != nil {
				// This error occurs during server creation, it's an internal error for this function
				return nil, tkerrors.Wrapf(err, "failed to parse endpoint for server creation")
			}
			kratosServerOptions = append(kratosServerOptions, transhttp.Endpoint(parsedEndpoint))
		}
	}

	serverOptsFromContext := FromServerOptions(svcOpts)

	httpOpts := append(kratosServerOptions, serverOptsFromContext...)

	return transhttp.NewServer(httpOpts...), nil
}

func init() {
	service.RegisterProtocol("http", &httpProtocolFactory{})
}
