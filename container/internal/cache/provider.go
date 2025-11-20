package cache

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/log"

	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/data/storage/cache"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
)

// Provider implements interfaces.CacheProvider
type Provider struct {
	config       *datav1.Caches
	log          *log.Helper
	opts         []options.Option
	cachedCaches map[string]interfaces.Cache
	onceCaches   sync.Once
}

func (p *Provider) RegisterCache(name string, cache interfaces.Cache) {
	//TODO implement me
	panic("implement me")
}

// NewProvider creates a new Provider.
func NewProvider(logger log.Logger, opts []options.Option) *Provider {
	helper := log.NewHelper(logger)
	return &Provider{
		log:          helper,
		opts:         opts,
		cachedCaches: make(map[string]interfaces.Cache),
	}
}

// SetConfig sets the cache configurations for the provider.
func (p *Provider) SetConfig(cfg *datav1.Caches) *Provider {
	p.config = cfg
	return p
}

// Caches returns all the configured caches.
func (p *Provider) Caches() (map[string]interfaces.Cache, error) {
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
func (p *Provider) Cache(name string) (interfaces.Cache, error) {
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
