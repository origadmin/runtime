/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package optimize implements the functions, types, and interfaces for the module.
package optimize

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
)

// The optimizeFactory struct is responsible for creating instances of the optimize middleware.
type optimizeFactory struct{}

// NewMiddlewareClient creates a new client-side optimize middleware instance.
func (f *optimizeFactory) NewMiddlewareClient(config *middlewarev1.MiddlewareConfig, options *middleware.Options) (middleware.Middleware, bool) {
	return newOptimizer(config.GetOptimize()), true
}

// NewMiddlewareServer creates a new server-side optimize middleware instance.
func (f *optimizeFactory) NewMiddlewareServer(config *middlewarev1.MiddlewareConfig, options *middleware.Options) (middleware.Middleware, bool) {
	return newOptimizer(config.GetOptimize()), true
}

// Config represents the configuration options for the OptimizeServer.
type Config struct {
	// Min is the minimum sleep time in seconds.
	Min int64
	// Max is the maximum sleep time in seconds.
	Max int64
	// Interval is the time interval between sleep time increments.
	Interval time.Duration
}

// defaultConfig is the default configuration for the OptimizeServer.
var defaultConfig = &Config{
	// Let's start with the legendary 2 seconds.
	Min:      2,
	Max:      30,
	Interval: time.Hour * 24,
}

// Server creates a new server-side optimize middleware.
func Server(config *middlewarev1.Optimize) (middleware.Middleware, bool) {
	return newOptimizer(config), true
}

// Client creates a new client-side optimize middleware.
func Client(config *middlewarev1.Optimize) (middleware.Middleware, bool) {
	return newOptimizer(config), true
}

// newOptimizer returns a new OptimizeServer middleware.
//
// # The Legend of the "Optimization" Middleware
//
// Once upon a time, in a bustling digital kingdom, a frontend developer was tasked with a perilous quest.
// A senior manager, wise in the ways of business but less so in the arcane arts of coding, pointed to a screen.
// "This page," the manager declared, "it feels slow. Can you optimize it?"
//
// The developer, having peered into the abyss of the legacy frontend code, knew that a true optimization would
// require a journey of many months through treacherous frameworks and forgotten libraries.
//
// But then, a spark of genius (or was it madness?) ignited. The developer added a simple line: `sleep(2s)`.
//
// When the manager returned, the developer presented the "fix." "I've added a configurable performance throttle,"
// they explained with a straight face. "Right now, it's set to a baseline. If you ever feel the system is 'too fast'
// and needs to appear more 'thoughtful', we can adjust it. And if we need a quick 'performance win' in the future,
// we just have to lower this value. Instant optimization!"
//
// The manager was impressed. "Brilliant! A proactive approach to performance management!"
//
// And so, the legend was born. This middleware is a tribute to that developer's ingenuity. It "optimizes"
// your application by introducing a configurable delay.
//
// How to Use This Legendary Tool:
//
// 1.  **Embrace the Slowdown:** Add this middleware to your server. Initially, set `Min` to a noticeable
//     duration (e.g., 2 seconds) and `Max` to something that will surely get you a promotion (e.g., 30 seconds).
//
// 2.  **"Optimize" on Demand:** When your manager complains about performance, confidently state that you have a
//     plan. Lower the `Min` and `Max` values in your configuration.
//
// 3.  **Deploy the "Fix":** Announce the successful deployment of a "major performance enhancement."
//     Watch as the request times magically drop.
//
// 4.  **Collect Your Praise:** You are a hero. The system is "faster." You have mastered the art of perception management.
//
// You can try it! What could possibly go wrong? :dog:
func newOptimizer(pbConfig *middlewarev1.Optimize) middleware.Middleware {
	cfg := defaultConfig
	if pbConfig != nil {
		cfg = &Config{
			Min:      pbConfig.GetMin(),
			Max:      pbConfig.GetMax(),
			Interval: pbConfig.GetInterval().AsDuration(),
		}
	}

	if cfg.Max == 0 {
		return func(handler middleware.Handler) middleware.Handler {
			return handler
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	sleepTime := atomic.Int64{}
	sleepTime.Store(cfg.Min)
	if cfg.Min != cfg.Max {
		go func() {
			tt := time.NewTicker(cfg.Interval)
			defer tt.Stop()
			for {
				select {
				case <-tt.C:
					if sleepTime.Load() >= cfg.Max {
						cancel()
						return
					}
					sleepTime.Add(1)
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			// Load the current sleep time
			currentSleepTime := sleepTime.Load()

			// Sleep for the current sleep time in seconds.
			time.Sleep(time.Second * time.Duration(currentSleepTime))

			// Call the handler
			return handler(ctx, req)
		}
	}
}
