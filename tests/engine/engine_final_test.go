package engine_test

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/engine"
	"github.com/origadmin/runtime/engine/metadata"
)

type customLogger struct {
	log.Logger
}

func TestCustomRegistryOverriding(t *testing.T) {
	ctx := context.Background()

	// 1. Initialize App with a custom Registry config option
	app := runtime.New("test-app", "1.0.0", runtime.WithRegistry(func(reg component.Registry) {
		reg.Register(metadata.CategoryLogger, func(root any) (*component.ModuleConfig, error) {
			return &component.ModuleConfig{
				Entries: []component.ConfigEntry{{Name: "logger", Value: nil}},
				Active:  "logger",
			}, nil
		}, func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
			t.Log("Creating custom logger from user registration")
			return &customLogger{}, nil
		}, engine.WithPriority(metadata.PriorityInfrastructure))
	}))

	// 2. Explicit Activation via WarmUp
	if err := app.Container().Init(ctx, struct{}{}); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 3. Verify
	l := app.Logger()
	if _, ok := l.(*customLogger); !ok {
		t.Fatalf("Expected customLogger, got %T.", l)
	}

	t.Log("Successfully verified that user-registered factory takes precedence over DefaultRegister.")
}
