package http

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service"
	"github.com/origadmin/runtime/service/tls"
	mw "github.com/origadmin/runtime/middleware"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// NewServer creates a new HTTP server with the given configuration and options.
// It is the recommended way to create a server when the protocol is known in advance.
func NewServer(cfg *transportv1.HTTPServer, opts ...service.Option) (*transhttp.Server, error) {
	ll := log.NewHelper(log.With(log.GetLogger(), "module", "service/http"))
	ll.Debugf("Creating new HTTP server instance with config: %+v", cfg)

	if cfg == nil {
		return nil, tkerrors.Errorf("HTTP server config is required for creation")
	}

	// 1. Process options to extract registrar.
	var sOpts service.Options
	sOpts.Apply(opts...)

	httpRegistrar, ok := sOpts.Value().registrar.(service.HTTPRegistrar)
	if !ok && sOpts.Value().registrar != nil {
		return nil, fmt.Errorf("invalid registrar: expected service.HTTPRegistrar, got %T", sOpts.Value().registrar)
	}

	// --- Server creation logic below uses the extracted, concrete 'cfg' ---

	var kOpts []transhttp.ServerOption
	var mws []middleware.Middleware

	// Build middleware chain
	for _, name := range cfg.Middlewares {
		m, ok := mw.Get(name)
		if !ok {
			return nil, fmt.Errorf("middleware '%s' not found in registry", name)
		}
		mws = append(mws, m)
	}
	if len(mws) > 0 {
		kOpts = append(kOpts, transhttp.Middleware(mws...))
	}

	// Apply other server options
	if cfg.Network != "" {
		kOpts = append(kOpts, transhttp.Network(cfg.Network))
	}
	if cfg.Addr != "" {
		kOpts = append(kOpts, transhttp.Address(cfg.Addr))
	}
	if cfg.Timeout != nil {
		kOpts = append(kOpts, transhttp.Timeout(cfg.Timeout.AsDuration()))
	}
	if cfg.ShutdownTimeout != nil {
		kOpts = append(kOpts, transhttp.ShutdownTimeout(cfg.ShutdownTimeout.AsDuration()))
	}

	// Apply TLS configuration
	if cfg.GetTls() != nil && cfg.GetTls().GetEnabled() {
		tlsConfig, err := tls.NewServerTLSConfig(cfg.GetTls())
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for server creation")
		}
		kOpts = append(kOpts, transhttp.TLSConfig(tlsConfig))
	}

	// Create the HTTP server instance
	srv := transhttp.NewServer(kOpts...)

	// Register business logic
	if httpRegistrar != nil {
		httpRegistrar.RegisterHTTP(context.Background(), srv)
	}

	return srv, nil
}
