package middleware

import (
	"fmt"
	"maps"
	"sync"

	kratosMiddleware "github.com/go-kratos/kratos/v2/middleware"
	"github.com/goexts/generic/cmp"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/interfaces/options"
	runtimelog "github.com/origadmin/runtime/log"
	runtimeMiddleware "github.com/origadmin/runtime/middleware"
)

// Provider manages the lifecycle of client and server middleware instances.
// It uses lazy-loading with sync.Once to ensure instances are created only when needed and in a concurrency-safe manner.
type Provider struct {
	mu                sync.RWMutex
	logger            *runtimelog.Helper
	orderedNames      []string
	clientMiddlewares map[string]kratosMiddleware.Middleware
	serverMiddlewares map[string]kratosMiddleware.Middleware
	config            *middlewarev1.Middlewares
	opts              []options.Option
	clientMWsOnce     sync.Once
	serverMWsOnce     sync.Once
	clientMWsErr      error
	serverMWsErr      error
}

// NewProvider creates a new, uninitialized Provider instance.
func NewProvider(logger runtimelog.Logger) *Provider {
	return &Provider{
		logger:            runtimelog.NewHelper(logger),
		clientMiddlewares: make(map[string]kratosMiddleware.Middleware),
		serverMiddlewares: make(map[string]kratosMiddleware.Middleware),
	}
}

// Initialize configures the provider with the necessary configuration and options.
func (p *Provider) Initialize(cfg *middlewarev1.Middlewares, opts ...options.Option) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.config = cfg
	p.orderedNames = make([]string, 0, len(cfg.GetConfigs()))
	for _, cfg := range cfg.GetConfigs() {
		name := cmp.Or(cfg.Name, cfg.Type)
		if name == "" {
			continue
		}
		p.orderedNames = append(p.orderedNames, name)
	}
	p.opts = opts
}

// RegisterClientMiddleware allows for manual registration of a client middleware instance.
func (p *Provider) RegisterClientMiddleware(name string, mw kratosMiddleware.Middleware) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.clientMiddlewares[name]; ok {
		p.logger.Warnf("client middleware '%s' is being overwritten by manual registration", name)
	}
	p.clientMiddlewares[name] = mw
}

// ClientMiddlewares returns a map of all available client middleware instances.
// On the first call, it lazily creates and caches instances based on the configuration.
func (p *Provider) ClientMiddlewares() (map[string]kratosMiddleware.Middleware, error) {
	p.clientMWsOnce.Do(func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.logger.Debugf("middleware config: %+v", p.config)
		if p.config == nil {
			return
		}
		var allErrors error

		for _, cfg := range p.config.GetConfigs() {
			name := cmp.Or(cfg.Name, cfg.Type)
			if name == "" {
				continue
			}
			if _, exists := p.clientMiddlewares[name]; exists {
				continue
			}
			// Use the new WithClientCarrier option to pass only client middlewares
			opts := append(p.opts, runtimeMiddleware.WithClientCarrier(p.clientMiddlewares))
			// Attempt to create a client middleware. If the factory returns ok=false,
			// it means this config is not for a client middleware, so we just skip it.
			if cm, ok := runtimeMiddleware.NewClient(cfg, opts...); ok {
				p.clientMiddlewares[name] = cm
			} else {
				// If NewClient returns false, it means this config was not for a client middleware.
				// We don't treat this as an error, just skip it.
				continue
			}
		}
		p.clientMWsErr = allErrors
	})

	p.mu.RLock()
	defer p.mu.RUnlock()
	return maps.Clone(p.clientMiddlewares), p.clientMWsErr
}

// ClientMiddleware returns a single client middleware instance by name.
func (p *Provider) ClientMiddleware(name string) (kratosMiddleware.Middleware, error) {
	middlewares, err := p.ClientMiddlewares()
	if err != nil {
		return nil, err
	}
	mw, ok := middlewares[name]
	if !ok {
		return nil, fmt.Errorf("client middleware '%s' not found", name)
	}
	return mw, nil
}

// RegisterServerMiddleware allows for manual registration of a server middleware instance.
func (p *Provider) RegisterServerMiddleware(name string, mw kratosMiddleware.Middleware) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.serverMiddlewares[name]; ok {
		p.logger.Warnf("server middleware '%s' is being overwritten by manual registration", name)
	}
	p.serverMiddlewares[name] = mw
	p.orderedNames = append(p.orderedNames, name)
}

// ServerMiddlewares returns a map of all available server middleware instances.
// It follows the same lazy-loading and caching logic as ClientMiddlewares.
func (p *Provider) ServerMiddlewares() (map[string]kratosMiddleware.Middleware, error) {
	p.serverMWsOnce.Do(func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.logger.Debugf("middleware config: %+v", p.config)
		if p.config == nil {
			return
		}
		var allErrors error

		for _, cfg := range p.config.GetConfigs() {
			p.logger.Debugf("processing middleware config: %s", cfg.Name)
			name := cmp.Or(cfg.Name, cfg.Type)
			if name == "" {
				continue
			}
			if _, exists := p.serverMiddlewares[name]; exists {
				continue
			}
			// Use the new WithServerCarrier option to pass only server middlewares
			opts := append(p.opts, runtimeMiddleware.WithServerCarrier(p.serverMiddlewares))
			// Attempt to create a server middleware. If the factory returns ok=false,
			// it means this config is not for a server middleware, so we just skip it.
			if sm, ok := runtimeMiddleware.NewServer(cfg, opts...); ok {
				p.serverMiddlewares[name] = sm
				p.logger.Debugf("registered server middleware '%s'", name)
			} else {
				// If NewServer returns false, it means this config was not for a server middleware.
				// We don't treat this as an error, just skip it.
				continue
			}
		}
		p.serverMWsErr = allErrors
	})

	p.mu.RLock()
	defer p.mu.RUnlock()
	return maps.Clone(p.serverMiddlewares), p.serverMWsErr
}

// Names returns the ordered list of middleware names.
func (p *Provider) Names() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.orderedNames
}

// ServerMiddleware returns a single server middleware instance by name.
func (p *Provider) ServerMiddleware(name string) (kratosMiddleware.Middleware, error) {
	middlewares, err := p.ServerMiddlewares()
	if err != nil {
		return nil, err
	}
	mw, ok := middlewares[name]
	if !ok {
		return nil, fmt.Errorf("server middleware '%s' not found", name)
	}
	return mw, nil
}
