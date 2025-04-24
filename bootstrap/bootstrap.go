/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package bootstrap is a package that provides the bootstrap information for the service.
package bootstrap

import (
	"path/filepath"
	"time"
)

// Constants for default paths and environment
const (
	DefaultConfigPath = "configs/config.toml"
	DefaultEnv        = "release"
	DefaultWorkDir    = "."
)

// Bootstrap struct to hold bootstrap information
type Bootstrap struct {
	Flags      Flags
	WorkDir    string
	ConfigPath string
	Env        string
	Daemon     bool
}

var (
	buildEnv = DefaultEnv
)

// SetFlags sets the flags for the bootstrap
func (b *Bootstrap) SetFlags(name, version string) {
	b.Flags.Version = version
	b.Flags.ServiceName = name
}

// ServiceID returns the service ID
func (b *Bootstrap) ServiceID() string {
	return b.Flags.ServiceID()
}

// ID returns the ID
func (b *Bootstrap) ID() string {
	return b.Flags.ID
}

// Version returns the version
func (b *Bootstrap) Version() string {
	return b.Flags.Version
}

// ServiceName returns the service name
func (b *Bootstrap) ServiceName() string {
	return b.Flags.ServiceName
}

// StartTime returns the start time
func (b *Bootstrap) StartTime() time.Time {
	return b.Flags.StartTime
}

// Metadata returns the metadata
func (b *Bootstrap) Metadata() map[string]string {
	return b.Flags.Metadata
}

// WorkPath returns the work path
func (b *Bootstrap) WorkPath() string {
	if b.WorkDir == "" {
		b.WorkDir = DefaultWorkDir
	}
	b.WorkDir = absPath(b.WorkDir)

	if b.ConfigPath == "" {
		return b.WorkDir
	}

	configPath := b.ConfigPath
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(b.WorkDir, configPath)
	}
	return absPath(configPath)
}

func absPath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	if abs, err := filepath.Abs(p); err == nil {
		return abs
	}
	return p
}

// DefaultBootstrap returns a default bootstrap
func DefaultBootstrap() *Bootstrap {
	return &Bootstrap{
		WorkDir:    DefaultWorkDir,
		ConfigPath: DefaultConfigPath,
		Env:        DefaultEnv,
		Daemon:     false,
		Flags:      DefaultFlags(),
	}
}

// New returns a new bootstrap
func New(dir, path string) *Bootstrap {
	return &Bootstrap{
		WorkDir:    dir,
		ConfigPath: path,
		Env:        DefaultEnv,
		Daemon:     false,
		Flags:      DefaultFlags(),
	}
}

func WithFlags(name string, version string) *Bootstrap {
	return &Bootstrap{
		WorkDir:    DefaultWorkDir,
		ConfigPath: DefaultConfigPath,
		Env:        DefaultEnv,
		Daemon:     false,
		Flags:      NewFlags(name, version),
	}
}
