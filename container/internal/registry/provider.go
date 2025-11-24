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

// Provider implements interfaces.RegistryProvider. It manages the lifecycle of discovery
// and registrar instances, caching them after first creation and allowing for reconfiguration.
// It is safe for concurrent use.
type Provider struct {
	mu                     sync.Mutex
	config                 *discoveryv1.Discoveries
	logger                 *log.Helper
	opts                   []options.Option
	defaultRegistrar       string
	discoveries            map[string]registry.KDiscovery
	registrars             map[string]registry.KRegistrar
	discoveriesInitialized bool
	registrarsInitialized  bool
}

// NewProvider creates a new Provider.
func NewProvider(logger log.Logger) *Provider {
	return &Provider{
		logger: log.NewHelper(logger),
	}
}

// SetConfig updates the provider's configuration. This will clear any previously
// cached instances and cause them to be recreated on the next access, using the new configuration.
func (p *Provider) SetConfig(cfg *discoveryv1.Discoveries, opts ...options.Option) *Provider {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.config = cfg
	p.opts = opts
	p.discoveriesInitialized = false
	p.registrarsInitialized = false
	p.discoveries = make(map[string]registry.KDiscovery)
	p.registrars = make(map[string]registry.KRegistrar)

	return p
}

// RegisterDiscovery allows for manual registration of a discovery instance.
// This instance will be available alongside any instances created from configuration.
func (p *Provider) RegisterDiscovery(name string, discovery registry.KDiscovery) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.discoveries[name] = discovery
}

// Discoveries returns a map of all available discovery clients.
// On first call, it creates instances from the configuration and caches them.
// Subsequent calls return the cached instances unless SetConfig has been called.
func (p *Provider) Discoveries() (map[string]registry.KDiscovery, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.discoveriesInitialized {
		return p.discoveries, nil
	}

	var allErrors error
	if p.config != nil {
		for _, cfg := range p.config.GetConfigs() {
			name := cmp.Or(cfg.Name, cfg.Type)
			if _, exists := p.discoveries[name]; exists {
				p.logger.Warnf("discovery '%s' is already registered, skipping config-based creation", name)
				continue
			}
			d, err := registry.NewDiscovery(cfg, p.opts...)
			if err != nil {
				p.logger.Errorf("failed to create discovery '%s': %v", name, err)
				allErrors = errors.Join(allErrors, fmt.Errorf("failed to create discovery '%s': %w", name, err))
				continue
			}
			p.discoveries[name] = d
		}
	}

	p.discoveriesInitialized = true
	return p.discoveries, allErrors
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
	p.registrars[name] = registrar
}

// Registrars returns a map of all available service registrars.
// It follows the same caching and creation logic as Discoveries.
func (p *Provider) Registrars() (map[string]registry.KRegistrar, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.registrarsInitialized {
		return p.registrars, nil
	}

	var allErrors error
	if p.config != nil {
		for _, cfg := range p.config.GetConfigs() {
			name := cmp.Or(cfg.Name, cfg.Type)
			if _, exists := p.registrars[name]; exists {
				p.logger.Warnf("registrar '%s' is already registered, skipping config-based creation", name)
				continue
			}
			reg, err := registry.NewRegistrar(cfg, p.opts...)
			if err != nil {
				p.logger.Errorf("failed to create registrar '%s': %v", name, err)
				allErrors = errors.Join(allErrors, fmt.Errorf("failed to create registrar '%s': %w", name, err))
				continue
			}
			p.registrars[name] = reg
		}
	}

	p.registrarsInitialized = true
	return p.registrars, allErrors
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
func (p *Provider) DefaultRegistrar() (registry.KRegistrar, error) {
	p.mu.Lock()
	name := p.defaultRegistrar
	p.mu.Unlock()

	if name == "" {
		return nil, errors.New("default registrar not set")
	}
	return p.Registrar(name)
}
