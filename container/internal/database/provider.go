package database

import (
	"cmp"
	"errors"
	"fmt"
	"sync"

	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/container/internal/util" // Import util package
	"github.com/origadmin/runtime/data/storage/database"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	runtimelog "github.com/origadmin/runtime/log"
)

// Provider implements storageiface.DatabaseProvider. It manages the lifecycle of database
// instances, caching them after first creation.
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

// NewProvider creates a new Provider instance, applying functional options immediately.
func NewProvider(logger runtimelog.Logger, opts ...options.Option) *Provider {
	p := &Provider{
		log:       runtimelog.NewHelper(logger),
		databases: make(map[string]storageiface.Database),
		opts:      opts, // Store functional options here
	}
	return p
}

// SetConfig updates the provider's structural configuration.
// This will clear any previously cached instances and cause them to be recreated on the next access,
// using the new structural configuration and the functional options provided at NewProvider time.
// It also provisionally determines the default instance name from the configuration.
func (p *Provider) SetConfig(cfg *datav1.Databases) *Provider {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.config = cfg
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
// Subsequent calls return the cached instances.
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
	databases, err := p.Databases()
	if err != nil {
		return nil, err
	}
	if len(databases) == 0 {
		return nil, errors.New("no databases available")
	}

	p.mu.Lock()
	configDefaultName := p.defaultDatabase
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
	name, value, err := util.DefaultComponent(databases, prioritizedNames...)
	if err == nil {
		p.log.Debugf("resolved default database to '%s'", name)
		return value, nil
	}

	// If util.DefaultComponent returned an error, handle it here.
	// The error from util.DefaultComponent already describes why a default wasn't found.
	return nil, fmt.Errorf("no default database found: %w", err)
}
