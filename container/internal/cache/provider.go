package cache

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/log"

	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/data/storage/cache"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
)

// Provider implements storageiface.CacheProvider
type Provider struct {
	config       *datav1.Caches
	log          *log.Helper
	opts         []options.Option // Now stores options passed to SetConfig
	defaultName  string
	cachedCaches map[string]storageiface.Cache
	onceCaches   sync.Once
}

func (p *Provider) DefaultCache() (storageiface.Cache, error) {
	// Check if defaultName is set
	if p.defaultName == "" {
		return nil, fmt.Errorf("default cache name is not set")
	}

	defaultCache, err := p.Cache(p.defaultName)
	if err != nil {
		return nil, err
	}
	return defaultCache, nil
}

func (p *Provider) RegisterCache(name string, cache storageiface.Cache) {
	// Register the cache in the provider
	p.cachedCaches[name] = cache
}

// SetConfig sets the cache configurations and dynamic options for the provider.
func (p *Provider) SetConfig(cfg *datav1.Caches, opts ...options.Option) *Provider {
	p.config = cfg
	p.opts = opts // Store the dynamically passed options
	return p
}

// Caches returns all the configured caches.
func (p *Provider) Caches() (map[string]storageiface.Cache, error) {
	var allErrors error
	p.onceCaches.Do(func() {
		if p.config == nil || len(p.config.GetConfigs()) == 0 {
			p.log.Infow("msg", "no cache configurations found")
			return
		}

		for _, cfg := range p.config.GetConfigs() {
			name := cfg.GetName()
			if name == "" {
				p.log.Warnf("cache configuration is missing a name, using driver as fallback: %s", cfg.GetDriver())
				name = cfg.GetDriver()
			}
			// Pass the stored options to the cache creation
			ca, err := cache.New(cfg, p.opts...)
			if err != nil {
				p.log.Errorf("failed to create cache '%s': %v", name, err)
				allErrors = errors.Join(allErrors, fmt.Errorf("failed to create cache '%s': %w", name, err))
				continue
			}
			p.cachedCaches[name] = ca
		}
	})
	return p.cachedCaches, allErrors
}

// Cache returns a specific cache by name.
func (p *Provider) Cache(name string) (storageiface.Cache, error) {
	s, err := p.Caches()
	if err != nil {
		return nil, err
	}
	ca, ok := s[name]
	if !ok {
		return nil, fmt.Errorf("cache '%s' not found", name)
	}
	return ca, nil
}

// NewProvider creates a new Provider.
// It no longer receives opts, as options are passed dynamically via SetConfig.
func NewProvider(logger log.Logger) *Provider {
	helper := log.NewHelper(logger)
	return &Provider{
		log:          helper,
		cachedCaches: make(map[string]storageiface.Cache),
	}
}
