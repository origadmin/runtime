package database

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/log"

	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/data/storage/database"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
)

// Provider implements interfaces.DatabaseProvider
type Provider struct {
	config          *datav1.Databases
	log             *log.Helper
	opts            []options.Option
	cachedDatabases map[string]interfaces.Database
	onceDatabases   sync.Once
}

func (p *Provider) RegisterDatabase(name string, db interfaces.Database) {
	//TODO implement me
	panic("implement me")
}

// NewProvider creates a new Provider.
func NewProvider(logger log.Logger, opts []options.Option) *Provider {
	helper := log.NewHelper(logger)
	return &Provider{
		log:             helper,
		opts:            opts,
		cachedDatabases: make(map[string]interfaces.Database),
	}
}

// SetConfig sets the database configurations for the provider.
func (p *Provider) SetConfig(cfg *datav1.Databases) *Provider {
	p.config = cfg
	return p
}

// Databases returns all the configured databases.
func (p *Provider) Databases() (map[string]interfaces.Database, error) {
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
func (p *Provider) Database(name string) (interfaces.Database, error) {
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
