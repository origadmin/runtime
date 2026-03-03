/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package context

import (
	"context"
)

//go:generate adptool .
//go:adapter:package context

type (
	traceKey     struct{}
	transKey     struct{}
	rowLockKey   struct{}
	idKey        struct{}
	tokenKey     struct{}
	createdByKey struct{}
)

// NewContext creates a new context with common runtime values.
func NewContext(ctx context.Context) Context {
	return ctx
}

// NewTrace creates a new context with the given trace ID.
func NewTrace(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceKey{}, traceID)
}

// FromTrace returns the trace ID from the context.
func FromTrace(ctx context.Context) string {
	if v, ok := ctx.Value(traceKey{}).(string); ok {
		return v
	}
	return ""
}

// NewRowLock creates a new context with the row lock flag.
func NewRowLock(ctx context.Context, lock bool) context.Context {
	return context.WithValue(ctx, rowLockKey{}, lock)
}

// FromRowLock returns the row lock flag from the context.
func FromRowLock(ctx context.Context) bool {
	if v, ok := ctx.Value(rowLockKey{}).(bool); ok {
		return v
	}
	return false
}

// NewCreatedBy creates a new context with the creator ID.
func NewCreatedBy(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, createdByKey{}, id)
}

// FromCreatedBy returns the creator ID from the context.
func FromCreatedBy(ctx context.Context) string {
	if v, ok := ctx.Value(createdByKey{}).(string); ok {
		return v
	}
	return ""
}

// NewID creates a new context with the object ID.
func NewID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, idKey{}, id)
}

// FromID returns the object ID from the context.
func FromID(ctx context.Context) string {
	if v, ok := ctx.Value(idKey{}).(string); ok {
		return v
	}
	return ""
}

// NewToken creates a new context with the token.
func NewToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey{}, token)
}

// FromToken returns the token from the context.
func FromToken(ctx context.Context) string {
	if v, ok := ctx.Value(tokenKey{}).(string); ok {
		return v
	}
	return ""
}

// NewTrans creates a new context with the transaction flag.
func NewTrans(ctx context.Context, trans bool) context.Context {
	return context.WithValue(ctx, transKey{}, trans)
}

// FromTrans returns the transaction flag from the context.
func FromTrans(ctx context.Context) bool {
	if v, ok := ctx.Value(transKey{}).(bool); ok {
		return v
	}
	return false
}
