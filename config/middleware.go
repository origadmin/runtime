/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"github.com/origadmin/toolkits/security"

	"github.com/origadmin/runtime/context"
)

type TokenParser func(context.Context, string) (security.Claims, error)
type ResponseWriter func(context.Context, security.Claims) (string, error)

type MiddlewareOption struct {
	Authorizer       security.Authorizer
	Authenticator    security.Authenticator
	Serializer       security.Serializer
	SecurityTokenKey string
	SecuritySkipKey  string
	PublicPaths      []string
	TokenParser      func(ctx context.Context) string
	Skipper          func(string) bool
}

type MiddlewareOptionSetting = func(opt *MiddlewareOption)
