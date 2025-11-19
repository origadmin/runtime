/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package config

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"

	sourcev1 "github.com/origadmin/runtime/api/gen/go/config/source/v1"
	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces/options"
)

const Module = "config"

var (
	ErrInvalidConfigType = runtimeerrors.NewStructured(Module, "invalid config type")
)

// SourceFunc is a function type that adapts a function to the SourceFactory interface.
// This allows registering a simple function as a factory, avoiding the need for a struct.
type SourceFunc func(*sourcev1.SourceConfig, ...options.Option) (kratosconfig.Source, error)

// NewSource makes SourceFunc implement the SourceFactory interface.
// The function itself becomes the factory method.
func (c SourceFunc) NewSource(config *sourcev1.SourceConfig, options ...options.Option) (kratosconfig.Source, error) {
	return c(config, options...)
}

// SourceFactory is the interface for creating configuration sources.
// It defines a single method, NewSource, which creates a new config source
// based on the provided configuration and options.
type SourceFactory interface {
	// NewSource creates a new config source.
	NewSource(*sourcev1.SourceConfig, ...options.Option) (kratosconfig.Source, error)
}

// Load loads configuration from the specified file path and scans it into the target struct.
// It returns the Kratos config instance, which should be closed by the caller when no longer needed.
func Load(configPath string, target interface{}) (kratosconfig.Config, error) {
	c := kratosconfig.New(
		kratosconfig.WithSource(
			file.NewSource(configPath),
		),
	)

	if err := c.Load(); err != nil {
		// Ensure config is closed on load error to prevent resource leaks
		c.Close()
		return nil, runtimeerrors.WrapStructured(err, Module, "failed to load config from %s", configPath).WithCaller()
	}

	if err := c.Scan(target); err != nil {
		// Ensure config is closed on scan error
		c.Close()
		return nil, runtimeerrors.WrapStructured(err, Module, "failed to scan config into target").WithCaller()
	}

	return c, nil
}
