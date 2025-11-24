package database

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/log"

	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/data/storage/database"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
)

// Provider implements storageiface.DatabaseProvider. It manages the lifecycle of database
// instances, caching them after first creation and allowing for reconfiguration.
// It is safe for concurrent use.
type Provider struct {
	mu              sync.Mutex
	config          *datav1.Databases
	log             *log.Helper
	opts            []options.Option
	defaultDatabase string
	databases       map[string]storageiface.Database
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
func (p *Provider) SetConfig(cfg *datav1.Databases, opts ...options.Option) *Provider {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.config = cfg
	p.opts = opts
	p.initialized = false
	p.databases = make(map[string]storageiface.Database)

	return p
}

// RegisterDatabase allows for manual registration of a database instance.
func (p *Provider) RegisterDatabase(name string, db storageiface.Database) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.databases[name] = db
}

// Databases returns a map of all available database instances.
// On first call, it creates instances from the configuration and caches them.
// Subsequent calls return the cached instances unless SetConfig has been called.
func (p *Provider) Databases() (map[string]storageiface.Database, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.initialized {
		return p.databases, nil
	}

	var allErrors error
	if p.config != nil {
		for _, cfg := range p.config.GetConfigs() {
			name := cfg.GetName()
			if name == "" {
				name = cfg.GetDialect()
				p.log.Warnf("database configuration is missing a name, using dialect as fallback: %s", name)
			}
			if _, exists := p.databases[name]; exists {
				p.log.Warnf("database '%s' is already registered, skipping config-based creation", name)
				continue
			}
			db, err := database.New(cfg, p.opts...)
			if err != nil {
				p.log.Errorf("failed to create database '%s': %v", name, err)
				allErrors = errors.Join(allErrors, fmt.Errorf("failed to create database '%s': %w", name, err))
				continue
			}
			p.databases[name] = db
		}
	}

	p.initialized = true
	return p.databases, allErrors
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
func (p *Provider) DefaultDatabase() (storageiface.Database, error) {
	p.mu.Lock()
	name := p.defaultDatabase
	p.mu.Unlock()

	if name == "" {
		return nil, errors.New("default database name is not set")
	}
	return p.Database(name)
}
