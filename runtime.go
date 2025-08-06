package runtime

import (
	"context"
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goexts/generic/settings"
	"github.com/google/uuid"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/config/resolver"
	"github.com/origadmin/runtime/interfaces"
	// logpkg "github.com/origadmin/runtime/log" // Alias for runtime/log
	// middlewarepkg "github.com/origadmin/runtime/middleware" // Alias for runtime/middleware
	"github.com/origadmin/runtime/service"
)

// App is a kratos app.
type App struct {
	kratosApp      *kratos.App
	resolvedConfig interfaces.Resolved // Add resolvedConfig field
}

// NewApp creates a new kratos app.
func NewApp(ctx context.Context, srv *service.ServiceInfo, opts ...Option) (*App, error) {
	o := settings.Apply(&options{
		id:      uuid.New().String(),
		name:    srv.Name,
		version: srv.Version,
		sigs:    []os.Signal{os.Interrupt, os.Kill},
		ctx:     ctx,
		logger:  log.DefaultLogger,
	}, opts)

	// Initialize and load configuration
	var resolved interfaces.Resolved
	if o.configSource != nil {
		kcfg, err := config.New(o.configSource)
		if err != nil {
			return nil, fmt.Errorf("failed to create config: %w", err)
		}

		if err := kcfg.Load(); err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}

		resolved = resolver.New(kcfg)
	}

	// Initialize and set logger
	// if resolved != nil && resolved.Logger() != nil {
	// 	o.logger = logpkg.New(resolved.Logger())
	// 	log.SetLogger(o.logger)
	// }

	// Build server middlewares
	// var serverMiddlewares []kratosmiddleware.Middleware
	// if resolved != nil && resolved.Middleware() != nil {
	// 	serverMiddlewares = defaultManager.MiddlewareProvider.BuildServer(resolved.Middleware())
	// }

	// Apply middlewares to servers
	// for i, s := range o.servers {
	// 	if hs, ok := s.(*http.Server); ok {
	// 		hs.Use(serverMiddlewares...)
	// 	} else if gs, ok := s.(*grpc.Server); ok {
	// 		gs.Use(serverMiddlewares...)
	// 	}
	// 	o.servers[i] = s
	// }

	kratosApp := kratos.New(
		kratos.ID(o.id),
		kratos.Name(o.name),
		kratos.Version(o.version),
		kratos.Metadata(o.metadata),
		kratos.Server(o.servers...),
		kratos.Signal(o.sigs...),
		kratos.Context(o.ctx),
		kratos.Logger(o.logger),
		kratos.Registrar(o.registrar),
	)

	return &App{kratosApp: kratosApp, resolvedConfig: resolved}, nil
}

// Run starts the application and waits for a stop signal.
func (a *App) Run() error {
	return a.kratosApp.Run()
}

// Stop stops the application.
func (a *App) Stop() error {
	return a.kratosApp.Stop()
}

// GetResolvedConfig returns the resolved configuration.
func (a *App) GetResolvedConfig() interfaces.Resolved {
	return a.resolvedConfig
}

// Option is an app option.
type Option func(*options)

type options struct {
	id           string
	name         string
	version      string
	metadata     map[string]string
	servers      []transport.Server
	sigs         []os.Signal
	ctx          context.Context
	logger       log.Logger
	registrar    registry.Registrar
	configSource *configv1.SourceConfig
}

// WithID sets the app id.
func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

// WithVersion sets the app version.
func WithVersion(version string) Option {
	return func(o *options) {
		o.version = version
	}
}

// WithMetadata sets the app metadata.
func WithMetadata(md map[string]string) Option {
	return func(o *options) {
		o.metadata = md
	}
}

// WithServer sets the app servers.
func WithServer(srv ...transport.Server) Option {
	return func(o *options) {
		o.servers = srv
	}
}

// WithSignal sets the app signals.
func WithSignal(sigs ...os.Signal) Option {
	return func(o *options) {
		o.sigs = sigs
	}
}

// WithLogger sets the app logger.
func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// WithRegistrar sets the app registrar.
func WithRegistrar(r registry.Registrar) Option {
	return func(o *options) {
		o.registrar = r
	}
}

// WithConfigSource sets the config source for the app.
func WithConfigSource(cfg *configv1.SourceConfig) Option {
	return func(o *options) {
		o.configSource = cfg
	}
}
