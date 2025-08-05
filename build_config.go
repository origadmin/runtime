/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces" // Add interfaces import
)

// LoadConfig loads the config file from the given path
func LoadConfig(path string, v any, ss ...interfaces.Option) error {
	sourceConfig, err := bootstrap.LoadSourceConfig(path)
	if err != nil {
		return err
	}
	runtimeConfig, err := NewConfig(sourceConfig, ss...)
	if err != nil {
		return err
	}
	if err := runtimeConfig.Load(); err != nil {
		return err
	}
	if err := runtimeConfig.Scan(v); err != nil {
		return err
	}
	return nil
}

func LoadConfigFromBootstrap(bs *bootstrap.Bootstrap, v any, ss ...interfaces.Option) error {
	sourceConfig, err := bootstrap.LoadSourceConfigFromBootstrap(bs)
	if err != nil {
		return err
	}
	runtimeConfig, err := NewConfig(sourceConfig, ss...)
	if err != nil {
		return err
	}
	if err := runtimeConfig.Load(); err != nil {
		return err
	}
	if err := runtimeConfig.Scan(v); err != nil {
		return err
	}
	return nil
}
