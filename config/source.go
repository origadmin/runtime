/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package config

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"

	runtimeerrors "github.com/origadmin/runtime/errors"
)

const Module = "config"

var (
	ErrInvalidConfigType = runtimeerrors.NewStructured(Module, "invalid config type")
)

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
