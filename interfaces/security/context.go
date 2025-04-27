/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"context"
)

type tokenCtx struct{}

func NewTokenContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenCtx{}, token)
}

func TokenFromContext(ctx context.Context) string {
	if token, ok := ctx.Value(tokenCtx{}).(string); ok {
		return token
	}
	return ""
}

type claimCtx struct{}

func NewClaimsContext(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimCtx{}, claims)
}

func ClaimsFromContext(ctx context.Context) Claims {
	if claims, ok := ctx.Value(claimCtx{}).(Claims); ok {
		return claims
	}
	return nil
}

type userClaimsCtx struct{}

func NewPolicyContext(ctx context.Context, policy Policy) context.Context {
	return context.WithValue(ctx, userClaimsCtx{}, policy)
}

func PolicyFromContext(ctx context.Context) Policy {
	if claims, ok := ctx.Value(userClaimsCtx{}).(Policy); ok {
		return claims
	}
	return nil
}

// rootCtxKey is a type used to store a boolean value in the context.
type rootCtxKey struct{}

// WithRootContext returns a new context with the rootCtxKey set to true.
func WithRootContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, rootCtxKey{}, true)
}

// ContextIsRoot checks if the context has the rootCtxKey set to true.
func ContextIsRoot(ctx context.Context) bool {
	// Try to get the value from the context.
	if _, ok := ctx.Value(rootCtxKey{}).(bool); ok {
		// If the value is present and is a boolean, return true.
		return true
	}
	// Otherwise, return false.
	return false
}
