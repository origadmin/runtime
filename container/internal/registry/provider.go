package registry

import (
	"cmp"
	"errors"
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/log"

	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/registry"
)

// Provider implements interfaces.RegistryProvider
type Provider struct {
	config                 *discoveryv1.Discoveries
	logger                 *log.Helper
	opts                   []options.Option
	discoveries            map[string]registry.KDiscovery
	registrars             map[string]registry.KRegistrar
	onceDiscoveries        sync.Once
	onceRegistrars         sync.Once
	onceDefaultRegistrar   sync.Once
	cachedDefaultRegistrar registry.KRegistrar
}

func (p *Provider) RegisterDiscovery(name string, discovery registry.KDiscovery) {
	//TODO implement me
	panic("implement me")
}

// NewProvider creates a new Provider.
func NewProvider(logger log.Logger, opts []options.Option) *Provider {
	return &Provider{
		logger:      log.NewHelper(logger),
		opts:        opts,
		discoveries: make(map[string]registry.KDiscovery),
		registrars:  make(map[string]registry.KRegistrar),
	}
}

// SetConfig sets the discovery configurations for the provider.
func (p *Provider) SetConfig(cfg *discoveryv1.Discoveries) *Provider {
	p.config = cfg
	return p
}

// Discoveries implements interfaces.RegistryProvider.
func (p *Provider) Discoveries() (map[string]registry.KDiscovery, error) {
	var allErrors error
	p.onceDiscoveries.Do(func() {
		if p.config == nil || len(p.config.GetConfigs()) == 0 {
			p.logger.Infow("msg", "no discovery configurations found")
			return
		}

		for _, cfg := range p.config.GetConfigs() {
			name := cmp.Or(cfg.Name, cfg.Type)
			d, err := registry.NewDiscovery(cfg, p.opts...)
			if err != nil {
				p.logger.Errorf("failed to create discovery '%s': %v", name, err)
				allErrors = errors.Join(allErrors, fmt.Errorf("failed to create discovery '%s': %w", name, err))
				continue
			}
			p.discoveries[name] = d
		}
	})
	return p.discoveries, allErrors
}

// Discovery implements interfaces.RegistryProvider.
func (p *Provider) Discovery(name string) (registry.KDiscovery, error) {
	discoveries, err := p.Discoveries()
	if err != nil {
		return nil, err
	}
	d, ok := discoveries[name]
	if !ok {
		return nil, fmt.Errorf("discovery '%s' not found", name)
	}
	return d, nil
}

// Registrars implements interfaces.RegistryProvider.
func (p *Provider) Registrars() (map[string]registry.KRegistrar, error) {
	var allErrors error
	p.onceRegistrars.Do(func() {
		if p.config == nil || len(p.config.GetConfigs()) == 0 {
			p.logger.Infow("msg", "no registrar configurations found in discoveries config")
			return
		}

		for _, cfg := range p.config.GetConfigs() {
			name := cmp.Or(cfg.Name, cfg.Type)
			reg, err := registry.NewRegistrar(cfg, p.opts...)
			if err != nil {
				p.logger.Errorf("failed to create registrar '%s': %v", name, err)
				allErrors = errors.Join(allErrors, fmt.Errorf("failed to create registrar '%s': %w", name, err))
				continue
			}
			p.registrars[name] = reg
		}
	})
	return p.registrars, allErrors
}

// Registrar implements interfaces.RegistryProvider.
func (p *Provider) Registrar(name string) (registry.KRegistrar, error) {
	registrars, err := p.Registrars()
	if err != nil {
		return nil, err
	}
	reg, ok := registrars[name]
	if !ok {
		return nil, fmt.Errorf("registrar '%s' not found", name)
	}
	return reg, nil
}

// DefaultRegistrar implements interfaces.RegistryProvider.
func (p *Provider) DefaultRegistrar() (registry.KRegistrar, error) {
	var currentErr error
	p.onceDefaultRegistrar.Do(func() {
		if p.config == nil {
			currentErr = errors.New("default registrar not configured")
			p.logger.Infow("msg", currentErr.Error())
			return
		}

		defaultRegistrarName := p.config.GetDefault()
		if defaultRegistrarName == "" {
			currentErr = errors.New("default registrar not configured")
			p.logger.Infow("msg", currentErr.Error())
			return
		}

		// Ensure registrars are loaded
		registrars, err := p.Registrars()
		if err != nil {
			currentErr = errors.Join(currentErr, fmt.Errorf("failed to load registrars for default: %w", err))
			p.logger.Errorf(currentErr.Error())
			return
		}

		reg, ok := registrars[defaultRegistrarName]
		if !ok {
			currentErr = errors.Join(currentErr, fmt.Errorf("default registrar '%s' not found", defaultRegistrarName))
			p.logger.Errorf(currentErr.Error())
			return
		}
		p.cachedDefaultRegistrar = reg
	})
	return p.cachedDefaultRegistrar, currentErr
}
