/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package config

import (
	"fmt" // Add fmt import for error formatting

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file" // Add file config source import

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces"
)

const Type = "config"

// NewConfig creates a new config instance.
func NewConfig(cfg *configv1.SourceConfig, opts ...interfaces.Option) (kratosconfig.Config, error) {
	return defaultConfigFactory.NewConfig(cfg, opts...)
}

// Register registers a config factory.
func Register(name string, factory interfaces.ConfigFactory) {
	defaultConfigFactory.Register(name, factory)
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
		return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
	}

	if err := c.Scan(target); err != nil {
		// Ensure config is closed on scan error
		c.Close()
		return nil, fmt.Errorf("failed to scan config into target: %w", err)
	}

	return c, nil
}
