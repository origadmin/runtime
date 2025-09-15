package config

import (
	"fmt"
	"sort"

	"github.com/goexts/generic/configure"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	sourcev1 "github.com/origadmin/runtime/api/gen/go/source/v1"
	"github.com/origadmin/runtime/interfaces/factory"
)

var (
	// defaultBuilder is the default config factory.
	defaultBuilder = NewBuilder()
)

// sourceFactory is a config factory that implements interfaces.ConfigBuilder.
type sourceFactory struct {
	factory.Registry[SourceFactory]
}

// RegisterConfigFunc registers a new ConfigBuilder with the given name and function.
func (f *sourceFactory) RegisterConfigFunc(name string, buildFunc BuildFunc) {
	f.Register(name, buildFunc)
}

// BuildFunc is a function type that takes a KConfig and a list of Options and returns a Selector and an error.
type BuildFunc func(*sourcev1.SourceConfig, *Options) (kratosconfig.Source, error)

// NewSource is a method that implements the ConfigBuilder interface for ConfigBuildFunc.
func (fn BuildFunc) NewSource(cfg *sourcev1.SourceConfig, opts *Options) (kratosconfig.Source, error) {
	// Call the function with the given KConfig and a list of Options.
	return fn(cfg, opts)
}

// NewConfig creates a new Selector object based on the given KConfig and options.
func (f *sourceFactory) NewConfig(srcs *sourcev1.Sources, opts ...Option) (kratosconfig.Config, error) {
	options := configure.Apply(&Options{}, opts) // Corrected: Use settings.Apply with a new interfaces.Options{}

	var sources []kratosconfig.Source

	// Get the list of sources from the protobuf config.
	userSources := srcs.GetSources()

	// Sort the sources by priority before creating them.
	sort.SliceStable(userSources, func(i, j int) bool {
		// Sources with lower priority values are loaded first.
		// Sources with higher priority values are loaded later, thus overriding earlier ones.
		return userSources[i].GetPriority() < userSources[j].GetPriority()
	})

	for _, src := range userSources {
		buildFactory, ok := f.Get(src.Type)
		if !ok {
			return nil, fmt.Errorf("unknown type: %s", src.Type)
		}
		source, err := buildFactory.NewSource(src, options)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	options.ConfigOptions = append(options.ConfigOptions, kratosconfig.WithSource(sources...))
	return kratosconfig.New(options.ConfigOptions...), nil
}

func (f *sourceFactory) SyncConfig(cfg *sourcev1.SourceConfig, v any, opts ...Option) error {
	// This method is a placeholder. Actual synchronization logic would go here.
	// For now, we'll just return nil or an error if needed.
	return nil
}

// NewBuilder creates a new config factory.
func NewBuilder() Builder {
	return &sourceFactory{
		Registry: factory.New[SourceFactory](),
	}
}
