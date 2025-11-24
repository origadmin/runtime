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

// Provider implements interfaces.ObjectStoreProvider
type Provider struct {
	config             *datav1.ObjectStores
	log                *log.Helper
	opts               []options.Option
	objectStoreName    string
	cachedObjectStores map[string]storageiface.ObjectStore
	onceObjectStores   sync.Once
}

func (p *Provider) DefaultObjectStore() (storageiface.ObjectStore, error) {
	// Check if objectStoreName is set
	if p.objectStoreName == "" {
		return nil, fmt.Errorf("object store name is not set")
	}

	return p.ObjectStore(p.objectStoreName)
}

func (p *Provider) RegisterObjectStore(name string, store storageiface.ObjectStore) {
	// Register the object store in the provider
	p.cachedObjectStores[name] = store
}

// NewProvider creates a new Provider.
// It no longer receives opts, as options are passed dynamically via SetConfig.
func NewProvider(logger log.Logger) *Provider {
	helper := log.NewHelper(logger)
	return &Provider{
		log:                helper,
		cachedObjectStores: make(map[string]storageiface.ObjectStore),
	}
}

// SetConfig sets the object store configurations and dynamic options for the provider.
func (p *Provider) SetConfig(cfg *datav1.ObjectStores, opts ...options.Option) *Provider {
	p.config = cfg
	p.opts = opts // Store the dynamically passed options
	return p
}

// ObjectStores returns all the configured object stores.
func (p *Provider) ObjectStores() (map[string]storageiface.ObjectStore, error) {
	var allErrors error
	p.onceObjectStores.Do(func() {
		if p.config == nil || len(p.config.GetConfigs()) == 0 {
			p.log.Infow("msg", "no object store configurations found")
			return
		}

		for _, cfg := range p.config.GetConfigs() {
			name := cfg.GetName()
			if name == "" {
				p.log.Warnf("object store configuration is missing a name, using driver as fallback: %s", cfg.GetDriver())
				name = cfg.GetDriver()
			}
			// Pass the stored options to the object store creation
			os, err := objectstore.New(cfg, p.opts...)
			if err != nil {
				p.log.Errorf("failed to create object store '%s': %v", name, err)
				allErrors = errors.Join(allErrors, fmt.Errorf("failed to create object store '%s': %w", name, err))
				continue
			}
			p.cachedObjectStores[name] = os
		}
	})
	return p.cachedObjectStores, allErrors
}

// ObjectStore returns a specific object store by name.
func (p *Provider) ObjectStore(name string) (storageiface.ObjectStore, error) {
	s, err := p.ObjectStores()
	if err != nil {
		return nil, err
	}
	os, ok := s[name]
	if !ok {
		return nil, fmt.Errorf("object store '%s' not found", name)
	}
	return os, nil
}
