package runtime

import (
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/google/uuid"

	"github.com/origadmin/runtime/interfaces"
)

// AppInfo represents the application's static, immutable identity information.
// It includes essential metadata such as the application's name, version, environment,
// and instance ID. This information is determined at startup and remains constant
// throughout the application's lifecycle.
type AppInfo interfaces.AppInfo

// NewAppInfo creates a new AppInfo instance with default values for ID, StartTime, and Metadata.
// It requires the application's name, version, and environment.
func NewAppInfo(name, version, env string) *AppInfo {
	return &AppInfo{
		Name:      name,
		Version:   version,
		ID:        uuid.New().String(),
		Env:       env,
		StartTime: time.Now(),
		Metadata:  make(map[string]string),
	}
}

// WithMetadata returns a new AppInfo with the provided key-value pair added to the metadata.
// This method supports fluent chaining, e.g., appInfo.WithMetadata(...).WithMetadata(...)
func (a AppInfo) WithMetadata(key, value string) AppInfo {
	if a.Metadata == nil {
		a.Metadata = make(map[string]string)
	}
	a.Metadata[key] = value
	return a
}

// String implements the fmt.Stringer interface for easy logging and identification.
// It returns a string in the format "name-version".
func (a AppInfo) String() string {
	return fmt.Sprintf("%s-%s", a.Name, a.Version)
}

// Options returns a slice of kratos.Option with the application's identity
// fields, suitable for passing to kratos.New().
func (a AppInfo) Options() []kratos.Option {
	// Ensure metadata for env is always present
	if a.Metadata == nil {
		a.Metadata = make(map[string]string)
	}
	a.Metadata["env"] = a.Env

	return []kratos.Option{
		kratos.ID(a.ID),
		kratos.Name(a.Name),
		kratos.Version(a.Version),
		kratos.Metadata(a.Metadata),
	}
}

// GetID returns the unique identifier of the application instance.
func (a AppInfo) GetID() string {
	return a.ID
}

// GetName returns the name of the application.
func (a AppInfo) GetName() string {
	return a.Name
}

// GetVersion returns the version of the application.
func (a AppInfo) GetVersion() string {
	return a.Version
}

// GetEnv returns the environment the application is running in.
func (a AppInfo) GetEnv() string {
	return a.Env
}

// GetStartTime returns the time the application was started.
func (a AppInfo) GetStartTime() time.Time {
	return a.StartTime
}

// GetUptime returns the duration since the application was started.
func (a AppInfo) GetUptime() time.Duration {
	return time.Since(a.StartTime)
}

// IsValid checks if the AppInfo instance contains essential, non-empty identification fields.
func (a AppInfo) IsValid() bool {
	return a.ID != "" && a.Name != "" && a.Version != ""
}
