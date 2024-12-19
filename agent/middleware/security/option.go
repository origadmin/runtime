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

const (
	// MetadataSecurityTokenKey is the default token key.
	MetadataSecurityTokenKey = "x-metadata-security-token-key"
	// MetadataSecuritySkipKey is the default skip key.
	MetadataSecuritySkipKey = "x-metadata-security-skip-key"
)

// TokenParser is a function that parses a token from the context.
type TokenParser func(context.Context, string) (security.Claims, error)

// ResponseWriter is a function that writes a response to the http.ResponseWriter.
type ResponseWriter func(context.Context, security.Claims) (string, error)

// Option is a struct that contains the settings for the security middleware.
type Option struct {
	// Authorizer is the authorizer used to authorize the request.
	Authorizer security.Authorizer
	// Authenticator is the authenticator used to authenticate the request.
	Authenticator security.Authenticator
	// Serializer is the serializer used to serialize the claims.
	Serializer security.Serializer
	// TokenKey is the key used to store the token in the context.
	TokenKey string
	// Scheme is the scheme used for the authorization header.
	Scheme string
	// HeaderAuthorize is the name of the authorization header.
	HeaderAuthorize string
	// SkipKey is the key used to skip authentication.
	SkipKey string
	// PublicPaths are the public paths that do not require authentication.
	PublicPaths []string
	// TokenParser is the parser used to parse the token from the context.
	TokenParser func(ctx context.Context) string
	// Parser is the parser used to parse the user claims.
	Parser security.UserClaimsParser
	// Skipper is the function used to skip authentication.
	Skipper func(string) bool
}

// OptionSetting is a function that sets an option.
type OptionSetting = func(option *Option)

// ApplyDefaults applies the default settings to the option.
func (o *Option) ApplyDefaults() {
	// Apply default token key if not set.
	if o.TokenKey == "" {
		o.TokenKey = MetadataSecurityTokenKey
	}
	// Apply default skip key if not set.
	if o.SkipKey == "" {
		o.SkipKey = MetadataSecuritySkipKey
	}
	// Apply default header authorize if not set.
	if o.HeaderAuthorize == "" {
		o.HeaderAuthorize = security.HeaderAuthorize
	}
	// Apply default scheme if not set.
	if o.Scheme == "" {
		o.Scheme = security.SchemeBearer.String()
	}
}

// WithConfig applies the configuration to the option.
func (o *Option) WithConfig(cfg *configv1.Security) *Option {
	paths := cfg.GetPublicPaths()
	paths = mergePublic(paths, o.PublicPaths...)
	o.PublicPaths = paths
	return o
}

// WithTokenParser sets the token parser.
func WithTokenParser(parser func(ctx context.Context) string) OptionSetting {
	return func(opt *Option) {
		opt.TokenParser = parser
	}
}

// ParserUserClaims parses the user claims from the context.
func (o *Option) ParserUserClaims(ctx context.Context, claims security.Claims) security.UserClaims {
	// TODO: implement parsing user claims
	return nil
}

// WithAuthenticator sets the authenticator.
func WithAuthenticator(authenticator security.Authenticator) OptionSetting {
	return func(opt *Option) {
		opt.Authenticator = authenticator
	}
}

// WithAuthorizer sets the authorizer.
func WithAuthorizer(authorizer security.Authorizer) OptionSetting {
	return func(opt *Option) {
		opt.Authorizer = authorizer
	}
}

// WithSkipper sets the public paths.
func WithSkipper(paths ...string) OptionSetting {
	return func(opt *Option) {
		opt.PublicPaths = mergePublic(paths)
	}
}

// WithTokenKey sets the token key.
func WithTokenKey(key string) OptionSetting {
	return func(opt *Option) {
		opt.TokenKey = key
	}
}

// WithSkipKey sets the skip key.
func WithSkipKey(key string) OptionSetting {
	return func(opt *Option) {
		opt.SkipKey = key
	}
}

// WithConfig sets the configuration.
func WithConfig(cfg *configv1.Security) OptionSetting {
	return func(option *Option) {
		option.WithConfig(cfg)
	}
}
