package runtime

import (
	"context"
	"errors"

	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/log"
)

var (
	// DefaultLoggerExtractor extracts logger config from the root.
	DefaultLoggerExtractor component.Extractor = func(root any) (*component.ModuleConfig, error) {
		if p, ok := root.(component.LoggerConfig); ok && p.GetLogger() != nil {
			l := p.GetLogger()
			return &component.ModuleConfig{
				Entries: []component.ConfigEntry{{Name: "logger", Value: l}},
				Active:  "logger",
			}, nil
		}
		return nil, nil
	}

	// DefaultLoggerProvider creates a logger instance.
	DefaultLoggerProvider component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		var cfg loggerv1.Logger
		if err := h.BindConfig(&cfg); err != nil {
			return log.DefaultLogger, nil
		}
		return log.NewLogger(&cfg), nil
	}

	// DefaultRegistryExtractor extracts registry config from the root.
	DefaultRegistryExtractor component.Extractor = func(root any) (*component.ModuleConfig, error) {
		if p, ok := root.(component.RegistryConfig); ok && p.GetDiscoveries() != nil {
			raw := p.GetDiscoveries()
			var entries []component.ConfigEntry
			for _, c := range raw.Configs {
				entries = append(entries, component.ConfigEntry{Name: c.Name, Value: c})
			}
			return &component.ModuleConfig{Entries: entries, Active: raw.GetActive()}, nil
		}
		return nil, nil
	}

	// DefaultRegistryProvider is a placeholder for registry factory.
	DefaultRegistryProvider component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		return nil, errors.New("registry provider not fully implemented")
	}
)
