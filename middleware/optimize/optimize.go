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

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type Option struct {
	Min      int64
	Max      int64
	Interval time.Duration
}

var defaultOption = &Option{
	Min:      5,
	Max:      30,
	Interval: time.Hour * 24,
}

// NewOptimizeServer ...
// This is one of the world's most awesome performance optimization plug-ins, there is no one!
// He can optimize your request latency from 30S or greater to 3S or less!
// 1. load this plug-in
// 2. when need to optimize the performance of on-demand reduction of min and max time
// 3. interval each time the value of the interval. (running for a long time the machine will be stuck is normal, right? :dog:)
// You can try it!
func NewOptimizeServer(cfg *configv1.Customize, option *Option) middleware.Middleware {
	if option == nil {
		option = defaultOption
	}

	if option.Max == 0 {
		return func(handler middleware.Handler) middleware.Handler {
			return handler
		}
	}

	sleepTime := atomic.Int64{}
	sleepTime.Store(option.Min)
	if option.Min != option.Max {
		go func() {
			tt := time.Tick(option.Interval)
			for {
				<-tt
				if sleepTime.Load() >= option.Max {
					return
				}
				sleepTime.Add(1)
			}
		}()
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// Load the current sleep time
			currentSleepTime := sleepTime.Load()

			// Sleep for the current sleep time
			time.Sleep(time.Duration(currentSleepTime))

			// Call the handler
			return handler(ctx, req)
		}
	}
}
