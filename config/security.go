/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"github.com/origadmin/toolkits/security"
)

func WithMiddlewareAuthenticator(authenticator security.Authenticator) MiddlewareOptionSetting {
	return func(opt *MiddlewareOption) {
		opt.Authenticator = authenticator
	}
}

func WithMiddlewareSkipper(paths ...string) MiddlewareOptionSetting {
	return func(opt *MiddlewareOption) {
		opt.PublicPaths = paths
	}
}
