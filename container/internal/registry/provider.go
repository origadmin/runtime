package registry

import (
	"errors"
	"fmt"
	"sync"

	"github.com/goexts/generic/cmp"
	"github.com/goexts/generic/maps"

	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	runtimelog "github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/registry"
)

// Provider manages the lifecycle of discovery and registrar instances.
// It uses lazy-loading with sync.Once to ensure instances are created only when needed and in a concurrency-safe manner.
type Provider struct {
	mu               sync.RWMutex
	logger           *runtimelog.Helper
	discoveries      map[string]registry.KDiscovery
	registrars       map[string]registry.KRegistrar
	config           *discoveryv1.Discoveries
	opts             []options.Option
	discoveriesOnce  sync.Once
	registrarsOnce   sync.Once
	discoveriesErr   error
	registrarsErr    error
	defaultRegistrar string
}

// NewProvider creates a new, uninitialized Provider instance.
func NewProvider(logger runtimelog.Logger) *Provider {
	return &Provider{
		logger:      runtimelog.NewHelper(logger),
		discoveries: make(map[string]registry.KDiscovery),
		registrars:  make(map[string]registry.KRegistrar),
	}
}

// Initialize configures the provider with the necessary configuration and options.
// This method is lightweight and only stores the provided values. It does not perform any instance creation.
func (p *Provider) Initialize(cfg *discoveryv1.Discoveries, opts ...options.Option) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.config = cfg
	p.defaultRegistrar = cmp.Or(cfg.GetActive(), cfg.GetDefault())
	p.opts = opts
}

// RegisterDiscovery allows for manual registration of a discovery instance.
func (p *Provider) RegisterDiscovery(name string, discovery registry.KDiscovery) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.discoveries[name]; ok {
		p.logger.Warnf("discovery '%s' is being overwritten by manual registration", name)
	}
	p.discoveries[name] = discovery
}

// Discoveries returns a map of all available discovery clients.
// On the first call, it lazily creates and caches instances based on the configuration.
// Any errors during creation are captured and returned on every subsequent call.
func (p *Provider) Discoveries() (map[string]registry.KDiscovery, error) {
	p.discoveriesOnce.Do(func() {
		p.mu.Lock()
		defer p.mu.Unlock()

		if p.config == nil {
			return
		}
		var allErrors error
		for _, cfg := range p.config.GetConfigs() {
			name := cmp.Or(cfg.GetName(), cfg.GetType())
			if name == "" {
				continue
			}
			if _, exists := p.discoveries[name]; exists {
				continue
			}
			d, err := registry.NewDiscovery(cfg, p.opts...)
			if err != nil {
				p.logger.Errorf("failed to create discovery '%s': %v", name, err)
				allErrors = errors.Join(allErrors, err)
				continue
			}
			p.discoveries[name] = d
		}
		p.discoveriesErr = allErrors
	})

	p.mu.RLock()
	defer p.mu.RUnlock()
	return maps.Clone(p.discoveries), p.discoveriesErr
}

// Discovery returns a single discovery client by name.
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

// RegisterRegistrar allows for manual registration of a registrar instance.
func (p *Provider) RegisterRegistrar(name string, registrar registry.KRegistrar) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.registrars[name]; ok {
		p.logger.Warnf("registrar '%s' is being overwritten by manual registration", name)
	}
	p.registrars[name] = registrar
}

// Registrars returns a map of all available service registrars.
// It follows the same lazy-loading and caching logic as Discoveries.
func (p *Provider) Registrars() (map[string]registry.KRegistrar, error) {
	p.registrarsOnce.Do(func() {
		p.mu.Lock()
		defer p.mu.Unlock()

		if p.config == nil {
			return
		}
		var allErrors error
		for _, cfg := range p.config.GetConfigs() {
			name := cmp.Or(cfg.GetName(), cfg.GetType())
			if name == "" {
				continue
			}
			if _, exists := p.registrars[name]; exists {
				continue
			}
			reg, err := registry.NewRegistrar(cfg, p.opts...)
			if err != nil {
				p.logger.Errorf("failed to create registrar '%s': %v", name, err)
				allErrors = errors.Join(allErrors, err)
				continue
			}
			p.registrars[name] = reg
		}
		p.registrarsErr = allErrors
	})

	p.mu.RLock()
	defer p.mu.RUnlock()
	return maps.Clone(p.registrars), p.registrarsErr
}

// Registrar returns a single service registrar by name.
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

// DefaultRegistrar returns the default service registrar.
func (p *Provider) DefaultRegistrar(globalDefaultName string) (registry.KRegistrar, error) {
	registrars, err := p.Registrars()
	if err != nil {
		return nil, err
	}
	if len(registrars) == 0 {
		return nil, errors.New("no registrars available")
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	var prioritizedNames []string
	if p.defaultRegistrar != "" {
		prioritizedNames = append(prioritizedNames, p.defaultRegistrar)
	}
	if globalDefaultName != "" {
		prioritizedNames = append(prioritizedNames, globalDefaultName)
	}
	prioritizedNames = append(prioritizedNames, interfaces.GlobalDefaultKey)
	defaultName, ok := maps.FirstKey(registrars, prioritizedNames...)
	if ok {
		p.logger.Debugf("resolved default registrar to '%s'", defaultName)
		return registrars[defaultName], nil
	}

	if len(registrars) == 1 {
		for name, comp := range registrars {
			p.logger.Debugf("no specific default found, falling back to the first available registrar: '%s'", name)
			return comp, nil
		}
	}

	return nil, errors.New("no default registrar could be determined")
}
