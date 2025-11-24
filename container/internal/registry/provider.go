package registry

import (
	"cmp"
	"errors"
	"fmt"
	"sync"

	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	"github.com/origadmin/runtime/interfaces/options"
	runtimelog "github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/registry"
)

// Provider implements interfaces.RegistryProvider. It manages the lifecycle of discovery
// and registrar instances, caching them after first creation and allowing for reconfiguration.
// It is safe for concurrent use.
type Provider struct {
	mu                     sync.Mutex
	config                 *discoveryv1.Discoveries
	logger                 *runtimelog.Helper
	opts                   []options.Option
	defaultRegistrar       string // defaultRegistrar from config (active -> default -> single)
	discoveries            map[string]registry.KDiscovery
	registrars             map[string]registry.KRegistrar
	discoveriesInitialized bool
	registrarsInitialized  bool
}

// NewProvider creates a new Provider.
func NewProvider(logger runtimelog.Logger) *Provider {
	return &Provider{
		logger:      runtimelog.NewHelper(logger),
		discoveries: make(map[string]registry.KDiscovery),
		registrars:  make(map[string]registry.KRegistrar),
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

	// Determine the provisional default registrar name based on config priority:
	// 1. 'active' field
	// 2. 'default' field
	// 3. single instance fallback
	var defaultName string
	if cfg != nil {
		defaultName = cmp.Or(cfg.GetActive(), cfg.GetDefault())
		if defaultName == "" && len(cfg.GetConfigs()) == 1 {
			defaultName = cmp.Or(cfg.GetConfigs()[0].GetName(), cfg.GetConfigs()[0].GetType())
		}
	}
	p.defaultRegistrar = defaultName

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
			name := cmp.Or(cfg.GetName(), cfg.GetType())
			if name == "" {
				p.logger.Warnf("discovery configuration is missing a name, using type as fallback: %s", cfg.GetType())
				continue
			}
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

// DefaultRegistrar returns the default service registrar. It performs validation and applies fallbacks.
// The globalDefaultName is provided by the container, having the lowest priority.
func (p *Provider) DefaultRegistrar(globalDefaultName string) (registry.KRegistrar, error) {
	// Ensure all registrars are initialized before we try to find the default.
	registrars, err := p.Registrars()
	if err != nil {
		return nil, err
	}

	p.mu.Lock()
	configDefaultName := p.defaultRegistrar // Default name determined from config (active -> default -> single)
	p.mu.Unlock()

	// Priority 1: Config-based default (active -> default -> single instance)
	if configDefaultName != "" {
		if registrar, ok := registrars[configDefaultName]; ok {
			return registrar, nil
		}
		p.logger.Warnf("config-based default registrar '%s' not found, attempting global default or fallback", configDefaultName)
	}

	// Priority 2: Global default name from options
	if globalDefaultName != "" {
		if registrar, ok := registrars[globalDefaultName]; ok {
			return registrar, nil
		}
		p.logger.Warnf("global default registrar '%s' not found, attempting single instance fallback", globalDefaultName)
	}

	// Priority 3: Fallback to single instance if only one exists
	if len(registrars) == 1 {
		for _, registrar := range registrars {
			return registrar, nil
		}
	}

	return nil, errors.New("no default registrar configured or found, and multiple registrars exist")
}
