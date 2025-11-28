package runtime

import (
	"time"

	"github.com/goexts/generic/maps"
	"github.com/google/uuid"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
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

// Merge combines the current appInfo with another interfaces.AppInfo.
// Values from the 'other' AppInfo will override the current appInfo's values
// if they are not empty or zero.
func (a *appInfo) Merge(other interfaces.AppInfo) {
	if other == nil {
		return
	}

	if other.ID() != "" {
		a.id = other.ID()
	}
	if other.Name() != "" {
		a.name = other.Name()
	}
	if other.Version() != "" {
		a.version = other.Version()
	}
	if other.Env() != "" {
		a.env = other.Env()
	}
	if other.Metadata() != nil {
		if a.metadata == nil {
			a.metadata = make(map[string]string)
		}
		for k, v := range other.Metadata() {
			a.metadata[k] = v
		}
	}
	// startTime is not merged as it represents the application's actual start time,
	// which should not be overridden by configuration.
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

// mergeAppInfoWithConfig merges application information from a protobuf configuration
// into an existing appInfo instance. Values from the configuration will override
// existing appInfo values if they are not empty. A new appInfo instance is returned.
func mergeAppInfoWithConfig(currentAppInfo *appInfo, config *appv1.App) *appInfo {
	if config == nil {
		return currentAppInfo
	}

	// Create a temporary AppInfo from the config
	configAppInfo := ConvertToAppInfo(config)
	currentAppInfo.Merge(configAppInfo) // Use the new Merge method
	return currentAppInfo
}

// ConvertToAppInfo converts a protobuf App message to an interfaces.AppInfo.
func ConvertToAppInfo(appConfig *appv1.App) interfaces.AppInfo {
	if appConfig == nil {
		return &appInfo{}
	}
	metadata := appConfig.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}

	return &appInfo{
		name:     appConfig.GetName(),
		version:  appConfig.GetVersion(),
		id:       appConfig.GetId(),
		env:      appConfig.GetEnv(),
		metadata: metadata,
	}
}

// --- Compile-time checks ---
var _ interfaces.AppInfo = (*appInfo)(nil)
