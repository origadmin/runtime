package runtime

import (
	"time"

	"github.com/goexts/generic/maps"
	"github.com/google/uuid"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	"github.com/origadmin/runtime/interfaces"
)

// AppInfo is the concrete, exported implementation that satisfies the AppInfo interface.
type AppInfo struct {
	id        string
	name      string
	version   string
	env       string
	startTime time.Time
	metadata  map[string]string
}

// NewAppInfo is the primary constructor for creating a pre-configured AppInfo instance.
func NewAppInfo(name, version string) *AppInfo {
	return &AppInfo{
		name:      name,
		version:   version,
		id:        uuid.New().String(),
		env:       "dev",
		metadata:  make(map[string]string),
		startTime: time.Now(),
	}
}

// NewAppInfoBuilder returns a new, blank AppInfo instance,
// serving as the entry point for chainable method configurations.
func NewAppInfoBuilder() *AppInfo {
	return &AppInfo{
		id:        uuid.New().String(),
		metadata:  make(map[string]string),
		startTime: time.Now(),
	}
}

// SetID sets the application ID.
func (a *AppInfo) SetID(id string) *AppInfo {
	a.id = id
	return a
}

// SetName sets the application name.
func (a *AppInfo) SetName(name string) *AppInfo {
	a.name = name
	return a
}

// SetVersion sets the application version.
func (a *AppInfo) SetVersion(version string) *AppInfo {
	a.version = version
	return a
}

// SetNameAndVersion sets both the application name and version.
func (a *AppInfo) SetNameAndVersion(name, version string) *AppInfo {
	a.name = name
	a.version = version
	return a
}

// SetEnv sets the application environment.
func (a *AppInfo) SetEnv(env string) *AppInfo {
	a.env = env
	return a
}

// SetStartTime sets the application start time.
func (a *AppInfo) SetStartTime(startTime time.Time) *AppInfo {
	a.startTime = startTime
	return a
}

// AddMetadata adds a single key-value pair to the application's metadata.
func (a *AppInfo) AddMetadata(key, value string) *AppInfo {
	if a.metadata == nil {
		a.metadata = make(map[string]string)
	}
	a.metadata[key] = value
	return a
}

// SetMetadata completely replaces the application's metadata with the provided map.
func (a *AppInfo) SetMetadata(metadata map[string]string) *AppInfo {
	a.metadata = metadata
	return a
}

// Merge combines the current AppInfo with another interfaces.AppInfo.
// Values from the 'other' AppInfo will override the current AppInfo's values
// if they are not empty or zero.
func (a *AppInfo) Merge(other interfaces.AppInfo) {
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

func (a *AppInfo) ID() string           { return a.id }
func (a *AppInfo) Name() string         { return a.name }
func (a *AppInfo) Version() string      { return a.version }
func (a *AppInfo) Env() string          { return a.env }
func (a *AppInfo) StartTime() time.Time { return a.startTime }

// Metadata returns a defensive copy of the metadata map to ensure immutability.
func (a *AppInfo) Metadata() map[string]string {
	return maps.Clone(a.metadata)
}

// mergeAppInfoWithConfig merges application information from a protobuf configuration
// into an existing AppInfo instance. Values from the configuration will override
// existing AppInfo values if they are not empty. A new AppInfo instance is returned.
func mergeAppInfoWithConfig(currentAppInfo *AppInfo, config *appv1.App) *AppInfo {
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
		return &AppInfo{}
	}
	metadata := appConfig.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}

	return &AppInfo{
		name:     appConfig.GetName(),
		version:  appConfig.GetVersion(),
		id:       appConfig.GetId(),
		env:      appConfig.GetEnv(),
		metadata: metadata,
	}
}

// --- Compile-time checks ---
var _ interfaces.AppInfo = (*AppInfo)(nil)
