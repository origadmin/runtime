package cache

import (
	"errors"
	"fmt"
	"sync"

	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/data/storage/cache"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	"github.com/origadmin/runtime/log"
)

// Provider implements storageiface.CacheProvider. It manages the lifecycle of cache
// instances, caching them after first creation and allowing for reconfiguration.
// It is safe for concurrent use.
type Provider struct {
	mu          sync.Mutex
	config      *datav1.Caches
	log         *log.Helper
	opts        []options.Option
	defaultName string
	caches      map[string]storageiface.Cache
	initialized bool
}

// NewProvider creates a new Provider.
func NewProvider(logger log.Logger) *Provider {
	return &Provider{
		log: log.NewHelper(logger),
	}
}

// SetConfig updates the provider's configuration. This will clear any previously
// cached instances and cause them to be recreated on the next access, using the new configuration.
func (p *Provider) SetConfig(cfg *datav1.Caches, opts ...options.Option) *Provider {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.config = cfg
	p.opts = opts
	p.initialized = false
	p.caches = make(map[string]storageiface.Cache)

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
			name := cfg.GetName()
			if name == "" {
				name = cfg.GetDriver()
				p.log.Warnf("cache configuration is missing a name, using driver as fallback: %s", name)
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

// DefaultCache returns the default cache instance.
func (p *Provider) DefaultCache() (storageiface.Cache, error) {
	p.mu.Lock()
	name := p.defaultName
	p.mu.Unlock()

	if name == "" {
		return nil, errors.New("default cache name is not set")
	}
	return p.Cache(name)
}
