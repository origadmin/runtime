package runtime

import (
	"time"

	"github.com/goexts/generic/maps"
	"github.com/google/uuid"

	"github.com/origadmin/runtime/interfaces"
)

// appInfo is the concrete, unexported implementation that satisfies the AppInfo interface.
type appInfo struct {
	id        string
	name      string
	version   string
	env       string
	startTime time.Time
	metadata  map[string]string
}

// AppInfoOption defines a functional option for configuring AppInfo.
type AppInfoOption func(*appInfo)

// NewAppInfo creates a new AppInfo instance using functional options.
func NewAppInfo(name, version string, opts ...AppInfoOption) interfaces.AppInfo {
	// Default values
	a := &appInfo{
		name:     name,
		version:  version,
		id:       uuid.New().String(), // Default ID
		env:      "dev",               // Default Env
		metadata: make(map[string]string),
	}

	// Apply options
	for _, opt := range opts {
		opt(a)
	}

	// Set start time at build time if it hasn't been set.
	if a.startTime.IsZero() {
		a.startTime = time.Now()
	}

	// Return a defensive copy to ensure the returned AppInfo is immutable.
	return &appInfo{
		id:        a.id,
		name:      a.name,
		version:   a.version,
		env:       a.env,
		startTime: a.startTime,
		metadata:  maps.Clone(a.metadata), // Clone metadata to ensure immutability
	}
}

// WithAppInfoEnv sets the environment for the application.
func WithAppInfoEnv(env string) AppInfoOption {
	return func(a *appInfo) {
		if env != "" {
			a.env = env
		}
	}
}

// WithAppInfoID sets a custom instance ID.
func WithAppInfoID(id string) AppInfoOption {
	return func(a *appInfo) {
		if id != "" {
			a.id = id
		}
	}
}

// WithAppInfoStartTime sets a custom start time.
func WithAppInfoStartTime(startTime time.Time) AppInfoOption {
	return func(a *appInfo) {
		if !startTime.IsZero() {
			a.startTime = startTime
		}
	}
}

// WithAppInfoMetadata adds a key-value pair to the application's metadata.
func WithAppInfoMetadata(key, value string) AppInfoOption {
	return func(a *appInfo) {
		if a.metadata == nil {
			a.metadata = make(map[string]string)
		}
		a.metadata[key] = value
	}
}

// --- Implementation of interfaces.AppInfo ---

func (a *appInfo) ID() string           { return a.id }
func (a *appInfo) Name() string         { return a.name }
func (a *appInfo) Version() string      { return a.version }
func (a *appInfo) Env() string          { return a.env }
func (a *appInfo) StartTime() time.Time { return a.startTime }
func (a *appInfo) Metadata() map[string]string {
	return maps.Clone(a.metadata)
}

// --- Compile-time checks ---
var _ interfaces.AppInfo = (*appInfo)(nil)
