/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"github.com/go-kratos/kratos/v2/metadata"

	"github.com/origadmin/runtime/context"
	"github.com/origadmin/toolkits/security"
)

type skipCtx struct{}

func NewSkipContext(ctx context.Context) context.Context {
	if SkipFromContext(ctx) {
		return ctx
	}
	return context.WithValue(ctx, skipCtx{}, true)
}

func SkipFromContext(ctx context.Context) bool {
	if _, ok := ctx.Value(skipCtx{}).(bool); ok {
		return true
	}
	return false
}

type claimCtx struct{}

func ClaimsFromContext(ctx context.Context) security.Claims {
	if claims, ok := ctx.Value(claimCtx{}).(security.Claims); ok {
		return claims
	}
	return nil
}

func NewClaimsContext(ctx context.Context, claims security.Claims) context.Context {
	return context.WithValue(ctx, claimCtx{}, claims)
}

func WithSkipContextClient(ctx context.Context, key string) context.Context {
	return metadata.AppendToClientContext(ctx, key, StringBoolTrue)
}

func WithSkipContextServer(ctx context.Context, key string) context.Context {
	if SkipFromContext(ctx) {
		return ctx
	}
	md, ok := metadata.FromServerContext(ctx)
	if ok && md.Get(key) == StringBoolTrue {
		return NewSkipContext(ctx)
	}
	return ctx
}
