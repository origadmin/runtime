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

const (
	SourceTypeFile SourceType = "file"
	SourceTypeEnv  SourceType = "env"
	//SourceTypeConsul = "consul"
	//SourceTypeNacos  = "nacos"
	//SourceTypeEtcd   = "etcd"
)

type SourceType string

var (
	// defaultBuilder is the default config factory.
	defaultBuilder = NewBuilder()
)

// sourceFactory is a config factory that implements interfaces.ConfigBuilder.
type sourceFactory struct {
	factory.Registry[SourceFactory]
}

// BuildFunc is a function type that takes a KConfig and a list of Options and returns a Selector and an error.
type BuildFunc func(*sourcev1.SourceConfig, ...options.Option) (kratosconfig.Source, error)

// NewSource is a method that implements the ConfigBuilder interface for ConfigBuildFunc.
func (fn BuildFunc) NewSource(cfg *sourcev1.SourceConfig, opts ...options.Option) (kratosconfig.Source, error) {
	// Call the function with the given KConfig and a list of Options.
	return fn(cfg, opts...)
}

// getDefaultPriorityForSourceType returns a default priority based on the source type.
// Higher values mean higher priority (overrides lower priority configs).
func getDefaultPriorityForSourceType(sourceType string) int32 {
	switch SourceType(sourceType) {
	case SourceTypeEnv:
		return 900 // Environment variables (highest common override)
	case SourceTypeFile:
		return 600 // Main application config file
	// Add other remote types as needed, e.g., "consul", "nacos", "etcd"
	// case SourceTypeConsul:
	// 	return 800 // Remote config service
	default:
		return 100 // Lowest default priority for unknown or base types
	}
}

// NewConfig creates a new configuration object that conforms to the interfaces.Config interface.
// It builds a Kratos config from sources, loads it, and immediately wraps it in an adapter
// to hide the underlying implementation from the rest of the framework.
func (f *sourceFactory) NewConfig(srcs *sourcev1.Sources, opts ...options.Option) (interfaces.Config, error) {
	logger := log.NewHelper(log.FromOptions(opts))
	fromOptions := FromOptions(opts...)
	var sources []kratosconfig.Source

	// Get the list of sources from the protobuf config.
	sourceConfigs := srcs.GetConfigs()

	// --- START: Assign Default Priorities if not set ---
	for _, src := range sourceConfigs {
		if src.GetPriority() == 0 {
			src.Priority = getDefaultPriorityForSourceType(src.GetType())
		}
	}
	// --- END: Assign Default Priorities if not set ---

	// Sort the sources by priority before creating them.
	sort.SliceStable(sourceConfigs, func(i, j int) bool {
		// Sources with lower priority values are loaded first.
		// Sources with higher priority values are loaded later, thus overriding earlier ones.
		return sourceConfigs[i].GetPriority() < sourceConfigs[j].GetPriority()
	})

	for _, src := range sourceConfigs {
		f, ok := f.Get(src.Type)
		if !ok {
			return nil, fmt.Errorf("unknown type: %s", src.Type)
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

func (f *sourceFactory) SyncConfig(cfg *sourcev1.SourceConfig, v any, opts ...options.Option) error {
	// This method is a placeholder. Actual synchronization logic would go here.
	// For now, we'll just return nil or an error if needed.
	return nil
}

// NewBuilder creates a new config factory.
func NewBuilder() Builder {
	return &sourceFactory{
		Registry: internalfactory.New[SourceFactory](),
	}
}
