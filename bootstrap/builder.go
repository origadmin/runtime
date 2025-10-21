package bootstrap

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/runtime/bootstrap/internal/container"
	"github.com/origadmin/runtime/interfaces"
)

// buildContainer creates the component container and initializes all registered components.
func buildContainer(sc interfaces.StructuredConfig, factories map[string]interfaces.ComponentFactory, opts ...Option) (interfaces.Container, log.Logger, error) {
	// 1. Create the component provider implementation.
	builder := container.NewBuilder(factories).WithConfig(sc)

	// 2. Initialize core components by consuming the config.
	c, err := builder.Build(opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize components: %w", err)
	}

	return c, builder.Logger(), nil
}