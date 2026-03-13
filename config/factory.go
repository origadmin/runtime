package config

import (
	"fmt"
	"sort"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	sourcev1 "github.com/origadmin/runtime/api/gen/go/config/source/v1"
	"github.com/origadmin/runtime/contracts/builder"
	"github.com/origadmin/runtime/contracts/options"
	internalfactory "github.com/origadmin/runtime/helpers/builder"
	"github.com/origadmin/runtime/log"
)

// Builder is the builder implementation for configurations. It is exported to allow
// for creating independent instances for testing or special use cases, while most
// users will interact with it via the package-level functions that use the default
// global instance.
type Builder struct {
	builder.Registry[SourceFactory]
}

// NewBuilder creates and returns a new, independent Builder instance.
func NewBuilder() *Builder {
	return &Builder{
		Registry: internalfactory.New[SourceFactory](),
	}
}

// NewConfig creates a new configuration object that conforms to the KConfig interface.
func (b *Builder) NewConfig(srcs *sourcev1.Sources, opts ...options.Option) (KConfig, error) {
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
	sort.SliceStable(sourceConfigs, func(i, j int) bool {
		return sourceConfigs[i].GetPriority() < sourceConfigs[j].GetPriority()
	})

	for _, src := range sourceConfigs {
		f, ok := b.Get(src.Type)
		if !ok {
			return nil, fmt.Errorf("unknown config source type: %s", src.Type)
		}
		source, err := f.NewSource(src, opts...)
		if err != nil {
			return nil, err
		}
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

	// Create the underlying Kratos config directly
	kc := kratosconfig.New(fromOptions.ConfigOptions...)

	return kc, nil
}

// defaultBuilder is the package-level global singleton builder instance.
var defaultBuilder = NewBuilder()

// New is a publicly exposed package-level function for creating config instances.
func New(srcs *sourcev1.Sources, opts ...options.Option) (KConfig, error) {
	return defaultBuilder.NewConfig(srcs, opts...)
}

// RegisterSourceFactory is a publicly exposed package-level function for registering a SourceFactory.
func RegisterSourceFactory(name string, factory SourceFactory) {
	defaultBuilder.Register(name, factory)
}

// GetSourceFactory is a publicly exposed package-level function for retrieving a SourceFactory.
func GetSourceFactory(name string) (SourceFactory, bool) {
	return defaultBuilder.Get(name)
}

// getDefaultPriorityForSourceType returns a default priority based on the source type.
func getDefaultPriorityForSourceType(sourceType string) int32 {
	switch SourceType(sourceType) {
	case SourceTypeEnv:
		return 900 // Environment variables
	case SourceTypeFile:
		return 600 // Main application config file
	default:
		return 100 // Lowest default priority
	}
}

type SourceType string

const (
	SourceTypeFile SourceType = "file"
	SourceTypeEnv  SourceType = "env"
)
