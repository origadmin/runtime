package middleware

import (
	"cmp"
	"errors"
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	kratosMiddleware "github.com/go-kratos/kratos/v2/middleware"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/interfaces/options"
	runtimeMiddleware "github.com/origadmin/runtime/middleware"
)

// Provider implements interfaces.ClientMiddlewareProvider and interfaces.MiddlewareProvider
type Provider struct {
	config                *middlewarev1.Middlewares
	log                   *log.Helper
	opts                  []options.Option // Now stores options passed to SetConfig
	clientMiddlewares     map[string]kratosMiddleware.Middleware
	serverMiddlewares     map[string]kratosMiddleware.Middleware
	onceClientMiddlewares sync.Once
	onceServerMiddlewares sync.Once
}

func (p *Provider) ServerMiddleware(name string) (kratosMiddleware.Middleware, error) {
	ms, err := p.ServerMiddlewares()
	if err != nil {
		return nil, err
	}
	sm, ok := ms[name]
	if !ok {
		return nil, fmt.Errorf("server middleware '%s' not found", name)
	}
	return sm, nil
}

func (p *Provider) RegisterServerMiddleware(name string, middleware kratosMiddleware.Middleware) {
	p.serverMiddlewares[name] = middleware
}

func (p *Provider) ClientMiddleware(name string) (kratosMiddleware.Middleware, error) {
	cm, err := p.ClientMiddlewares()
	if err != nil {
		return nil, err
	}
	cmw, ok := cm[name]
	if !ok {
		return nil, fmt.Errorf("client middleware '%s' not found", name)
	}
	return cmw, nil
}

func (p *Provider) RegisterClientMiddleware(name string, middleware kratosMiddleware.Middleware) {
	p.clientMiddlewares[name] = middleware
}

// SetConfig sets the middleware configurations and dynamic options for the provider.
func (p *Provider) SetConfig(cfg *middlewarev1.Middlewares, opts ...options.Option) *Provider {
	p.config = cfg
	p.opts = opts // Store the dynamically passed options
	return p
}

func (p *Provider) ClientMiddlewares() (map[string]kratosMiddleware.Middleware, error) {
	var allErrors error
	p.onceClientMiddlewares.Do(func() {
		for _, cfg := range p.config.GetConfigs() {
			name := cmp.Or(cfg.Name, cfg.Type)
			// Pass the stored options to the client middleware creation
			cm, ok := runtimeMiddleware.NewClientMiddleware(cfg, p.opts...)
			if !ok {
				allErrors = errors.Join(allErrors, fmt.Errorf("failed to create client middleware '%s'", name))
				continue
			}
			p.clientMiddlewares[name] = cm
		}
	})
	return p.clientMiddlewares, allErrors
}

func (p *Provider) ServerMiddlewares() (map[string]kratosMiddleware.Middleware, error) {
	var allErrors error
	p.onceServerMiddlewares.Do(func() {
		for _, cfg := range p.config.GetConfigs() {
			name := cmp.Or(cfg.Name, cfg.Type)
			// Pass the stored options to the server middleware creation
			sm, ok := runtimeMiddleware.NewServerMiddleware(cfg, p.opts...)
			if !ok {
				allErrors = errors.Join(allErrors, fmt.Errorf("failed to create server middleware '%s'", name))
				continue
			}
			p.serverMiddlewares[name] = sm
		}
	})
	return p.serverMiddlewares, allErrors
}

// NewProvider creates a new Provider.
// It no longer receives opts, as options are passed dynamically via SetConfig.
func NewProvider(logger log.Logger) *Provider {
	helper := log.NewHelper(logger)
	return &Provider{
		log:               helper,
		clientMiddlewares: make(map[string]kratosMiddleware.Middleware),
		serverMiddlewares: make(map[string]kratosMiddleware.Middleware),
	}
}
