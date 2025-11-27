package objectstore

import (
	"errors"
	"fmt"
	"maps"
	"sync"

	"github.com/goexts/generic/cmp"

	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/data/storage/objectstore"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	runtimelog "github.com/origadmin/runtime/log"
)

// Provider manages the lifecycle of object store instances.
// It uses lazy-loading with sync.Once to ensure instances are created only when needed and in a concurrency-safe manner.
type Provider struct {
	mu             sync.RWMutex
	logger         *runtimelog.Helper
	objectStores   map[string]storageiface.ObjectStore
	config         *datav1.ObjectStores
	opts           []options.Option
	objectStoresOnce sync.Once
	objectStoresErr  error
	defaultName    string
}

// NewProvider creates a new, uninitialized Provider instance.
func NewProvider(logger runtimelog.Logger) *Provider {
	return &Provider{
		logger:       runtimelog.NewHelper(logger),
		objectStores: make(map[string]storageiface.ObjectStore),
	}
}

// Initialize configures the provider with the necessary configuration and options.
func (p *Provider) Initialize(cfg *datav1.ObjectStores, opts ...options.Option) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.config = cfg
	p.opts = opts
	if cfg != nil {
		p.defaultName = cmp.Or(cfg.GetActive(), cfg.GetDefault())
	}
}

// RegisterObjectStore allows for manual registration of an object store instance.
func (p *Provider) RegisterObjectStore(name string, store storageiface.ObjectStore) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.objectStores[name]; ok {
		p.logger.Warnf("object store '%s' is being overwritten by manual registration", name)
	}
	p.objectStores[name] = store
}

// ObjectStores returns a map of all available object store instances.
// On the first call, it lazily creates and caches instances based on the configuration.
func (p *Provider) ObjectStores() (map[string]storageiface.ObjectStore, error) {
	p.objectStoresOnce.Do(func() {
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
			if _, exists := p.objectStores[name]; exists {
				continue
			}
			os, err := objectstore.New(cfg, p.opts...)
			if err != nil {
				p.logger.Errorf("failed to create object store '%s': %v", name, err)
				allErrors = errors.Join(allErrors, err)
				continue
			}
			p.objectStores[name] = os
		}
		p.objectStoresErr = allErrors
	})

	p.mu.RLock()
	defer p.mu.RUnlock()
	return maps.Clone(p.objectStores), p.objectStoresErr
}

// ObjectStore returns a single object store instance by name.
func (p *Provider) ObjectStore(name string) (storageiface.ObjectStore, error) {
	stores, err := p.ObjectStores()
	if err != nil {
		return nil, err
	}
	os, ok := stores[name]
	if !ok {
		return nil, fmt.Errorf("object store '%s' not found", name)
	}
	return os, nil
}

// DefaultObjectStore returns the default object store instance.
func (p *Provider) DefaultObjectStore(globalDefaultName string) (storageiface.ObjectStore, error) {
	stores, err := p.ObjectStores()
	if err != nil {
		return nil, err
	}
	if len(stores) == 0 {
		return nil, errors.New("no object stores available")
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
		if comp, ok := stores[name]; ok {
			p.logger.Debugf("resolved default object store to '%s'", name)
			return comp, nil
		}
	}

	if len(stores) == 1 {
		for name, comp := range stores {
			p.logger.Debugf("no specific default found, falling back to the first available object store: '%s'", name)
			return comp, nil
		}
	}

	return nil, errors.New("no default object store could be determined")
}
