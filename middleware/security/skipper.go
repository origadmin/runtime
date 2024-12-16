/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goexts/generic/settings"

	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

func Skipper(cfg *configv1.Security, ss ...OptionSetting) (middleware.Middleware, bool) {
	option := settings.ApplyDefaultsOrZero(ss...)
	if option.Skipper == nil {
		return nil, false
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if option.Skipper != nil {
				if tr, ok := transport.FromClientContext(ctx); ok {
					// todo: check the request method
					if option.Skipper(tr.Operation()) {
						return handler(NewSkipContext(ctx), req)
					}
				}
			}
			return handler(ctx, req)
		}
	}, true
}
