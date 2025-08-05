package runtime

import (
	"context"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goexts/generic/settings"
	"github.com/google/uuid"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/service"
)

// App is a kratos app.
type App struct {
	kratosApp    *kratos.App
	configLoader *config.Loader
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
	var loader *config.Loader
	if o.configSource != nil {
		loader = config.New() // Use the default builder
		if err := loader.Load(o.configSource); err != nil {
			return nil, err
		}
	}

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

	return &App{kratosApp: kratosApp, configLoader: loader}, nil
}

// Run starts the application and waits for a stop signal.
func (a *App) Run() error {
	return a.kratosApp.Run()
}

// Stop stops the application.
func (a *App) Stop() error {
	return a.kratosApp.Stop()
}

// GetConfigLoader returns the config loader.
func (a *App) GetConfigLoader() *config.Loader {
	return a.configLoader
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
