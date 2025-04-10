/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"errors"

	"github.com/origadmin/toolkits/security"

	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

const (
	// MetadataSecurityTokenKey is the default token key.
	MetadataSecurityTokenKey = "x-md-global-security-token-key"
	// MetadataSecuritySkipKey is the default skip key.
	MetadataSecuritySkipKey = "x-md-global-security-skip-key"
)

// TokenParser is a function that parses a token from the context.
type TokenParser func(context.Context, string) (security.Claims, error)

// ResponseWriter is a function that writes a response to the http.ResponseWriter.
type ResponseWriter func(context.Context, security.Claims) (string, error)

// Option is a struct that contains the settings for the security middleware.
type Options struct {
	// Authorizer is the authorizer used to authorize the request.
	Authorizer security.Authorizer
	// Tokenizer is the authenticator used to authenticate the request.
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
	// PolicyParser is the parser used to parse the user claims.
	PolicyParser security.PolicyParser
	// Skipper is the function used to skip authentication.
	Skipper func(string) bool
	// IsRoot is the function used to check if the request is root.
	IsRoot func(ctx context.Context, claims security.Claims) bool
}

// Option is a function that sets an option.
type Option = func(option *Options)

// ApplyDefaults applies the default settings to the option.
func (o *Options) ApplyDefaults() {
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
	if o.TokenParser == nil {
		o.TokenParser = aggregateTokenParsers(
			TokenFromTransportClient(o.HeaderAuthorize, o.Scheme),
			TokenFromTransportServer(o.HeaderAuthorize, o.Scheme))
	}
	if o.IsRoot == nil {
		o.IsRoot = func(ctx context.Context, claims security.Claims) bool {
			return false
		}
	}
}

// WithConfig applies the configuration to the option.
func (o *Options) WithConfig(cfg *configv1.Security) *Options {
	paths := cfg.GetPublicPaths()
	paths = mergePublic(paths, o.PublicPaths...)
	o.PublicPaths = paths
	return o
}

// WithTokenParser sets the token parser.
func WithTokenParser(parser func(ctx context.Context) string) Option {
	return func(opt *Options) {
		opt.TokenParser = parser
	}
}

// ParsePolicy parses the user claims from the context.
func (o *Options) ParsePolicy(ctx context.Context, claims security.Claims) (security.Policy, error) {
	if o.PolicyParser == nil {
		return nil, errors.New("user claims parser is nil")
	}
	if claims == nil {
		claims = security.ClaimsFromContext(ctx)
	}
	return o.PolicyParser(ctx, claims)
}

// WithAuthenticator sets the token.
func WithAuthenticator(authenticator security.Authenticator) Option {
	return func(opt *Options) {
		opt.Authenticator = authenticator
	}
}

// WithAuthorizer sets the authorizer.
func WithAuthorizer(authorizer security.Authorizer) Option {
	return func(opt *Options) {
		opt.Authorizer = authorizer
	}
}

// WithSkipper sets the public paths.
func WithSkipper(paths ...string) Option {
	return func(opt *Options) {
		opt.PublicPaths = mergePublic(paths)
	}
}

// WithTokenKey sets the token key.
func WithTokenKey(key string) Option {
	return func(opt *Options) {
		opt.TokenKey = key
	}
}

// WithSkipKey sets the skip key.
func WithSkipKey(key string) Option {
	return func(opt *Options) {
		opt.SkipKey = key
	}
}

// WithConfig sets the configuration.
func WithConfig(cfg *configv1.Security) Option {
	return func(option *Options) {
		option.WithConfig(cfg)
	}
}

type AuthNSetting = func(authenticator *Authenticator)

func WithCache(cache security.CacheStorage) AuthNSetting {
	return func(authenticator *Authenticator) {
		authenticator.Cache = cache
	}
}

func WithScheme(scheme security.Scheme) AuthNSetting {
	return func(authenticator *Authenticator) {
		authenticator.Scheme = scheme
	}
}
