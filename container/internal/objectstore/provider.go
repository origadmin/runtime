package objectstore

import (
	"errors"
	"fmt"
	"maps"
	"sync"

	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/contracts"
	"github.com/origadmin/runtime/contracts/options"
	storageiface "github.com/origadmin/runtime/contracts/storage"
	"github.com/origadmin/runtime/helpers/configutil"
	runtimelog "github.com/origadmin/runtime/log"
)

// Provider manages the lifecycle of object store instances.
type Provider struct {
	mu           sync.RWMutex
	logger       *runtimelog.Helper
	objectStores map[string]storageiface.ObjectStore
	config       *datav1.ObjectStores
	opts         []options.Option
	defaultName  string
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
		// Normalize the configuration to determine the correct default and config list.
		normalizedDefault, normalizedConfigs, err := configutil.Normalize(cfg.GetActive(), cfg.GetDefault(), cfg.GetConfigs())
		if err != nil {
			// Log the error and proceed with a best-effort approach.
			p.logger.Warnf("failed to normalize object store configuration: %v. Using active name as default.", err)
			p.defaultName = cfg.GetActive() // Fallback to original behavior
			return
		}

		// Update the provider's state with the normalized configuration.
		if normalizedDefault != nil {
			p.defaultName = normalizedDefault.GetName()
		}
		// The provider uses p.config to create object stores, so we update it.
		// NOTE: This modifies the original config object passed to Initialize.
		// This is acceptable as the provider takes ownership of the config.
		p.config.Default = normalizedDefault
		p.config.Configs = normalizedConfigs
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
// It only returns instances that have been manually registered.
func (p *Provider) ObjectStores() (map[string]storageiface.ObjectStore, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return maps.Clone(p.objectStores), nil
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
	prioritizedNames = append(prioritizedNames, contracts.GlobalDefaultKey)

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
