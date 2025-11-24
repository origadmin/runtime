package middleware

import (
	"cmp"
	"errors"
	"fmt"
	"sync"

	kratosMiddleware "github.com/go-kratos/kratos/v2/middleware"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/interfaces/options"
	runtimelog "github.com/origadmin/runtime/log"
	runtimeMiddleware "github.com/origadmin/runtime/middleware"
)

// Provider implements interfaces.ClientMiddlewareProvider and interfaces.MiddlewareProvider.
// It manages the lifecycle of middleware instances, caching them after first creation and
// allowing for reconfiguration. It is safe for concurrent use.
type Provider struct {
	mu                   sync.Mutex
	config               *middlewarev1.Middlewares
	log                  *runtimelog.Helper // Changed to runtimelog.Helper
	opts                 []options.Option
	clientMiddlewares    map[string]kratosMiddleware.Middleware
	serverMiddlewares    map[string]kratosMiddleware.Middleware
	clientMWsInitialized bool
	serverMWsInitialized bool
}

// NewProvider creates a new Provider.
func NewProvider(logger runtimelog.Logger) *Provider {
	return &Provider{
		log: runtimelog.NewHelper(logger), // Changed to runtimelog.NewHelper
	}
}

// SetConfig updates the provider's configuration. This will clear any previously
// cached instances and cause them to be recreated on the next access, using the new configuration.
func (p *Provider) SetConfig(cfg *middlewarev1.Middlewares, opts ...options.Option) *Provider {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.config = cfg
	p.opts = opts
	p.clientMWsInitialized = false
	p.serverMWsInitialized = false
	p.clientMiddlewares = make(map[string]kratosMiddleware.Middleware)
	p.serverMiddlewares = make(map[string]kratosMiddleware.Middleware)

	return p
}

// RegisterClientMiddleware allows for manual registration of a client middleware instance.
func (p *Provider) RegisterClientMiddleware(name string, middleware kratosMiddleware.Middleware) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clientMiddlewares[name] = middleware
}

// ClientMiddlewares returns a map of all available client middleware instances.
// On first call, it creates instances from the configuration and caches them.
// Subsequent calls return the cached instances unless SetConfig has been called.
func (p *Provider) ClientMiddlewares() (map[string]kratosMiddleware.Middleware, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.clientMWsInitialized {
		return p.clientMiddlewares, nil
	}

	var allErrors error
	if p.config != nil {
		for _, cfg := range p.config.GetConfigs() {
			name := cmp.Or(cfg.Name, cfg.Type)
			if _, exists := p.clientMiddlewares[name]; exists {
				p.log.Warnf("client middleware '%s' is already registered, skipping config-based creation", name)
				continue
			}
			cm, ok := runtimeMiddleware.NewClientMiddleware(cfg, p.opts...)
			if ok {
				p.clientMiddlewares[name] = cm
			} else {
				allErrors = errors.Join(allErrors, fmt.Errorf("failed to create client middleware '%s'", name))
			}
		}
	}

	p.clientMWsInitialized = true
	return p.clientMiddlewares, allErrors
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
func (p *Provider) RegisterServerMiddleware(name string, middleware kratosMiddleware.Middleware) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.serverMiddlewares[name] = middleware
}

// ServerMiddlewares returns a map of all available server middleware instances.
// It follows the same caching and creation logic as ClientMiddlewares.
func (p *Provider) ServerMiddlewares() (map[string]kratosMiddleware.Middleware, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.serverMWsInitialized {
		return p.serverMiddlewares, nil
	}

	var allErrors error
	if p.config != nil {
		for _, cfg := range p.config.GetConfigs() {
			name := cmp.Or(cfg.Name, cfg.Type)
			if _, exists := p.serverMiddlewares[name]; exists {
				p.log.Warnf("server middleware '%s' is already registered, skipping config-based creation", name)
				continue
			}
			sm, ok := runtimeMiddleware.NewServerMiddleware(cfg, p.opts...)
			if ok {
				p.serverMiddlewares[name] = sm
			} else {
				allErrors = errors.Join(allErrors, fmt.Errorf("failed to create server middleware '%s'", name))
			}
		}
	}

	p.serverMWsInitialized = true
	return p.serverMiddlewares, allErrors
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
