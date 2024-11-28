/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package bootstrap

import (
	"path/filepath"
	"time"
)

// Constants for default paths and environment
const (
	DefaultConfigPath = "configs/config.toml"
	DefaultEnv        = "dev"
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
	// set workdir to current directory if not set
	if b.WorkDir == "" {
		b.WorkDir = DefaultWorkDir
	}
	// if the workdir is not absolute path, make it absolute
	b.WorkDir, _ = filepath.Abs(b.WorkDir)

	// return the workdir, if the config path is empty
	if b.ConfigPath == "" {
		return b.WorkDir
	}

	// if the config path is absolute path, return it
	if filepath.IsAbs(b.ConfigPath) {
		return b.ConfigPath
	}

	// point the path to the `workdir path/config path`, and make it absolute
	path := filepath.Join(b.WorkDir, b.ConfigPath)
	path, _ = filepath.Abs(path)
	return path
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
