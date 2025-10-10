package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/logging"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/log"
)

// loggingFactory implements middleware.Factory for the logging middleware.
type loggingFactory struct{}

// NewMiddlewareClient creates a new client-side logging middleware.
func (f *loggingFactory) NewMiddlewareClient(cfg *middlewarev1.MiddlewareConfig, opts ...options.Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)

	// Get logging-specific configuration from the Protobuf config.
	loggingConfig := cfg.GetLogging()
	if loggingConfig == nil || !loggingConfig.GetEnabled() {
		return nil, false
	}

	helper.Info("enabling client logging middleware")

	// Kratos logging middleware expects kratosLog.Logger.
	// Assuming origadmin/runtime/log.Logger is compatible with kratos/v2/log.Logger interface.
	return logging.Client(mwOpts.Logger), true
}

// NewMiddlewareServer creates a new server-side logging middleware.
func (f *loggingFactory) NewMiddlewareServer(cfg *middlewarev1.MiddlewareConfig, opts ...options.Option) (KMiddleware, bool) {
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)

	// Get logging-specific configuration from the Protobuf config.
	loggingConfig := cfg.GetLogging()
	if loggingConfig == nil || !loggingConfig.GetEnabled() {
		return nil, false
	}
	helper.Info("enabling server logging middleware")
	// Kratos logging middleware expects kratosLog.Logger.
	// Assuming origadmin/runtime/log.Logger is compatible with kratos/v2/log.Logger interface.
	return logging.Server(mwOpts.Logger), true
}
