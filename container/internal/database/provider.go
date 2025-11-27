package database

import (
	"errors"
	"fmt"
	"maps"
	"sync"

	"github.com/goexts/generic/cmp"

	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/data/storage/database"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	runtimelog "github.com/origadmin/runtime/log"
)

// Provider manages the lifecycle of database instances.
// It uses lazy-loading with sync.Once to ensure instances are created only when needed and in a concurrency-safe manner.
type Provider struct {
	mu           sync.RWMutex
	logger       *runtimelog.Helper
	databases    map[string]storageiface.Database
	config       *datav1.Databases
	opts         []options.Option
	databasesOnce sync.Once
	databasesErr  error
	defaultName  string
}

// NewProvider creates a new, uninitialized Provider instance.
func NewProvider(logger runtimelog.Logger) *Provider {
	return &Provider{
		logger:    runtimelog.NewHelper(logger),
		databases: make(map[string]storageiface.Database),
	}
}

// Initialize configures the provider with the necessary configuration and options.
func (p *Provider) Initialize(cfg *datav1.Databases, opts ...options.Option) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.config = cfg
	p.opts = opts
	if cfg != nil {
		p.defaultName = cmp.Or(cfg.GetActive(), cfg.GetDefault())
	}
}

// RegisterDatabase allows for manual registration of a database instance.
func (p *Provider) RegisterDatabase(name string, db storageiface.Database) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.databases[name]; ok {
		p.logger.Warnf("database '%s' is being overwritten by manual registration", name)
	}
	p.databases[name] = db
}

// Databases returns a map of all available database instances.
// On the first call, it lazily creates and caches instances based on the configuration.
func (p *Provider) Databases() (map[string]storageiface.Database, error) {
	p.databasesOnce.Do(func() {
		p.mu.Lock()
		defer p.mu.Unlock()

		if p.config == nil {
			return
		}
		var allErrors error
		for _, cfg := range p.config.GetConfigs() {
			name := cmp.Or(cfg.GetName(), cfg.GetDialect())
			if name == "" {
				continue
			}
			if _, exists := p.databases[name]; exists {
				continue
			}
			db, err := database.New(cfg, p.opts...)
			if err != nil {
				p.logger.Errorf("failed to create database '%s': %v", name, err)
				allErrors = errors.Join(allErrors, err)
				continue
			}
			p.databases[name] = db
		}
		p.databasesErr = allErrors
	})

	p.mu.RLock()
	defer p.mu.RUnlock()
	return maps.Clone(p.databases), p.databasesErr
}

// Database returns a single database instance by name.
func (p *Provider) Database(name string) (storageiface.Database, error) {
	databases, err := p.Databases()
	if err != nil {
		return nil, err
	}
	db, ok := databases[name]
	if !ok {
		return nil, fmt.Errorf("database '%s' not found", name)
	}
	return db, nil
}

// DefaultDatabase returns the default database instance.
func (p *Provider) DefaultDatabase(globalDefaultName string) (storageiface.Database, error) {
	databases, err := p.Databases()
	if err != nil {
		return nil, err
	}
	if len(databases) == 0 {
		return nil, errors.New("no databases available")
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
		if comp, ok := databases[name]; ok {
			p.logger.Debugf("resolved default database to '%s'", name)
			return comp, nil
		}
	}

	if len(databases) == 1 {
		for name, comp := range databases {
			p.logger.Debugf("no specific default found, falling back to the first available database: '%s'", name)
			return comp, nil
		}
	}

	return nil, errors.New("no default database could be determined")
}
