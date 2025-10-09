package grpc

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
)

// DefaultServerMiddlewares returns the default server middlewares as raw middleware.Middleware slice.
func DefaultServerMiddlewares() []middleware.Middleware {
	return []middleware.Middleware{
		recovery.Recovery(),
	}
}

// DefaultClientMiddlewares returns the default client middlewares as raw middleware.Middleware slice.
func DefaultClientMiddlewares() []middleware.Middleware {
	return []middleware.Middleware{
		recovery.Recovery(),
	}
}
