/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package log

import (
	"context"

	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	"github.com/origadmin/runtime/contracts"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/helpers/comp"
)

// Resolve resolves the logger configuration.
func Resolve(source any, _ component.Category) (*component.ModuleConfig, error) {
	if c, ok := source.(contracts.LoggerConfig); ok {
		logger := c.GetLogger()
		if logger == nil {
			return nil, nil
		}
		// Priority: Name -> Type
		name := comp.ExtractName(logger)
		if name == "" {
			name = "logger"
		}
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{{Name: name, Value: logger}},
			Active:  name,
		}, nil
	}
	return nil, nil
}

// DefaultProvider is the engine-compatible provider for logger components.
var DefaultProvider component.Provider = func(ctx context.Context, h component.Handle) (any, error) {
	cfg, err := comp.AsConfig[loggerv1.Logger](h)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return DefaultLogger, nil
	}
	return NewLogger(cfg), nil
}
