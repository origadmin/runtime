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
)

// Config represents the configuration options for the OptimizeServer.
type Config struct {
	// Min is the minimum sleep time in seconds.
	Min int64
	// Max is the maximum sleep time in seconds.
	Max int64
	// Interval is the time interval between sleep time increments.
	Interval time.Duration
}

// defaultOption is the default configuration for the OptimizeServer.
var defaultConfig = &Config{
	Min:      5,
	Max:      30,
	Interval: time.Hour * 24,
}

// NewOptimizeServer returns a new OptimizeServer middleware.
//
// This is one of the world's most awesome performance optimization plug-ins, there is no one!
// He can optimize your request latency from 30S or greater to 3S or less!
//
// The OptimizeServer takes a configv1.Customize and a Config as input, and returns a middleware.CreateApp.
// If the Config is nil, it defaults to the defaultOption.
// If the Max sleep time is 0, the OptimizeServer returns a no-op middleware.
// Otherwise, it returns a middleware that sleeps for the current sleep time before calling the handler.
//
// 1. load this plug-in
// 2. when need to optimize the performance of on-demand reduction of min and max time
// 3. interval each time the value of the interval. (running for a long time the machine will be stuck is normal, right? :dog:)
// You can try it!
func NewOptimizeServer(ctx context.Context, config *Config) middleware.Middleware {
	if ctx == nil {
		ctx = context.Background()
	}
	if config == nil {
		config = defaultConfig
	}
	if config.Max == 0 {
		return func(handler middleware.Handler) middleware.Handler {
			return handler
		}
	}

	sleepTime := atomic.Int64{}
	sleepTime.Store(config.Min)
	if config.Min != config.Max {
		go func() {
			tt := time.Tick(config.Interval)
			for {
				select {
				case <-tt:
					if sleepTime.Load() >= config.Max {
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

			// Sleep for the current sleep time
			time.Sleep(time.Duration(currentSleepTime))

			// Call the handler
			return handler(ctx, req)
		}
	}
}
