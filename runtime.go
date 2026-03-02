package runtime

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goexts/generic/configure"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	"github.com/origadmin/runtime/bootstrap"
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/engine"
	runtimelog "github.com/origadmin/runtime/log"
)

// Standard configuration interfaces
type (
	AppConfig        interface{ GetApp() *appv1.App }
	LoggerConfig     interface{ GetLogger() *loggerv1.Logger }
	MiddlewareConfig interface {
		GetMiddlewares() *middlewarev1.Middlewares
	}
	DataConfig     interface{ GetData() *datav1.Data }
	RegistryConfig interface {
		GetDiscoveries() *discoveryv1.Discoveries
	}
)

// App defines the application's runtime environment powered by engine.
type App struct {
	appInfo    *appv1.App
	result     bootstrap.Result
	engine     engine.Registry
	mu         sync.RWMutex
	globalOpts []options.Option
	ctx        context.Context
	cancel     context.CancelFunc
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
	reg := engine.NewContainer(nil)
	app := &App{
		appInfo: info,
		engine:  reg,
		ctx:     ctx,
		cancel:  cancel,
	}
	app.registerCoreFactories()
	return configure.Apply(app, opts)
}

func (r *App) registerCoreFactories() {
	// Logger (Infrastructure)
	r.engine.Register(engine.CategoryInfrastructure, func(root any) (*engine.ModuleConfig, error) {
		if p, ok := root.(LoggerConfig); ok && p.GetLogger() != nil {
			l := p.GetLogger()
			return &engine.ModuleConfig{Entries: []engine.ConfigEntry{{Name: "logger", Value: l}}, Active: "logger"}, nil
		}
		return nil, nil
	}, func(ctx context.Context, h engine.Handle, opts ...options.Option) (any, error) {
		var cfg loggerv1.Logger
		if err := engine.BindConfig(h, &cfg); err != nil {
			return runtimelog.DefaultLogger, nil
		}
		return runtimelog.NewLogger(&cfg), nil
	}, engine.WithPriority(engine.PriorityInfrastructure))

	// Registry (Registry)
	r.engine.Register(engine.CategoryRegistry, func(root any) (*engine.ModuleConfig, error) {
		if p, ok := root.(RegistryConfig); ok && p.GetDiscoveries() != nil {
			raw := p.GetDiscoveries()
			var entries []engine.ConfigEntry
			for _, c := range raw.Configs {
				entries = append(entries, engine.ConfigEntry{Name: c.Name, Value: c})
			}
			return &engine.ModuleConfig{Entries: entries, Active: raw.GetActive()}, nil
		}
		return nil, nil
	}, func(ctx context.Context, h engine.Handle, opts ...options.Option) (any, error) {
		return nil, errors.New("registry provider not implemented")
	}, engine.WithPriority(engine.PriorityRegistry))

	// Middleware (Stack)
	r.engine.Register(engine.CategoryMiddleware, func(root any) (*engine.ModuleConfig, error) {
		if p, ok := root.(MiddlewareConfig); ok && p.GetMiddlewares() != nil {
			raw := p.GetMiddlewares()
			var entries []engine.ConfigEntry
			for _, c := range raw.Configs {
				entries = append(entries, engine.ConfigEntry{Name: c.Name, Value: c})
			}
			return &engine.ModuleConfig{Entries: entries}, nil
		}
		return nil, nil
	}, func(ctx context.Context, h engine.Handle, opts ...options.Option) (any, error) {
		return nil, errors.New("middleware provider not implemented")
	}, engine.WithPriority(engine.PriorityServerStack), engine.WithScope(engine.ServerScope))
}

// Load loads and merges configuration metadata.
func (r *App) Load(path string, bootOpts ...bootstrap.Option) error {
	res, err := bootstrap.New(path, bootOpts...)
	if err != nil {
		return fmt.Errorf("bootstrap failed: %w", err)
	}
	r.result = res
	r.engine.BindRoot(res.Config())

	// --- Metadata Refresh Flow ---

	// 1. Initial metadata from framework bootstrap source.
	if boot := res.Bootstrap(); boot != nil && boot.GetApp() != nil {
		UpdateAppInfo(r.appInfo, boot.GetApp())
	}

	// 2. Refresh metadata from Loader (raw KConfig).
	// This captures environment-specific overrides (e.g., config_dev.yaml).
	if loader := res.Loader(); loader != nil {
		var meta struct {
			App *appv1.App `json:"app" yaml:"app"`
		}
		if err := loader.Scan(&meta); err == nil && meta.App != nil {
			UpdateAppInfo(r.appInfo, meta.App)
		}
	}

	// 3. Sniff from business configuration object.
	if biz := res.Config(); biz != nil {
		if p, ok := biz.(AppConfig); ok && p.GetApp() != nil {
			UpdateAppInfo(r.appInfo, p.GetApp())
		}
	}

	if r.appInfo.GetName() == "" || r.appInfo.GetVersion() == "" {
		return errors.New("runtime: application metadata missing after load")
	}

	return nil
}

// WarmUp initializes standard components.
func (r *App) WarmUp() error {
	return r.engine.Init(r.ctx)
}

// Getters
func (r *App) Config() runtimeconfig.KConfig { return r.result.Loader() }
func (r *App) BusinessConfig() any           { return r.result.Config() }
func (r *App) Logger() log.Logger {
	l, err := engine.Cast[log.Logger](r.ctx, r.engine.In(engine.CategoryInfrastructure), "logger")
	if err != nil {
		return runtimelog.DefaultLogger
	}
	return l
}
func (r *App) Result() bootstrap.Result { return r.result }
func (r *App) Context() context.Context { return r.ctx }

// Stop cancels the root context to shut down the application.
func (r *App) Stop() {
	if r.cancel != nil {
		r.cancel()
	}
}

// NewApp creates a native Kratos application instance.
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
	return engine.Cast[registry.Registrar](r.ctx, r.engine.In(engine.CategoryRegistry), "")
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
