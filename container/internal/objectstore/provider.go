package objectstore

import (
	"errors"
	"fmt"
	"sync"

	"github.com/goexts/generic/cmp"

	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/data/storage/objectstore"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	runtimelog "github.com/origadmin/runtime/log"
)

// Provider implements interfaces.ObjectStoreProvider. It manages the lifecycle of object store
// instances, caching them after first creation and allowing for reconfiguration.
// It is safe for concurrent use.
type Provider struct {
	mu              sync.Mutex
	config          *datav1.ObjectStores
	log             *runtimelog.Helper
	opts            []options.Option
	objectStoreName string // objectStoreName from config (active -> default -> single)
	objectStores    map[string]storageiface.ObjectStore
	initialized     bool
}

// NewProvider creates a new Provider.
func NewProvider(logger runtimelog.Logger) *Provider { // Changed logger type
	return &Provider{
		log:          runtimelog.NewHelper(logger),
		objectStores: make(map[string]storageiface.ObjectStore),
	}
}

// SetConfig updates the provider's configuration. This will clear any previously
// cached instances and cause them to be recreated on the next access, using the new configuration.
// It also provisionally determines the default instance name from the configuration.
func (p *Provider) SetConfig(cfg *datav1.ObjectStores, opts ...options.Option) *Provider {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.config = cfg
	p.opts = opts
	p.initialized = false
	p.objectStores = make(map[string]storageiface.ObjectStore)

	// Determine the provisional default object store name based on config priority:
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
	p.objectStoreName = defaultName

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
			name := cmp.Or(cfg.GetName(), cfg.GetDriver())
			if name == "" {
				p.log.Warnf("object store configuration is missing a name, using driver as fallback: %s", cfg.GetDriver())
				continue
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

// DefaultObjectStore returns the default object store instance. It performs validation and applies fallbacks.
// The globalDefaultName is provided by the container, having the lowest priority.
func (p *Provider) DefaultObjectStore(globalDefaultName string) (storageiface.ObjectStore, error) {
	// Ensure all stores are initialized before we try to find the default.
	stores, err := p.ObjectStores()
	if err != nil {
		return nil, err
	}

	p.mu.Lock()
	configDefaultName := p.objectStoreName // Default name determined from config (active -> default -> single)
	p.mu.Unlock()

	// Priority 1: Config-based default (active -> default -> single instance)
	if configDefaultName != "" {
		if store, ok := stores[configDefaultName]; ok {
			return store, nil
		}
		p.log.Warnf("config-based default object store '%s' not found, attempting global default or fallback", configDefaultName)
	}

	// Priority 2: Global default name from options
	if globalDefaultName != "" {
		if store, ok := stores[globalDefaultName]; ok {
			return store, nil
		}
		p.log.Warnf("global default object store '%s' not found, attempting single instance fallback", globalDefaultName)
	}

	// Priority 3: Fallback to single instance if only one exists
	if len(stores) == 1 {
		for _, store := range stores {
			return store, nil
		}
	}

	return nil, errors.New("no default object store configured or found, and multiple object stores exist")
}
