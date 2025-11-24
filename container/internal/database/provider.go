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

// Provider implements storageiface.DatabaseProvider
type Provider struct {
	config          *datav1.Databases
	log             *log.Helper
	opts            []options.Option // Now stores options passed to SetConfig
	defaultDatabase string
	cachedDatabases map[string]storageiface.Database
	onceDatabases   sync.Once
}

func (p *Provider) DefaultDatabase() (storageiface.Database, error) {
	// Check if defaultDatabase is set
	if p.defaultDatabase == "" {
		return nil, fmt.Errorf("default database name is not set")
	}

	// Return the default database
	return p.Database(p.defaultDatabase)
}

func (p *Provider) RegisterDatabase(name string, db storageiface.Database) {
	// Register the database in the provider
	p.cachedDatabases[name] = db
}

// SetConfig sets the database configurations and dynamic options for the provider.
func (p *Provider) SetConfig(cfg *datav1.Databases, opts ...options.Option) *Provider {
	p.config = cfg
	p.opts = opts // Store the dynamically passed options
	return p
}

// Databases returns all the configured databases.
func (p *Provider) Databases() (map[string]storageiface.Database, error) {
	var allErrors error
	p.onceDatabases.Do(func() {
		if p.config == nil || len(p.config.GetConfigs()) == 0 {
			p.log.Infow("msg", "no database configurations found")
			return
		}

		for _, cfg := range p.config.GetConfigs() {
			name := cfg.GetName()
			if name == "" {
				p.log.Warnf("database configuration is missing a name, using driver as fallback: %s", cfg.GetDialect())
				name = cfg.GetDialect()
			}
			// Pass the stored options to the database creation
			db, err := database.New(cfg, p.opts...)
			if err != nil {
				p.log.Errorf("failed to create database '%s': %v", name, err)
				allErrors = errors.Join(allErrors, fmt.Errorf("failed to create database '%s': %w", name, err))
				continue
			}
			p.cachedDatabases[name] = db
		}
	})
	return p.cachedDatabases, allErrors
}

// Database returns a specific database by name.
func (p *Provider) Database(name string) (storageiface.Database, error) {
	s, err := p.Databases()
	if err != nil {
		return nil, err
	}
	db, ok := s[name]
	if !ok {
		return nil, fmt.Errorf("database '%s' not found", name)
	}
	return db, nil
}

// NewProvider creates a new Provider.
// It no longer receives opts, as options are passed dynamically via SetConfig.
func NewProvider(logger log.Logger) *Provider {
	helper := log.NewHelper(logger)
	return &Provider{
		log:             helper,
		cachedDatabases: make(map[string]storageiface.Database),
	}
}
