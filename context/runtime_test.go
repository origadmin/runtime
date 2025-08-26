/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package context

import (
	"fmt"
	"testing"
)

// unused check
var (
	_ = WithContext
	_ = TODO
	_ = WithDeadlineCause
	_ = WithCancelCause
	_ = WithTimeoutCause
	_ = WithValue
	_ = WithCancel
	_ = Cause
	_ = AfterFunc
	_ = WithoutCancel
	_ = WithTimeout
	_ = WithDeadline
	_ = NewRowLock
	_ = FromRowLock
	_ = NewCreatedBy
	_ = FromCreatedBy
	_ = Background
	_ = NewTrace
	_ = FromTrans
	_ = FromTrace
	_ = NewTrans
	_ = NewID
	_ = FromID
	_ = NewToken
	_ = FromToken
	_ = NewUserCache
	_ = FromUserCache
	_ = NewTag
	_ = FromTag
	_ = NewStack
	_ = FromStack
	_ = NewDB
	_ = FromDB
)

// Creates a new context with a trace ID
func TestCreatesNewContextWithTraceID(t *testing.T) {
	// Arrr! Let's create a new context and see if it has a trace ID!
	ctx := Background()
	traceID := "trace-123"
	newCtx := NewTrace(ctx, traceID)

	if newCtx == ctx {
		t.Error("Shiver me timbers! The context should be new!")
	}
}

// Uses the provided trace ID correctly
func TestUsesProvidedTraceIDCorrectly(t *testing.T) {
	// Ahoy! Let's check if the trace ID is used correctly!
	ctx := Background()
	traceID := "trace-456"
	newCtx := NewTrace(ctx, traceID)

	if newCtx.Value(traceCtx{}) != traceID {
		t.Error("Blimey! The trace ID wasn't set correctly!")
	}
}

// Returns a context that includes the trace ID
func TestReturnsContextWithTraceID(t *testing.T) {
	// Yo-ho-ho! Let's see if the context includes the trace ID!
	ctx := Background()
	traceID := "trace-789"
	newCtx := NewTrace(ctx, traceID)
	newCtx = NewTrace(ctx, "trace-788")
	newCtx = NewTrace(ctx, "trace-788")
	newCtx = NewTrace(ctx, traceID)
	if FromTrace(newCtx) != traceID {
		t.Error("Arrr! The context doesn't include the trace ID!")
	}
}

// Manages empty string trace ID
func TestManagesEmptyStringTraceID(t *testing.T) {
	// Ahoy matey! Let's see how it handles an empty string trace ID!
	ctx := Background()
	traceID := ""
	newCtx := NewTrace(ctx, traceID)

	if newCtx.Value(traceCtx{}) != traceID {
		t.Error("Blimey! The empty string trace ID wasn't set correctly!")
	}
}

// Deals with extremely long trace ID
func TestDealsWithExtremelyLongTraceID(t *testing.T) {
	// Yo-ho-ho! Let's see how it handles an extremely long trace ID!
	ctx := Background()
	traceID := string(make([]byte, 10000))
	newCtx := NewTrace(ctx, traceID)

	if newCtx.Value(traceCtx{}) != traceID {
		t.Error("Arrr! The extremely long trace ID wasn't set correctly!")
	}
}

// Handles special characters in trace ID
func TestHandlesSpecialCharactersInTraceID(t *testing.T) {
	// Avast ye! Let's see how it handles special characters in the trace ID!
	ctx := Background()
	traceID := "!@#$%^&*()_+"
	newCtx := NewTrace(ctx, traceID)
	newCtx = NewTrans(newCtx, "123")
	newCtx = NewID(newCtx, "123")
	newCtx = NewToken(newCtx, "NewToken")
	newCtx = NewUserCache(newCtx, "NewUserCache")
	newCtx = NewRowLock(newCtx)
	newCtx = WithValue(newCtx, "456", "789")
	newCtx = WithMapValue(newCtx, "456", "789")
	newCtx = WithMapValue(newCtx, "ggb", "ggb")
	t.Log(FromTrace(newCtx))
	if FromTrace(newCtx) != traceID {
		t.Error("Blimey! The special characters in the trace ID weren't set correctly!")
	} else {

	}
	if v, ok := FromTrans(newCtx); !ok {
		t.Error("Blimey! The special characters in the trace ID weren't set correctly!")
	} else {
		t.Log(v)
	}

	if v := newCtx.Value("456"); v == nil {
		t.Error("Blimey! The special characters in the trace ID weren't set correctly!")
	} else {
		t.Log(v)
	}

	if v := FromID(newCtx); v == "" {
		t.Error("Blimey! The special characters in the trace ID weren't set correctly!")
	} else {
		t.Log(v)
	}

	if v := FromToken(newCtx); v == "" {
		t.Error("Blimey! The special characters in the trace ID weren't set correctly!")
	} else {
		t.Log(v)
	}

	if v, ok := FromUserCache(newCtx); !ok || v == nil {
		t.Error("Blimey! The special characters in the trace ID weren't set correctly!")
	} else {
		t.Log(v)
	}
	if truth := FromRowLock(newCtx); !truth {
		t.Error("Blimey! The special characters in the trace ID weren't set correctly!")
	}
}

const N = 50000

func BenchmarkWithValue(b *testing.B) {
	ctx := Background()
	for i := 0; i < N; i++ {
		ctx = WithValue(ctx, fmt.Sprintf("key %d", i), fmt.Sprintf("value %d", i))
	}
	count := 0
	for i := 0; i < N; i++ {
		if v := ctx.Value(fmt.Sprintf("key %d", i)); v != nil {
			count++
		}
	}
	if count != N {
		b.Fatalf("expected %d, got %d", b.N, count)
	}
	b.StopTimer()
}

func BenchmarkWithMapValue(b *testing.B) {
	ctx := Background()
	for i := 0; i < N; i++ {
		ctx = WithMapValue(ctx, fmt.Sprintf("key %d", i), fmt.Sprintf("value %d", i))
	}
	count := 0
	for i := 0; i < N; i++ {
		if v := ctx.Value(fmt.Sprintf("key %d", i)); v != nil {
			count++
		}
	}
	if count != N {
		b.Fatalf("expected %d, got %d", b.N, count)
	}
	b.StopTimer()
}

func BenchmarkWithRandValue(b *testing.B) {
	ctx := Background()
	for i := 0; i < N; i++ {
		if i%2 == 0 {
			ctx = WithValue(ctx, fmt.Sprintf("key %d", i), fmt.Sprintf("value %d", i))
		} else {
			ctx = WithMapValue(ctx, fmt.Sprintf("key %d", i), fmt.Sprintf("value %d", i))
		}
	}
	count := 0
	for i := 0; i < N; i++ {
		if v := ctx.Value(fmt.Sprintf("key %d", i)); v != nil {
			count++
		}
	}
	if count != N {
		b.Fatalf("expected %d, got %d", b.N, count)
	}
	b.StopTimer()
}
