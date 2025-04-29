/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package bootstrap is a package that provides the bootstrap information for the service.
package bootstrap

import (
	"path/filepath"
	"time"

	"github.com/origadmin/toolkits/errors"
)

// Constants for default paths and environment
const (
	DefaultConfigPath = "configs/config.toml"
	DefaultEnv        = "release"
	DefaultWorkDir    = "."
)

// Bootstrap struct to hold bootstrap information
type Bootstrap struct {
	daemon      bool
	env         string
	workDir     string
	configPath  string
	version     string
	startTime   time.Time
	metadata    map[string]string
	serviceID   string
	serviceName string
}

func (b *Bootstrap) ConfigFilePath() string {
	if b.workDir == "" {
		return absPath(b.configPath)
	}
	workDir := absPath(b.workDir)
	if b.configPath == "" {
		return workDir
	}

	configPath := b.ConfigPath()
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(workDir, configPath)
	}
	return absPath(configPath)
}

func (b *Bootstrap) ConfigPath() string {
	return b.configPath
}

func (b *Bootstrap) Version() string {
	return b.version
}

func (b *Bootstrap) StartTime() time.Time {
	return b.startTime
}

func (b *Bootstrap) Metadata() map[string]string {
	return b.metadata
}

func (b *Bootstrap) ServiceID() string {
	return b.serviceID
}

func (b *Bootstrap) ServiceName() string {
	return b.serviceName
}

func (b *Bootstrap) Daemon() bool {
	return b.daemon
}

func (b *Bootstrap) WorkDir() string {
	return b.workDir
}

func (b *Bootstrap) SetWorkDir(workDir string) {
	b.workDir = workDir
}

func (b *Bootstrap) SetDaemon(daemon bool) {
	b.daemon = daemon
}

func (b *Bootstrap) SetConfigPath(configPath string) {
	b.configPath = configPath
}

func (b *Bootstrap) SetVersion(version string) {
	b.version = version
}

func (b *Bootstrap) SetStartTime(startTime time.Time) {
	b.startTime = startTime
}

func (b *Bootstrap) SetMetadata(metadata map[string]string) {
	b.metadata = metadata
}

func (b *Bootstrap) SetServiceID(serviceID string) {
	b.serviceID = serviceID
}

func (b *Bootstrap) SetServiceName(serviceName string) {
	b.serviceName = serviceName
}

var (
	buildEnv = DefaultEnv
)

func (b *Bootstrap) SetEnv(env string) error {
	if env != "debug" && env != "release" {
		return errors.New("invalid env value")
	}
	b.env = env
	return nil
}

func (b *Bootstrap) Env() string {
	return b.env
}

func (b *Bootstrap) SetServiceInfo(name, version string) {
	b.serviceName = name
	b.version = version
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

// New returns a new bootstrap
func New() *Bootstrap {
	return &Bootstrap{
		//WorkDir:     DefaultWorkDir,
		//ConfigPath:  DefaultConfigPath,
		env:       buildEnv,
		serviceID: RandomID(),
		//version:     version,
		//serviceName: name,
		startTime: time.Now(),
		metadata:  make(map[string]string),
	}
}

func WithFlags(name string, version string) *Bootstrap {
	bs := New()
	bs.serviceName = name
	bs.version = version
	return bs
}
