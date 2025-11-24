package database

import (
	"cmp"
	"errors"
	"fmt"
	"sync"

	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/data/storage/database"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	runtimelog "github.com/origadmin/runtime/log"
)

// Provider implements storageiface.DatabaseProvider. It manages the lifecycle of database
// instances, caching them after first creation and allowing for reconfiguration.
// It is safe for concurrent use.
type Provider struct {
	mu              sync.Mutex
	config          *datav1.Databases
	log             *runtimelog.Helper
	opts            []options.Option
	defaultDatabase string // defaultDatabase from config (active -> default -> single)
	databases       map[string]storageiface.Database
	initialized     bool
}

// NewProvider creates a new Provider.
func NewProvider(logger runtimelog.Logger) *Provider {
	return &Provider{
		log:       runtimelog.NewHelper(logger),
		databases: make(map[string]storageiface.Database),
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

	// Determine the provisional default database name based on config priority:
	// 1. 'active' field
	// 2. 'default' field
	// 3. single instance fallback
	var defaultName string
	if cfg != nil {
		defaultName = cmp.Or(cfg.GetActive(), cfg.GetDefault())
		if defaultName == "" && len(cfg.GetConfigs()) == 1 {
			defaultName = cmp.Or(cfg.GetConfigs()[0].GetName(), cfg.GetConfigs()[0].GetDialect())
		}
	}
	p.defaultDatabase = defaultName

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
			name := cmp.Or(cfg.GetName(), cfg.GetDialect())
			if name == "" {
				p.log.Warnf("database configuration is missing a name, using dialect as fallback: %s", cfg.GetDialect())
				continue
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

// DefaultDatabase returns the default database instance. It performs validation and applies fallbacks.
// The globalDefaultName is provided by the container, having the lowest priority.
func (p *Provider) DefaultDatabase(globalDefaultName string) (storageiface.Database, error) {
	// Ensure all databases are initialized before we try to find the default.
	databases, err := p.Databases()
	if err != nil {
		return nil, err
	}

	p.mu.Lock()
	configDefaultName := p.defaultDatabase // Default name determined from config (active -> default -> single)
	p.mu.Unlock()

	// Priority 1: Config-based default (active -> default -> single instance)
	if configDefaultName != "" {
		if db, ok := databases[configDefaultName]; ok {
			return db, nil
		}
		p.log.Warnf("config-based default database '%s' not found, attempting global default or fallback", configDefaultName)
	}

	// Priority 2: Global default name from options
	if globalDefaultName != "" {
		if db, ok := databases[globalDefaultName]; ok {
			return db, nil
		}
		p.log.Warnf("global default database '%s' not found, attempting single instance fallback", globalDefaultName)
	}

	// Priority 3: Fallback to single instance if only one exists
	if len(databases) == 1 {
		for _, db := range databases {
			return db, nil
		}
	}

	return nil, errors.New("no default database configured or found, and multiple databases exist")
}
