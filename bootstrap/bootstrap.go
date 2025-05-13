// Package bootstrap is a package that provides the bootstrap information for the service.
package bootstrap

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Constants for default paths and environment
const (
	EnvRelease        = "release"
	EnvDebug          = "debug"
	DefaultConfigPath = "configs/config.toml"
	defaultEnv        = EnvRelease
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

func (b *Bootstrap) SetDaemon(daemon bool) {
	b.daemon = daemon
}

func (b *Bootstrap) SetWorkDir(workDir string) {
	b.workDir = workDir
}

func (b *Bootstrap) SetConfigPath(configPath string) {
	b.configPath = configPath
}

func (b *Bootstrap) SetPath(dir, configPath string) {
	b.workDir = dir
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
	buildEnv = defaultEnv
)

func (b *Bootstrap) SetEnv(env string) {
	b.env = env
}

func (b *Bootstrap) Env() string {
	return b.env
}

func (b *Bootstrap) IsDebug() bool {
	return b.env == EnvDebug
}

func (b *Bootstrap) SetServiceInfo(name, version string) {
	b.serviceName = name
	b.version = version
}

func (b *Bootstrap) ServiceInfo() ServiceInfo {
	return ServiceInfo{
		Name:      b.serviceName,
		Version:   b.version,
		ID:        b.serviceID,
		Metadata:  b.metadata,
		StartTime: b.startTime,
	}
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
		workDir:     DefaultWorkDir,
		configPath:  DefaultConfigPath,
		version:     DefaultVersion,
		serviceName: DefaultServiceName,
		env:         buildEnv,
		serviceID:   RandomID(),
		startTime:   time.Now(),
		metadata:    make(map[string]string),
	}
}

func WithFlags(name string, version string) *Bootstrap {
	bs := New()
	bs.serviceName = name
	bs.version = version
	return bs
}

func RandomID() string {
	id, err := os.Hostname()
	if err != nil {
		id = "unknown"
	}

	b := make([]byte, 4)
	if _, err := rand.Read(b); err == nil {
		return fmt.Sprintf("%s.%x", id, b)
	}
	return id + "." + RandomSuffix
}
