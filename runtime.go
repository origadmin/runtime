package runtime

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goexts/generic/configure"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/contracts"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/engine"
	"github.com/origadmin/runtime/engine/bootstrap"
	enginecontext "github.com/origadmin/runtime/engine/context"
	"github.com/origadmin/runtime/helpers/comp"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/registry"
)

// App defines the application's runtime environment powered by engine.
type App struct {
	appInfo *appv1.App
	result  bootstrap.Result
	engine  component.Container
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

	// Create engine registry at startup with standard resolvers and global registrations
	reg := engine.NewContainer(
		engine.WithCategoryResolvers(DefaultResolvers),
		engine.WithGlobalRegistrations(),
	)

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

// WithContainer sets a callback to configure the internal engine container.
func WithContainer(fn func(component.Container)) Option {
	return func(a *App) {
		fn(a.engine)
	}
}

// Register adds a component registration to the engine.
func Register(cat Category, p Provider, opts ...RegisterOption) {
	engine.Register(cat, p, opts...)
}

func (r *App) registerDefaultFactories() {
	// Logger Default
	r.engine.Register(CategoryLogger,
		log.DefaultProvider,
		engine.WithPriority(component.PriorityFramework))

	// Registry components are self-registered by the registry package init()
}

// Load loads configuration into Result.
func (r *App) Load(path string, bootOpts ...bootstrap.Option) error {
	res, err := bootstrap.New(path, bootOpts...)
	if err != nil {
		return fmt.Errorf("bootstrap failed: %w", err)
	}
	r.result = res

	// 1. Foundation: Bootstrap Metadata (Base layer)
	if boot := res.Bootstrap(); boot != nil && boot.GetApp() != nil {
		UpdateAppInfo(r.appInfo, boot.GetApp())
	}

	// 2. Override: Business Object (High priority)
	if biz := res.Config(); biz != nil {
		if p, ok := biz.(contracts.AppConfig); ok && p.GetApp() != nil {
			UpdateAppInfo(r.appInfo, p.GetApp())
		} else {
			// 3. Fallback: Scan from Decoder (if not strong-typed)
			if loader := res.Decoder(); loader != nil {
				var meta struct {
					App *appv1.App `json:"app" yaml:"app"`
				}
				if err := loader.Scan(&meta); err == nil && meta.App != nil {
					UpdateAppInfo(r.appInfo, meta.App)
				}
			}
		}
	}

	if r.appInfo.GetName() == "" || r.appInfo.GetVersion() == "" {
		return errors.New("runtime: application metadata missing after load")
	}

	// Auto warm-up the engine if business configuration is available
	if r.Config() != nil {
		if err := r.WarmUp(); err != nil {
			return fmt.Errorf("warm-up failed during load: %w", err)
		}
	}

	return nil
}

// WarmUp activates the engine with the loaded configuration.
func (r *App) WarmUp() error {
	if r.result == nil || r.result.Config() == nil {
		return errors.New("runtime: cannot warm-up without loaded configuration")
	}
	return r.engine.Load(r.ctx, r.result.Config())
}

// Getters
func (r *App) Decoder() runtimeconfig.KConfig { return r.result.Decoder() }
func (r *App) Config() any                    { return r.result.Config() }
func (r *App) Logger() log.Logger {
	l, err := comp.Get[log.Logger](r.ctx, r.engine.In(CategoryLogger))
	if err != nil {
		return log.DefaultLogger
	}
	return l
}
func (r *App) Result() bootstrap.Result       { return r.result }
func (r *App) Container() component.Container { return r.engine }

func (r *App) In(cat Category, opts ...InOption) component.Registry {
	return r.engine.In(cat, opts...)
}

// Context returns the app context.
func (r *App) Context() context.Context { return r.ctx }

// NewContext creates a new context from the app context.
func NewContext(ctx context.Context) context.Context {
	return enginecontext.NewContext(ctx)
}

// NewTrace creates a new context with the given trace ID.
func NewTrace(ctx context.Context, traceID string) context.Context {
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

func (r *App) DefaultRegistrar() (registry.KRegistrar, error) {
	// Directly obtain from CategoryRegistrar with standard Kratos interface
	return comp.Get[registry.KRegistrar](r.ctx, r.engine.In(CategoryRegistrar))
}

func (r *App) Discoveries() (map[string]registry.KDiscovery, error) {
	return registry.GetDiscoveries(r.ctx, r.engine.In(CategoryDiscovery))
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
	fmt.Printf("[%s] %s\n", ts, ai.Name)
	fmt.Printf("  PID:        %d@%s\n", pid, host)
	fmt.Printf("  Version:    %s\n", ai.Version)
	fmt.Printf("  AppId:      %s\n", ai.Id)
	fmt.Printf("  InstanceId: %s\n", ai.InstanceId)
}

// --- Wire Providers ---

// ProvideLogger is a Wire provider function that extracts the logger from the App.
func ProvideLogger(rt *App) log.Logger {
	return rt.Logger()
}

// ProvideDefaultRegistrar is a Wire provider function that extracts the registrar from the App.
func ProvideDefaultRegistrar(rt *App) (registry.KRegistrar, error) {
	return rt.DefaultRegistrar()
}
