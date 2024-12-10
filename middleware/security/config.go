/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"github.com/origadmin/toolkits/security"

	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type TokenParser func(context.Context, string) (security.Claims, error)
type ResponseWriter func(context.Context, security.Claims) (string, error)

type ConfigOption struct {
	Authorizer       security.Authorizer
	Authenticator    security.Authenticator
	Serializer       security.Serializer
	SecurityTokenKey string
	SecuritySkipKey  string
	PublicPaths      []string
	TokenParser      func(ctx context.Context) string
	Skipper          func(string) bool
}

type ConfigOptionSetting = func(opt *ConfigOption)

func WithAuthenticator(authenticator security.Authenticator) ConfigOptionSetting {
	return func(opt *ConfigOption) {
		opt.Authenticator = authenticator
	}
}

func WithAuthorizer(authorizer security.Authorizer) ConfigOptionSetting {
	return func(opt *ConfigOption) {
		opt.Authorizer = authorizer
	}
}

func WithSkipper(paths ...string) ConfigOptionSetting {
	return func(opt *ConfigOption) {
		opt.PublicPaths = mergePublic(paths)
	}
}

func WithSecurityTokenKey(key string) ConfigOptionSetting {
	return func(opt *ConfigOption) {
		opt.SecurityTokenKey = key
	}
}

func WithSecuritySkipKey(key string) ConfigOptionSetting {
	return func(opt *ConfigOption) {
		opt.SecuritySkipKey = key
	}
}

func WithConfig(cfg *configv1.Security) ConfigOptionSetting {
	return func(opt *ConfigOption) {
		paths := cfg.GetPublicPaths()
		paths = mergePublic(paths, opt.PublicPaths...)
		opt.PublicPaths = paths
	}
}
