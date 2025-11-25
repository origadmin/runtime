package objectstore

import (
	"errors"
	"fmt"
	"sync"

	"github.com/goexts/generic/cmp"

	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/container/internal/util" // Import util package
	"github.com/origadmin/runtime/data/storage/objectstore"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	runtimelog "github.com/origadmin/runtime/log"
)

// Provider implements interfaces.ObjectStoreProvider. It manages the lifecycle of object store
// instances, caching them after first creation.
// It is safe for concurrent use.
type Provider struct {
	mu              sync.Mutex
	config          *datav1.ObjectStores
	log             *runtimelog.Helper
	opts            []options.Option // This field stores options set via SetOptions
	objectStoreName string // objectStoreName from config (active -> default -> single)
	objectStores    map[string]storageiface.ObjectStore
	initialized     bool
}

// NewProvider creates a new Provider instance. Options are set via SetOptions.
func NewProvider(logger runtimelog.Logger) *Provider {
	p := &Provider{
		log:          runtimelog.NewHelper(logger),
		objectStores: make(map[string]storageiface.ObjectStore),
	}
	return p
}

// SetOptions sets the functional options for the provider.
// This will clear any previously cached instances and cause them to be recreated on the next access,
// using the new options and the structural configuration.
func (p *Provider) SetOptions(opts ...options.Option) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.opts = opts
	p.initialized = false
	p.objectStores = make(map[string]storageiface.ObjectStore)
}

// SetConfig updates the provider's structural configuration.
// This will clear any previously cached instances and cause them to be recreated on the next access,
// using the new structural configuration and the functional options provided at NewProvider time.
// It also provisionally determines the default instance name from the configuration.
func (p *Provider) SetConfig(cfg *datav1.ObjectStores) *Provider {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.config = cfg
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
// Subsequent calls return the cached instances.
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
	stores, err := p.ObjectStores()
	if err != nil {
		return nil, err
	}
	if len(stores) == 0 {
		return nil, errors.New("no object stores available")
	}

	p.mu.Lock()
	configDefaultName := p.objectStoreName
	p.mu.Unlock()

	var prioritizedNames []string

	// Priority 1: Config-based default
	if configDefaultName != "" {
		prioritizedNames = append(prioritizedNames, configDefaultName)
	}

	// Priority 2: External globalDefaultName
	if globalDefaultName != "" {
		prioritizedNames = append(prioritizedNames, globalDefaultName)
	}

	// Priority 3: GlobalDefaultKey (as a final fallback)
	prioritizedNames = append(prioritizedNames, interfaces.GlobalDefaultKey)

	// Call the utility function to determine the default component
	name, value, err := util.DefaultComponent(stores, prioritizedNames...)
	if err == nil {
		p.log.Debugf("resolved default object store to '%s'", name)
		return value, nil
	}

	// If util.DefaultComponent returned an error, handle it here.
	// The error from util.DefaultComponent already describes why a default wasn't found.
	return nil, fmt.Errorf("no default object store found: %w", err)
}
