package cache

import (
	"errors"
	"fmt"
	"sync"

	"github.com/goexts/generic/cmp"

	runtimelog "github.com/origadmin/runtime/log" // Corrected import alias

	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/data/storage/cache"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
)

// Provider implements storageiface.CacheProvider. It manages the lifecycle of cache
// instances, caching them after first creation and allowing for reconfiguration.
// It is safe for concurrent use.
type Provider struct {
	mu          sync.Mutex
	config      *datav1.Caches
	log         *runtimelog.Helper
	opts        []options.Option
	defaultName string // defaultName from config (active -> default -> single)
	caches      map[string]storageiface.Cache
	initialized bool
}

// NewProvider creates a new Provider.
func NewProvider(logger runtimelog.Logger) *Provider { // Changed logger type
	return &Provider{
		log:    runtimelog.NewHelper(logger),
		caches: make(map[string]storageiface.Cache),
	}
}

// SetConfig updates the provider's configuration. This will clear any previously
// cached instances and cause them to be recreated on the next access, using the new configuration.
// It also provisionally determines the default instance name from the configuration.
func (p *Provider) SetConfig(cfg *datav1.Caches, opts ...options.Option) *Provider {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.config = cfg
	p.opts = opts
	p.initialized = false
	p.caches = make(map[string]storageiface.Cache)

	// Determine the provisional default cache name based on config priority:
	// 1. 'active' field
	// 2. 'default' field
	// 3. single instance fallback
	var defaultName string
	if cfg != nil {
		defaultName = cmp.Or(cfg.GetActive(), cfg.GetDefault())
		if defaultName == "" && len(cfg.GetConfigs()) == 1 {
			defaultName = cmp.Or(cfg.GetConfigs()[0].GetName(), cfg.GetConfigs()[0].GetDriver())
		}
	}
	p.defaultName = defaultName

	return p
}

// RegisterCache allows for manual registration of a cache instance.
func (p *Provider) RegisterCache(name string, cache storageiface.Cache) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.caches[name] = cache
}

// Caches returns a map of all available cache instances.
// On first call, it creates instances from the configuration and caches them.
// Subsequent calls return the cached instances unless SetConfig has been called.
func (p *Provider) Caches() (map[string]storageiface.Cache, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.initialized {
		return p.caches, nil
	}

	var allErrors error
	if p.config != nil {
		for _, cfg := range p.config.GetConfigs() {
			name := cmp.Or(cfg.GetName(), cfg.GetDriver())
			if name == "" {
				p.log.Warnf("cache configuration is missing a name, using driver as fallback: %s", cfg.GetDriver())
				continue
			}
			if _, exists := p.caches[name]; exists {
				p.log.Warnf("cache '%s' is already registered, skipping config-based creation", name)
				continue
			}
			ca, err := cache.New(cfg, p.opts...)
			if err != nil {
				p.log.Errorf("failed to create cache '%s': %v", name, err)
				allErrors = errors.Join(allErrors, fmt.Errorf("failed to create cache '%s': %w", name, err))
				continue
			}
			p.caches[name] = ca
		}
	}

	p.initialized = true
	return p.caches, allErrors
}

// Cache returns a single cache instance by name.
func (p *Provider) Cache(name string) (storageiface.Cache, error) {
	caches, err := p.Caches()
	if err != nil {
		return nil, err
	}
	ca, ok := caches[name]
	if !ok {
		return nil, fmt.Errorf("cache '%s' not found", name)
	}
	return ca, nil
}

// DefaultCache returns the default cache instance. It performs validation and applies fallbacks.
// The globalDefaultName is provided by the container, having the lowest priority.
func (p *Provider) DefaultCache(globalDefaultName string) (storageiface.Cache, error) {
	// Ensure all caches are initialized before we try to find the default.
	caches, err := p.Caches()
	if err != nil {
		return nil, err
	}

	p.mu.Lock()
	configDefaultName := p.defaultName // Default name determined from config (active -> default -> single)
	p.mu.Unlock()

	// Priority 1: Config-based default (active -> default -> single instance)
	if configDefaultName != "" {
		if cache, ok := caches[configDefaultName]; ok {
			return cache, nil
		}
		p.log.Warnf("config-based default cache '%s' not found, attempting global default or fallback", configDefaultName)
	}

	// Priority 2: Global default name from options
	if globalDefaultName != "" {
		if cache, ok := caches[globalDefaultName]; ok {
			return cache, nil
		}
		p.log.Warnf("global default cache '%s' not found, attempting single instance fallback", globalDefaultName)
	}

	// Priority 3: Fallback to single instance if only one exists
	if len(caches) == 1 {
		for _, cache := range caches {
			return cache, nil
		}
	}

	return nil, errors.New("no default cache configured or found, and multiple caches exist")
}
