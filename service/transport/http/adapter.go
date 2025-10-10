package http

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
)

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
