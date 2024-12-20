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
	"github.com/origadmin/runtime/log"
)

// SkipperServer returns a middleware that skips certain operations based on the provided configuration.
// It takes a Security configuration and a variable number of OptionSettings.
// If the Skipper is not configured, it returns nil and false.
func SkipperServer(cfg *configv1.Security, ss ...OptionSetting) (middleware.Middleware, bool) {
	log.Debugf("Skipper: creating middleware with config: %+v", cfg)
	// Apply default settings to the options
	option := settings.ApplyDefaultsOrZero(ss...)
	log.Debugf("Skipper: applied default settings to options: %+v", option)

	// If the Skipper is not configured, return immediately
	if option.Skipper == nil {
		log.Debugf("Skipper: skipper is not configured, returning nil and false")
		return nil, false
	}

	// Return a middleware that wraps the provided handler
	log.Debugf("Skipper: returning middleware with skipper: %+v", option.Skipper)
	return func(handler middleware.Handler) middleware.Handler {
		// Return a new handler that checks if the operation should be skipped
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			log.Debugf("Skipper: handling request: %+v", req)
			// If the Skipper is configured, check if the operation should be skipped
			if option.Skipper != nil {
				// Get the transport from the client context
				if tr, ok := transport.FromServerContext(ctx); ok {
					log.Debugf("Skipper: got transport from client context: %+v", tr)
					// todo: check the request method
					// If the operation should be skipped, create a new skip context and call the next handler
					if option.Skipper(tr.Operation()) {
						log.Debugf("Skipper: skipping request, creating new skip context")
						return handler(NewSkipContext(ctx), req)
					}
				} else {
					log.Debugf("Skipper: unable to get transport from client context")
				}
			} else {
				log.Debugf("Skipper: skipper is nil, not skipping request")
			}
			// If the operation should not be skipped, call the next handler with the original context
			log.Debugf("Skipper: not skipping request, calling next handler with original context")
			return handler(ctx, req)
		}
	}, true
}
