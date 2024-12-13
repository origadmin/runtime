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

type Option struct {
	Authorizer    security.Authorizer
	Authenticator security.Authenticator
	Serializer    security.Serializer
	TokenKey      string
	SkipKey       string
	PublicPaths   []string
	TokenParser   func(ctx context.Context) string
	Skipper       func(string) bool
}

type OptionSetting = func(option *Option)

func (o *Option) ApplyDefaults() {
	if o.TokenKey == "" {
		o.TokenKey = MetadataSecurityTokenKey
	}
	if o.SkipKey == "" {
		o.SkipKey = MetadataSecuritySkipKey
	}
}
func (o *Option) WithConfig(cfg *configv1.Security) *Option {
	paths := cfg.GetPublicPaths()
	paths = mergePublic(paths, o.PublicPaths...)
	o.PublicPaths = paths
	return o
}

func WithAuthenticator(authenticator security.Authenticator) OptionSetting {
	return func(opt *Option) {
		opt.Authenticator = authenticator
	}
}

func WithAuthorizer(authorizer security.Authorizer) OptionSetting {
	return func(opt *Option) {
		opt.Authorizer = authorizer
	}
}

func WithSkipper(paths ...string) OptionSetting {
	return func(opt *Option) {
		opt.PublicPaths = mergePublic(paths)
	}
}

func WithTokenKey(key string) OptionSetting {
	return func(opt *Option) {
		opt.TokenKey = key
	}
}

func WithSkipKey(key string) OptionSetting {
	return func(opt *Option) {
		opt.SkipKey = key
	}
}

func WithConfig(cfg *configv1.Security) OptionSetting {
	return func(option *Option) {
		option.WithConfig(cfg)
	}
}
