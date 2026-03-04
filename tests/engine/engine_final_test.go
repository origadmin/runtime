package engine_test

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/engine"
)

type customLogger struct {
	log.Logger
}

func TestCustomRegistryOverriding(t *testing.T) {
	ctx := context.Background()

	// 1. Create App instance
	app := runtime.New("test-app", "1.0.0")
	reg := app.Container()

	// 2. Register Custom Logger (Simulating init() behavior)
	reg.Register(runtime.CategoryLogger, func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		t.Log("Creating custom logger from manual registration")
		return &customLogger{}, nil
	}, engine.WithExtractor(func(root any) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{{Name: "logger", Value: nil}},
			Active:  "logger",
		}, nil
	}), engine.WithPriority(100))

	// 3. Directly load configuration into the container (Injecting, not loading from file)
	if err := reg.Load(ctx, struct{}{}); err != nil {
		t.Fatalf("Container Load failed: %v", err)
	}

	// 4. Verify Override
	l := app.Logger()
	if _, ok := l.(*customLogger); !ok {
		t.Fatalf("Expected customLogger, got %T. User-registered factory should take precedence.", l)
	}

	t.Log("Successfully verified that manual/init registration takes precedence.")
}
