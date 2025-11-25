package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/logging"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/log"
)

// loggingFactory implements middleware.Factory for the logging middleware.
type loggingFactory struct{}

// NewMiddlewareClient creates a new client-side logging middleware.
func (f *loggingFactory) NewMiddlewareClient(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	if !cfg.GetEnabled() {
		return nil, false
	}
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)

	//// Get logging-specific configuration from the Protobuf config.
	//loggingConfig := cfg.GetLogging()
	//if loggingConfig == nil {
	//	return nil, false
	//}

	helper.Debugf("enabling logging client middleware")

	// Kratos logging middleware expects kratosLog.Logger.
	// Assuming origadmin/runtime/log.Logger is compatible with kratos/v2/log.Logger interface.
	return logging.Client(mwOpts.Logger), true
}

// NewMiddlewareServer creates a new server-side logging middleware.
func (f *loggingFactory) NewMiddlewareServer(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	if !cfg.GetEnabled() {
		return nil, false
	}
	// Resolve common options once at the factory level.
	mwOpts := FromOptions(opts...)
	helper := log.NewHelper(mwOpts.Logger)

	// Get logging-specific configuration from the Protobuf config.
	//loggingConfig := cfg.GetLogging()
	//if loggingConfig == nil {
	//	return nil, false
	//}
	helper.Debugf("enabling logging server middleware")
	// Kratos logging middleware expects kratosLog.Logger.
	// Assuming origadmin/runtime/log.Logger is compatible with kratos/v2/log.Logger interface.
	return logging.Server(mwOpts.Logger), true
}
