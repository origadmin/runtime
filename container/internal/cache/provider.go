package cache

import (
	"errors"
	"fmt"
	"maps"
	"sync"

	"github.com/goexts/generic/cmp"

	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/data/storage/cache"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	runtimelog "github.com/origadmin/runtime/log"
)

// Provider manages the lifecycle of cache instances.
// It uses lazy-loading with sync.Once to ensure instances are created only when needed and in a concurrency-safe manner.
type Provider struct {
	mu          sync.RWMutex
	logger      *runtimelog.Helper
	caches      map[string]storageiface.Cache
	config      *datav1.Caches
	opts        []options.Option
	cachesOnce  sync.Once
	cachesErr   error
	defaultName string
}

// NewProvider creates a new, uninitialized Provider instance.
func NewProvider(logger runtimelog.Logger) *Provider {
	return &Provider{
		logger: runtimelog.NewHelper(logger),
		caches: make(map[string]storageiface.Cache),
	}
}

// Initialize configures the provider with the necessary configuration and options.
func (p *Provider) Initialize(cfg *datav1.Caches, opts ...options.Option) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.config = cfg
	p.opts = opts
	if cfg != nil {
		p.defaultName = cmp.Or(cfg.GetActive(), cfg.GetDefault())
	}
}

// RegisterCache allows for manual registration of a cache instance.
func (p *Provider) RegisterCache(name string, c storageiface.Cache) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.caches[name]; ok {
		p.logger.Warnf("cache '%s' is being overwritten by manual registration", name)
	}
	p.caches[name] = c
}

// Caches returns a map of all available cache instances.
// On the first call, it lazily creates and caches instances based on the configuration.
func (p *Provider) Caches() (map[string]storageiface.Cache, error) {
	p.cachesOnce.Do(func() {
		p.mu.Lock()
		defer p.mu.Unlock()

		if p.config == nil {
			return
		}
		var allErrors error
		for _, cfg := range p.config.GetConfigs() {
			name := cmp.Or(cfg.GetName(), cfg.GetDriver())
			if name == "" {
				continue
			}
			if _, exists := p.caches[name]; exists {
				continue
			}
			ca, err := cache.New(cfg, p.opts...)
			if err != nil {
				p.logger.Errorf("failed to create cache '%s': %v", name, err)
				allErrors = errors.Join(allErrors, err)
				continue
			}
			p.caches[name] = ca
		}
		p.cachesErr = allErrors
	})

	p.mu.RLock()
	defer p.mu.RUnlock()
	return maps.Clone(p.caches), p.cachesErr
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
func (p *Provider) DefaultCache(globalDefaultName string) (storageiface.Cache, error) {
	caches, err := p.Caches()
	if err != nil {
		return nil, err
	}
	if len(caches) == 0 {
		return nil, errors.New("no caches available")
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	var prioritizedNames []string
	if p.defaultName != "" {
		prioritizedNames = append(prioritizedNames, p.defaultName)
	}
	if globalDefaultName != "" {
		prioritizedNames = append(prioritizedNames, globalDefaultName)
	}
	prioritizedNames = append(prioritizedNames, interfaces.GlobalDefaultKey)

	for _, name := range prioritizedNames {
		if comp, ok := caches[name]; ok {
			p.logger.Debugf("resolved default cache to '%s'", name)
			return comp, nil
		}
	}

	if len(caches) == 1 {
		for name, comp := range caches {
			p.logger.Debugf("no specific default found, falling back to the first available cache: '%s'", name)
			return comp, nil
		}
	}

	return nil, errors.New("no default cache could be determined")
}
