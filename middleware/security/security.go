/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package security

import (
	"github.com/go-kratos/kratos/v2/middleware"

	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

func Middleware(cfg *configv1.Security) (middleware.Middleware, error) {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
	}, nil
}
