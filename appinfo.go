package runtime

import (
	"time"

	"github.com/goexts/generic/maps"
	"github.com/google/uuid"

	"github.com/origadmin/runtime/interfaces"
)

// appInfo is the concrete, unexported implementation that satisfies both
// the AppInfo and AppInfoBuilder interfaces.
type appInfo struct {
	id        string
	name      string
	version   string
	env       string
	startTime time.Time
	metadata  map[string]string
}

// NewAppInfoBuilder is the public constructor. It returns the AppInfoBuilder
// interface, hiding the concrete implementation.
func NewAppInfoBuilder(name, version string) interfaces.AppInfoBuilder {
	return &appInfo{
		name:     name,
		version:  version,
		id:       uuid.New().String(), // Default ID
		env:      "dev",               // Default Env
		metadata: make(map[string]string),
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

// --- Implementation of interfaces.AppInfoBuilder ---

func (a *appInfo) WithEnv(env string) interfaces.AppInfoBuilder {
	if env != "" {
		a.env = env
	}
	return a
}

func (a *appInfo) WithID(id string) interfaces.AppInfoBuilder {
	if id != "" {
		a.id = id
	}
	return a
}

func (a *appInfo) WithStartTime(startTime time.Time) interfaces.AppInfoBuilder {
	if !startTime.IsZero() {
		a.startTime = startTime
	}
	return a
}

func (a *appInfo) WithMetadata(key, value string) interfaces.AppInfoBuilder {
	if a.metadata == nil {
		a.metadata = make(map[string]string)
	}
	a.metadata[key] = value
	return a
}

func (a *appInfo) Build() interfaces.AppInfo {
	// Set start time at build time if it hasn't been set.
	if a.startTime.IsZero() {
		a.startTime = time.Now()
	}
	// Return a defensive copy to ensure the returned AppInfo is immutable
	// and cannot be affected by further changes to the builder.
	return &appInfo{
		id:        a.id,
		name:      a.name,
		version:   a.version,
		env:       a.env,
		startTime: a.startTime,
		metadata:  a.Metadata(),
	}
}

// --- Compile-time checks ---
var _ interfaces.AppInfo = (*appInfo)(nil)
var _ interfaces.AppInfoBuilder = (*appInfo)(nil)
