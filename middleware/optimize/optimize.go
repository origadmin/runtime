package optimize

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"google.golang.org/protobuf/types/known/durationpb"

	optimizev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/optimize/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	"github.com/origadmin/runtime/extension/customize"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/log"
)

// optimizeFactory implements middleware.Factory for the optimize middleware.
type optimizeFactory struct{}

// NewMiddlewareClient creates a new client-side optimize middleware instance.
// Modified to comply with the Factory interface definition
func (f *optimizeFactory) NewMiddlewareClient(cfg *middlewarev1.Middleware,
	opts ...options.Option) (middleware.Middleware, bool) {
	// Parse options using FromOptions
	logger := log.FromOptions(opts) // FIX: Changed opts... to opts
	helper := log.NewHelper(logger)
	helper.Debugf("enabling client optimize middleware")

	// Get optimize configuration
	return nil, false
}

// NewMiddlewareServer creates a new server-side optimize middleware instance.
// Modified to comply with the Factory interface definition
func (f *optimizeFactory) NewMiddlewareServer(cfg *middlewarev1.Middleware, opts ...options.Option) (middleware.Middleware, bool) {
	// Parse options using FromOptions
	logger := log.FromOptions(opts) // FIX: Changed opts... to opts
	helper := log.NewHelper(logger)
	helper.Debugf("enabling server optimize middleware")

	// Check if custom configuration is enabled
	if !cfg.GetEnabled() || cfg.GetType() != "customize" || cfg.GetCustomize() == nil {
		return nil, false
	}
	// Get custom configuration
	config := cfg.GetCustomize()

	// Create Optimize configuration object
	// Try to decode Optimize configuration from config.Value
	optimizeConfig, err := customize.NewFromStruct[optimizev1.Optimize](config)
	if err != nil {
		helper.Errorf("failed to unmarshal optimize config: %v", err)
		optimizeConfig = defaultOptimize
	}

	return newOptimizer(optimizeConfig), true
}

// defaultOptimize is the default configuration for the OptimizeServer.
var defaultOptimize = &optimizev1.Optimize{
	// Let's start with the legendary 2 seconds.
	Min:      2,
	Max:      30,
	Interval: durationpb.New(time.Hour * 24),
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
func newOptimizer(config *optimizev1.Optimize) middleware.Middleware {
	if config == nil || config.Max == 0 {
		return func(handler middleware.Handler) middleware.Handler {
			return handler
		}
	}

	sleepTime := atomic.Int64{}
	sleepTime.Store(config.Min)
	if config.Min != config.Max {
		go func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			tt := time.NewTicker(config.Interval.AsDuration())
			defer tt.Stop()
			for {
				select {
				case <-tt.C:
					if sleepTime.Load() >= config.Max {
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
