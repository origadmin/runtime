package config

import (
	"fmt"
	"sort"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	sourcev1 "github.com/origadmin/runtime/api/gen/go/config/source/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/factory"
	"github.com/origadmin/runtime/interfaces/options"
	internalfactory "github.com/origadmin/runtime/internal/factory"
	"github.com/origadmin/runtime/log"
)

// sourceFactory is the internal builder implementation for configurations.
// It holds a registry of source factories and is responsible for creating
// a final configuration object from multiple sources.
type sourceFactory struct {
	registry factory.Registry[SourceFactory]
}

// New creates a new configuration object that conforms to the interfaces.Config interface.
// It builds a Kratos config from sources, loads it, and immediately wraps it in an adapter
// to hide the underlying implementation from the rest of the framework.
func (sf *sourceFactory) New(srcs *sourcev1.Sources, opts ...options.Option) (interfaces.Config, error) {
	logger := log.NewHelper(log.FromOptions(opts))
	fromOptions := FromOptions(opts...)
	var sources []kratosconfig.Source

	// Get the list of sources from the protobuf config.
	sourceConfigs := srcs.GetConfigs()

	// Assign default priorities if not set.
	for _, src := range sourceConfigs {
		if src.GetPriority() == 0 {
			src.Priority = getDefaultPriorityForSourceType(src.GetType())
		}
	}

	// Sort the sources by priority before creating them.
	// Sources with lower priority values are loaded first.
	// Sources with higher priority values are loaded later, thus overriding earlier ones.
	sort.SliceStable(sourceConfigs, func(i, j int) bool {
		return sourceConfigs[i].GetPriority() < sourceConfigs[j].GetPriority()
	})

	for _, src := range sourceConfigs {
		f, ok := sf.registry.Get(src.Type)
		if !ok {
			return nil, fmt.Errorf("unknown config source type: %s", src.Type)
		}
		source, err := f.NewSource(src, opts...)
		if err != nil {
			return nil, err
		}
		// Defensively check if the factory returned a nil source, which would cause a panic later.
		if source == nil {
			return nil, fmt.Errorf("config source factory for type '%s' returned a nil source", src.Type)
		}
		logger.Infof("Created source: %s with priority: %d", src.Type, src.Priority)
		sources = append(sources, source)
	}
	if fromOptions.Sources != nil {
		sources = append(sources, fromOptions.Sources...)
		logger.Infof("Added %d sources from options", len(fromOptions.Sources))
	}
	fromOptions.ConfigOptions = append(fromOptions.ConfigOptions, kratosconfig.WithSource(sources...))

	// Create the underlying Kratos config
	kc := kratosconfig.New(fromOptions.ConfigOptions...)

	// Wrap it in our adapter and return the interface
	return &adapter{kc: kc}, nil
}

// defaultSourceBuilder is the package-level internal singleton builder instance.
var defaultSourceBuilder = &sourceFactory{
	registry: internalfactory.New[SourceFactory](),
}

// NewConfig is a publicly exposed package-level function for creating config instances.
// It delegates the call to the internal defaultSourceBuilder.
func NewConfig(srcs *sourcev1.Sources, opts ...options.Option) (interfaces.Config, error) {
	return defaultSourceBuilder.New(srcs, opts...)
}

// RegisterSourceFactory is a publicly exposed package-level function for registering a SourceFactory.
// It delegates the call to the internal defaultSourceBuilder.
func RegisterSourceFactory(name string, factory SourceFactory) {
	defaultSourceBuilder.registry.Register(name, factory)
}

// GetSourceFactory is a publicly exposed package-level function for retrieving a SourceFactory.
// It delegates the call to the internal defaultSourceBuilder.
func GetSourceFactory(name string) (SourceFactory, bool) {
	return defaultSourceBuilder.registry.Get(name)
}

// getDefaultPriorityForSourceType returns a default priority based on the source type.
// Higher values mean higher priority (overrides lower priority configs).
func getDefaultPriorityForSourceType(sourceType string) int32 {
	switch SourceType(sourceType) {
	case SourceTypeEnv:
		return 900 // Environment variables (highest common override)
	case SourceTypeFile:
		return 600 // Main application config file
	default:
		return 100 // Lowest default priority for unknown or base types
	}
}

type SourceType string

const (
	SourceTypeFile SourceType = "file"
	SourceTypeEnv  SourceType = "env"
)
