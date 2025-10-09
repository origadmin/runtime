package http

import (
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service/tls"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// NewHTTPServer creates a new concrete HTTP server instance based on the provided configuration.
// It returns *transhttp.Server, not the generic interfaces.Server.
func NewHTTPServer(httpConfig *transportv1.HttpServerConfig, serverOpts *ServerOptions) (*transhttp.Server, error) {
	ll := log.NewHelper(log.With(log.GetLogger(), "module", "service/http"))
	ll.Debugf("Creating new HTTP server instance with config: %+v", httpConfig)

	if httpConfig == nil {
		return nil, tkerrors.Errorf("HTTP server config is required for creation")
	}

	// Get the container instance. It will be nil if not provided in options.
	var c interfaces.Container
	if serverOpts.ServiceOptions != nil {
		c = serverOpts.ServiceOptions.Container
	}

	// Check if middlewares are configured.
	hasMiddlewaresConfigured := len(httpConfig.GetMiddlewares()) > 0

	// If middlewares are configured but no container is provided, return an error.
	// This consolidates the nil check for the container.
	if hasMiddlewaresConfigured && c == nil {
		return nil, fmt.Errorf("application container is required for server middlewares but not found in options")
	}

	var kOpts []transhttp.ServerOption

	// Configure middlewares.
	var mws []middleware.Middleware
	if hasMiddlewaresConfigured {
		// 'c' is guaranteed to be non-nil at this point due to the early check above.
		for _, name := range httpConfig.GetMiddlewares() {
			m, ok := c.ServerMiddleware(name)
			if !ok {
				return nil, fmt.Errorf("server middleware '%s' not found in container", name)
			}
			mws = append(mws, m)
		}
	} else {
		// If no specific middlewares are configured, use default ones from adapter.go.
		mws = DefaultServerMiddlewares()
	}

	if len(mws) > 0 {
		kOpts = append(kOpts, transhttp.Middleware(mws...))
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
	if httpConfig.ShutdownTimeout != nil {
		kratosOpts = append(kOpts, transhttp.ShutdownTimeout(httpConfig.ShutdownTimeout.AsDuration()))
	}

	// Apply TLS configuration
	if httpConfig.GetTls() != nil && httpConfig.GetTls().GetEnabled() {
		tlsConfig, err := tls.NewServerTLSConfig(httpConfig.GetTls())
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for server creation")
		}
		kratosOpts = append(kratosOpts, transhttp.TLSConfig(tlsConfig))
	}

	// Apply any external Kratos HTTP server options passed via functional options.
	// These are applied last, allowing them to override previous options if needed.
	if len(serverOpts.HttpServerOptions) > 0 {
		kratosOpts = append(kratosOpts, serverOpts.HttpServerOptions...)
	}

	// Create the HTTP server instance.
	srv := transhttp.NewServer(kratosOpts...)

	return srv, nil
}
