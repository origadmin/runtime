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

// newAppInfo is the internal constructor that returns a concrete *appInfo struct.
// It is the core logic for creating and configuring an appInfo instance.
func newAppInfo(name, version string) *appInfo {
	a := &appInfo{
		name:     name,
		version:  version,
		id:       uuid.New().String(),
		env:      "dev",
		metadata: make(map[string]string),
	}

	// Set start time at build time if it hasn't been set by an option.
	if a.startTime.IsZero() {
		a.startTime = time.Now()
	}

	return a
}

// NewAppInfo is the public constructor that returns an interfaces.AppInfo.
// It acts as a wrapper around the internal newAppInfo, hiding the concrete
// implementation from the outside world.
func NewAppInfo(name, version string, opts ...AppInfoOption) interfaces.AppInfo {
	a := newAppInfo(name, version)
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// --- Implementation of interfaces.AppInfo ---

func (a *appInfo) ID() string           { return a.id }
func (a *appInfo) Name() string         { return a.name }
func (a *appInfo) Version() string      { return a.version }
func (a *appInfo) Env() string          { return a.env }
func (a *appInfo) StartTime() time.Time { return a.startTime }

// Metadata returns a defensive copy of the metadata map to ensure immutability.
func (a *appInfo) Metadata() map[string]string {
	return maps.Clone(a.metadata)
}

// --- Compile-time checks ---
var _ interfaces.AppInfo = (*appInfo)(nil)
