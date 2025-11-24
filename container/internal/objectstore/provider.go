package objectstore

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/log"

	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/data/storage/objectstore"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
)

// Provider implements interfaces.ObjectStoreProvider. It manages the lifecycle of object store
// instances, caching them after first creation and allowing for reconfiguration.
// It is safe for concurrent use.
type Provider struct {
	mu              sync.Mutex
	config          *datav1.ObjectStores
	log             *log.Helper
	opts            []options.Option
	objectStoreName string
	objectStores    map[string]storageiface.ObjectStore
	initialized     bool
}

// NewProvider creates a new Provider.
func NewProvider(logger log.Logger) *Provider {
	return &Provider{
		log: log.NewHelper(logger),
	}
}

// SetConfig updates the provider's configuration. This will clear any previously
// cached instances and cause them to be recreated on the next access, using the new configuration.
func (p *Provider) SetConfig(cfg *datav1.ObjectStores, opts ...options.Option) *Provider {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.config = cfg
	p.opts = opts
	p.initialized = false
	p.objectStores = make(map[string]storageiface.ObjectStore)

	return p
}

// RegisterObjectStore allows for manual registration of an object store instance.
func (p *Provider) RegisterObjectStore(name string, store storageiface.ObjectStore) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.objectStores[name] = store
}

// ObjectStores returns a map of all available object store instances.
// On first call, it creates instances from the configuration and caches them.
// Subsequent calls return the cached instances unless SetConfig has been called.
func (p *Provider) ObjectStores() (map[string]storageiface.ObjectStore, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.initialized {
		return p.objectStores, nil
	}

	var allErrors error
	if p.config != nil {
		for _, cfg := range p.config.GetConfigs() {
			name := cfg.GetName()
			if name == "" {
				name = cfg.GetDriver()
				p.log.Warnf("object store configuration is missing a name, using driver as fallback: %s", name)
			}
			if _, exists := p.objectStores[name]; exists {
				p.log.Warnf("object store '%s' is already registered, skipping config-based creation", name)
				continue
			}
			os, err := objectstore.New(cfg, p.opts...)
			if err != nil {
				p.log.Errorf("failed to create object store '%s': %v", name, err)
				allErrors = errors.Join(allErrors, fmt.Errorf("failed to create object store '%s': %w", name, err))
				continue
			}
			p.objectStores[name] = os
		}
	}

	p.initialized = true
	return p.objectStores, allErrors
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
func (p *Provider) DefaultObjectStore() (storageiface.ObjectStore, error) {
	p.mu.Lock()
	name := p.objectStoreName
	p.mu.Unlock()

	if name == "" {
		return nil, errors.New("default object store name is not set")
	}
	return p.ObjectStore(name)
}
