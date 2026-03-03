package runtime

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goexts/generic/configure"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	"github.com/origadmin/runtime/engine/bootstrap"
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/engine"
	enginecontext "github.com/origadmin/runtime/engine/context"
	"github.com/origadmin/runtime/engine/metadata"
)

// Standard configuration interfaces (Contract Re-export)
type (
	AppConfig        = component.AppConfig
	LoggerConfig     = component.LoggerConfig
	MiddlewareConfig = component.MiddlewareConfig
	DataConfig       = component.DataConfig
	RegistryConfig   = component.RegistryConfig
)

// App defines the application's runtime environment powered by engine.
type App struct {
	appInfo *appv1.App
	result  bootstrap.Result
	engine  component.Registry
	ctx     context.Context
	cancel  context.CancelFunc
}

// New creates a new App instance.
func New(name, version string, opts ...Option) *App {
	return NewWithAppInfo(NewAppInfo(name, version), opts...)
}

// NewWithAppInfo creates a new App instance using a pre-configured App info.
func NewWithAppInfo(info *appv1.App, opts ...Option) *App {
	ctx, cancel := context.WithCancel(context.Background())
	if info == nil {
		info = NewAppInfoBuilder()
	}

	// Create engine registry at startup
	reg := engine.NewContainer()

	app := &App{
		appInfo: info,
		engine:  reg,
		ctx:     ctx,
		cancel:  cancel,
	}

	// Apply options
	app = configure.Apply(app, opts)

	// Apply framework DEFAULTS last
	app.registerDefaultFactories()

	return app
}

// WithRegistry sets a callback to configure the internal engine registry.
func WithRegistry(fn func(component.Registry)) Option {
	return func(a *App) {
		fn(a.engine)
	}
}

func (r *App) registerDefaultFactories() {
	// Logger Default
	if !r.engine.Has(metadata.CategoryLogger) {
		r.engine.Register(metadata.CategoryLogger,
			DefaultLoggerExtractor,
			DefaultLoggerProvider,
			engine.WithPriority(metadata.PriorityInfrastructure))
	}

	// Registry Default
	if !r.engine.Has(metadata.CategoryRegistry) {
		r.engine.Register(metadata.CategoryRegistry,
			DefaultRegistryExtractor,
			DefaultRegistryProvider,
			engine.WithPriority(metadata.PriorityRegistry))
	}
}

// Load loads configuration into Result.
func (r *App) Load(path string, bootOpts ...bootstrap.Option) error {
	res, err := bootstrap.New(path, bootOpts...)
	if err != nil {
		return fmt.Errorf("bootstrap failed: %w", err)
	}
	r.result = res

	// Refresh app info
	if boot := res.Bootstrap(); boot != nil && boot.GetApp() != nil {
		UpdateAppInfo(r.appInfo, boot.GetApp())
	}
	if loader := res.Loader(); loader != nil {
		var meta struct {
			App *appv1.App `json:"app" yaml:"app"`
		}
		if err := loader.Scan(&meta); err == nil && meta.App != nil {
			UpdateAppInfo(r.appInfo, meta.App)
		}
	}
	if biz := res.Config(); biz != nil {
		if p, ok := biz.(component.AppConfig); ok && p.GetApp() != nil {
			UpdateAppInfo(r.appInfo, p.GetApp())
		}
	}

	if r.appInfo.GetName() == "" || r.appInfo.GetVersion() == "" {
		return errors.New("runtime: application metadata missing after load")
	}

	return nil
}

// WarmUp activates the engine with the loaded configuration.
func (r *App) WarmUp() error {
	if r.result == nil || r.result.Config() == nil {
		return errors.New("runtime: cannot warm-up without loaded configuration")
	}
	return r.engine.Init(r.ctx, r.result.Config())
}

// Getters
func (r *App) Config() runtimeconfig.KConfig { return r.result.Loader() }
func (r *App) BusinessConfig() any           { return r.result.Config() }
func (r *App) Logger() log.Logger {
	// Navigate via CategoryLogger
	l, err := engine.Cast[log.Logger](r.ctx, r.engine.In(metadata.CategoryLogger), "logger")
	if err != nil {
		return log.DefaultLogger
	}
	return l
}
func (r *App) Result() bootstrap.Result      { return r.result }
func (r *App) Container() component.Registry { return r.engine }

// Context returns the app context.
func (r *App) Context() context.Context { return r.ctx }

// NewContext creates a new context from the app context.
func (r *App) NewContext(ctx context.Context) context.Context {
	return enginecontext.NewContext(ctx)
}

// NewTrace creates a new context with the given trace ID.
func (r *App) NewTrace(ctx context.Context, traceID string) context.Context {
	return enginecontext.NewTrace(ctx, traceID)
}

func (r *App) Stop() {
	if r.cancel != nil {
		r.cancel()
	}
}

func (r *App) NewApp(servers []transport.Server, options ...kratos.Option) *kratos.App {
	info := r.appInfo
	md := info.GetMetadata()
	if md == nil {
		md = make(map[string]string)
	}
	if info.GetEnv() != "" {
		md["env"] = info.GetEnv()
	}
	opts := []kratos.Option{
		kratos.Context(r.ctx),
		kratos.Logger(r.Logger()),
		kratos.Server(servers...),
		kratos.ID(info.GetId()),
		kratos.Name(info.GetName()),
		kratos.Version(info.GetVersion()),
		kratos.Metadata(md),
	}
	if registrar, _ := r.DefaultRegistrar(); registrar != nil {
		opts = append(opts, kratos.Registrar(registrar))
	}
	opts = append(opts, options...)
	return kratos.New(opts...)
}

func (r *App) DefaultRegistrar() (registry.Registrar, error) {
	return engine.Cast[registry.Registrar](r.ctx, r.engine.In(metadata.CategoryRegistry), "")
}

func (r *App) AppInfo() *appv1.App { return r.appInfo }

func (r *App) ShowAppInfo() {
	ai := r.appInfo
	if ai == nil {
		return
	}
	ts := time.Now().Format(time.RFC3339)
	host, _ := os.Hostname()
	pid := os.Getpid()
	fmt.Printf("[%s] %s (pid:%d@%s)\n  Version: %s\n  AppId: %s\n  InstanceId: %s\n", ts, ai.Name, pid, host, ai.Version, ai.Id, ai.InstanceId)
}
