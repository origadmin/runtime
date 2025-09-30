/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config provides adapters for Kratos config types and functions.
package config

import (
	_ "github.com/go-kratos/kratos/v2/config"
	kratosconfig "github.com/go-kratos/kratos/v2/config"
)

//go:generate adptool ./kratos.go
//go:adapter:package github.com/go-kratos/kratos/v2/config
//go:adapter:package:type *
//go:adapter:package:type:prefix K
//go:adapter:package:func *
//go:adapter:package:func:prefix K
//go:adapter:package:func New
//go:adapter:package:func:rename NewKkc

// --- Adapter Layer ---

// adapter implements the basic interfaces.kc by wrapping a kratos.kc instance.
// Its sole responsibility is to adapt the Kratos config to our internal interface.
type adapter struct {
	kc kratosconfig.Config
}

// Load implements the interfaces.kc interface.
func (a *adapter) Load() error {
	return a.kc.Load()
}

// Decode implements the interfaces.kc interface.
func (a *adapter) Decode(key string, value any) error {
	if key == "" {
		return a.kc.Scan(value)
	}
	return a.kc.Value(key).Scan(value)
}

// Raw implements the interfaces.kc interface.
func (a *adapter) Raw() any {
	return a.kc
}

// Close implements the interfaces.kc interface.
func (a *adapter) Close() error {
	return a.kc.Close()
}
