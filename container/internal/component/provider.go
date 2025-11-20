package component

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/runtime/interfaces" // Ensure interfaces is imported
	"github.com/origadmin/runtime/interfaces/options"
	runtimelog "github.com/origadmin/runtime/log"
)

// Provider implements interfaces.ComponentProvider
type Provider struct {
	config                    interfaces.StructuredConfig
	log                       *log.Helper
	opts                      []options.Option
	genericComponentFactories interfaces.GenericComponentFactory // Needed for generic components
	cachedComponents          map[string]interfaces.Component
	onceComponents            sync.Once
}

// NewProvider creates a new Provider.
func NewProvider(logger log.Logger, opts []options.Option, genericFactories interfaces.GenericComponentFactory) *Provider {
	helper := log.NewHelper(logger)
	return &Provider{
		log:                       helper,
		opts:                      opts,
		genericComponentFactories: genericFactories,
		cachedComponents:          make(map[string]interfaces.Component),
	}
}

// SetConfig sets the structured configurations for the provider.
func (p *Provider) SetConfig(cfg interfaces.StructuredConfig) *Provider {
	p.config = cfg
	return p
}

// Components returns all the configured components.
func (p *Provider) Components() (map[string]interfaces.Component, error) {
	var allErrors error
	p.onceComponents.Do(func() {
		componentConfigs := make(map[string]*interfaces.ComponentConfig)
		err := p.config.Decode("component", &componentConfigs) // Using generic decode
		if err != nil {
			p.log.Errorf("failed to decode component configurations: %v", err)
			allErrors = errors.Join(allErrors, fmt.Errorf("failed to decode component configurations: %w", err))
			return
		}
		if componentConfigs == nil || len(componentConfigs) == 0 {
			p.log.Infow("msg", "no component configurations found")
			return
		}

		for name, cfg := range componentConfigs {
			if p.genericComponentFactories == nil {
				allErrors = errors.Join(allErrors, fmt.Errorf("generic component factory is not available"))
				break
			}

			comp, err := p.genericComponentFactories.NewComponent(name, cfg, p.opts...)
			if err != nil {
				p.log.Errorf("failed to create component '%s' of type '%s': %v", name, cfg.Type, err)
				allErrors = errors.Join(allErrors, fmt.Errorf("failed to create component '%s' with type '%s': %w", name, cfg.Type, err))
				continue
			}
			p.cachedComponents[name] = comp
		}
	})
	return p.cachedComponents, allErrors
}

// Component returns a specific component by name.
func (p *Provider) Component(name string) (interfaces.Component, error) {
	s, err := p.Components()
	if err != nil {
		return nil, err
	}
	comp, ok := s[name]
	if !ok {
		return nil, fmt.Errorf("component '%s' not found", name)
	}
	return comp, nil
}
