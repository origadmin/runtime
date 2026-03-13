/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package log

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	"github.com/origadmin/runtime/contracts"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/helpers/comp"
	"github.com/origadmin/runtime/helpers/configutil"
)

// Resolve resolves the logger configuration.
func Resolve(ctx context.Context, source any, opts *component.LoadOptions) (*component.ModuleConfig, error) {
	if c, ok := source.(contracts.LoggerConfig); ok {
		logger := c.GetLogger()
		if logger == nil {
			return nil, nil
		}
		// Priority: Name -> Type
		name := configutil.ExtractName(logger)
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
		log.Context(ctx).Warnf("logger: no config found, using default logger")
		return DefaultLogger, nil
	}
	return NewLogger(cfg), nil
}
