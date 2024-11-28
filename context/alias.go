/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package context provides the context functions
package context

import (
	"context"
	"time"
)

// All context types are defined in the context package.
type (
	Context         = context.Context
	CancelFunc      = context.CancelFunc
	CancelCauseFunc = context.CancelCauseFunc
)

// TODO returns an empty Context. It is never canceled, has no deadline, and has no values.
func TODO() Context {
	return context.TODO()
}

// Background returns a non-nil, empty Context. It is never canceled, has no deadline, and has no values.
// Background is typically used by the main function, initialization, and tests, and as the top-level
// Context for incoming requests.
func Background() Context {
	return context.Background()
}

// WithCancel returns a copy of the parent Context with a new Done channel. The returned
// context's Done channel is closed when the returned cancel function is called or when the parent
// context's Done channel is closed, whichever happens first.
func WithCancel(parent Context) (Context, CancelFunc) {
	return context.WithCancel(parent)
}

// WithCancelCause returns a copy of the parent Context with a new Done channel and an associated
// cause. The returned context's Done channel is closed when the returned cancel function is called,
// when the parent context's Done channel is closed, or when the provided cause is returned by the
// parent context's Err method, whichever happens first.
func WithCancelCause(parent Context) (Context, CancelCauseFunc) {
	return context.WithCancelCause(parent)
}

// Cause returns the underlying cause of the context's cancellation or nil if it was not caused
// by another error.
func Cause(ctx Context) error {
	return context.Cause(ctx)
}

// AfterFunc waits for the duration to elapse and then calls f in its own goroutine. If the context
// is canceled before the duration elapses, f will not be called.
func AfterFunc(ctx Context, f func()) func() bool {
	return context.AfterFunc(ctx, f)
}

// WithoutCancel returns a non-cancellable Context derived from parent. It is never canceled, has no
// deadline, and has the same values as the parent. If the parent is already non-cancellable,
// WithoutCancel returns the parent.
func WithoutCancel(parent Context) Context {
	return context.WithoutCancel(parent)
}

// WithDeadline returns a copy of the parent Context with a deadline adjusted to be no later than d.
// If the parent's deadline is already earlier than d, WithDeadline(parent, d) is semantically
// equivalent to parent. The returned context's Done channel is closed when the deadline expires,
// when the returned cancel function is called, or when the parent context's Done channel is closed,
// whichever happens first.
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc) {
	return context.WithDeadline(parent, deadline)
}

// WithDeadlineCause returns a copy of the parent Context with a deadline adjusted to be no later than d
// and an associated cause. If the parent's deadline is already earlier than d, WithDeadlineCause(parent, d)
// is semantically equivalent to parent. The returned context's Done channel is closed when the deadline
// expires, when the returned cancel function is called, or when the provided cause is returned by the
// parent context's Err method, whichever happens first.
func WithDeadlineCause(parent Context, deadline time.Time, cause error) (Context, CancelFunc) {
	return context.WithDeadlineCause(parent, deadline, cause)
}

// WithTimeout returns a copy of the parent Context with a deadline adjusted to be no later than when
// the duration d elapses. It is equivalent to WithDeadline(parent, time.Now().Add(d)).
// The returned context's Done channel is closed when the timeout elapses, when the returned cancel
// function is called, or when the parent context's Done channel is closed, whichever happens first.
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

// WithTimeoutCause returns a copy of the parent Context with a deadline adjusted to be no later than
// when the duration d elapses and an associated cause. It is equivalent to WithDeadline(parent, time.Now().Add(d)).
// The returned context's Done channel is closed when the timeout elapses, when the returned cancel
// function is called, or when the provided cause is returned by the parent context's Err method,
// whichever happens first.
func WithTimeoutCause(parent Context, timeout time.Duration, cause error) (
	Context, CancelFunc,
) {
	return context.WithTimeoutCause(parent, timeout, cause)
}

// WithValue returns a copy of the parent Context that carries the given key-value pair. The value's
// type must be comparable for equality; otherwise, WithValue will panic. Use context Values only
// for request-scoped data that transits processes and API boundaries, not for passing optional parameters
// to functions.
func WithValue(parent Context, key, val any) Context {
	return context.WithValue(parent, key, val)
}
