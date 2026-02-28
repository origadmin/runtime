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
	"github.com/origadmin/runtime/bootstrap"
	runtimeconfig "github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/container"
	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/contracts"
	"github.com/origadmin/runtime/contracts/constant"
	"github.com/origadmin/runtime/contracts/options"
	runtimelog "github.com/origadmin/runtime/log"
)

// Standard configuration sources for interface sniffing.
type (
	AppConfig        interface{ GetApp() *appv1.App }
	LoggerConfig     interface{ GetLogger() *loggerv1.Logger }
	MiddlewareConfig interface {
		GetMiddlewares() *middlewarev1.Middlewares
	}
)

// App defines the application's runtime environment.
type App struct {
	appInfo       *appv1.App
	result        bootstrap.Result
	container     container.Container
	mu            sync.RWMutex
	globalOpts    []options.Option
	componentOpts map[string][]options.Option
	containerOpts []options.Option
	ctx           context.Context
	cancel        context.CancelFunc
}

// New creates a new App instance with default metadata.
func New(name, version string, opts ...Option) *App {
	return NewWithAppInfo(NewAppInfo(name, version), opts...)
}

// NewWithAppInfo creates a new App instance using a pre-configured App info.
func NewWithAppInfo(info *appv1.App, opts ...Option) *App {
	ctx, cancel := context.WithCancel(context.Background())
	if info == nil {
		info = NewAppInfoBuilder()
	}
	return configure.Apply(&App{
		appInfo:       info,
		componentOpts: make(map[string][]options.Option),
		containerOpts: make([]options.Option, 0),
		ctx:           ctx,
		cancel:        cancel,
	}, opts)
}

// Load loads and merges configuration metadata according to the defined structure.
func (r *App) Load(path string, bootOpts ...bootstrap.Option) error {
	res, err := bootstrap.New(path, bootOpts...)
	if err != nil {
		return fmt.Errorf("bootstrap failed: %w", err)
	}
	r.result = res

	// --- Metadata Refresh Flow ---

	// 1. [Infrastructure Layer] Initial metadata from framework bootstrap source.
	// This captures the base 'app' node defined in bootstrap.yaml.
	if boot := res.Bootstrap(); boot != nil && boot.GetApp() != nil {
		UpdateAppInfo(r.appInfo, boot.GetApp())
	}

	// 2. [Content Layer] Refresh metadata from Loader.
	// This ensures environment-specific overrides (e.g., config_dev.yaml) are captured.
	if loader := res.Loader(); loader != nil {
		var meta struct {
			App *appv1.App `json:"app" yaml:"app"`
		}
		if err := loader.Scan(&meta); err == nil && meta.App != nil {
			UpdateAppInfo(r.appInfo, meta.App)
		}
	}

	// 3. [Binding Layer] Sniff from explicitly bound business configuration objects.
	// This has the highest priority and overrides all previous values.
	biz := res.Config()
	if biz != nil {
		if p, ok := biz.(AppConfig); ok && p.GetApp() != nil {
			UpdateAppInfo(r.appInfo, p.GetApp())
		}
		if p, ok := biz.(MiddlewareConfig); ok && p.GetMiddlewares() != nil {
			r.containerOpts = append(r.containerOpts, container.WithMiddlewareConfig(p.GetMiddlewares()))
		}
		if p, ok := biz.(LoggerConfig); ok && p.GetLogger() != nil {
			r.containerOpts = append(r.containerOpts, container.WithLoggerConfig(p.GetLogger()))
		}
	}

	// 4. Validate core application metadata.
	if r.appInfo.GetName() == "" || r.appInfo.GetVersion() == "" {
		return errors.New("runtime: application metadata missing after load")
	}

	// 5. Initialize the container.
	ctnOpts := append(r.containerOpts, container.WithAppInfo(r.appInfo))
	r.container = container.New(res.StructuredConfig(), ctnOpts...)
	r.globalOpts = append(r.globalOpts, container.WithContainer(r.container), runtimelog.WithLogger(r.Logger()))
	return nil
}

// WarmUp initializes all registered components in the container.
func (r *App) WarmUp() error {
	var initErrors []error
	if _, err := r.RegistryProvider(); err != nil {
		initErrors = append(initErrors, err)
	}
	if _, err := r.DatabaseProvider(); err != nil {
		initErrors = append(initErrors, err)
	}
	if _, err := r.CacheProvider(); err != nil {
		initErrors = append(initErrors, err)
	}
	if _, err := r.ObjectStoreProvider(); err != nil {
		initErrors = append(initErrors, err)
	}
	if _, err := r.MiddlewareProvider(); err != nil {
		initErrors = append(initErrors, err)
	}
	r.mu.RLock()
	for name := range r.componentOpts {
		if name == string(constant.ComponentRegistries) || name == string(constant.ComponentDatabases) || name == string(constant.ComponentCaches) || name == string(constant.ComponentObjectStores) || name == string(constant.ComponentMiddlewares) {
			continue
		}
		if _, err := r.Component(name); err != nil {
			initErrors = append(initErrors, err)
		}
	}
	r.mu.RUnlock()
	if len(initErrors) > 0 {
		return fmt.Errorf("warm-up failed: %v", initErrors)
	}
	return nil
}

// AddGlobalOptions adds functional options to all components.
func (r *App) AddGlobalOptions(opts ...options.Option) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.globalOpts = append(r.globalOpts, opts...)
}

// AddComponentOptions adds functional options to a specific component.
func (r *App) AddComponentOptions(name string, opts ...options.Option) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.componentOpts[name] = append(r.componentOpts[name], opts...)
}

// Convenience methods for component configuration.
func (r *App) ConfigureRegistry(opts ...options.Option) {
	r.AddComponentOptions(string(constant.ComponentRegistries), opts...)
}
func (r *App) ConfigureDatabase(opts ...options.Option) {
	r.AddComponentOptions(string(constant.ComponentDatabases), opts...)
}
func (r *App) ConfigureCache(opts ...options.Option) {
	r.AddComponentOptions(string(constant.ComponentCaches), opts...)
}
func (r *App) ConfigureObjectStore(opts ...options.Option) {
	r.AddComponentOptions(string(constant.ComponentObjectStores), opts...)
}
func (r *App) ConfigureMiddleware(opts ...options.Option) {
	r.AddComponentOptions(string(constant.ComponentMiddlewares), opts...)
}

func (r *App) getMergedOptions(name string) []options.Option {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append(append([]options.Option{}, r.globalOpts...), r.componentOpts[name]...)
}

// Getters for internal state.
func (r *App) Config() runtimeconfig.KConfig                { return r.result.Loader() }
func (r *App) StructuredConfig() contracts.StructuredConfig { return r.result.StructuredConfig() }
func (r *App) Logger() log.Logger                           { return r.container.Logger() }
func (r *App) Container() container.Container               { return r.container }
func (r *App) Result() bootstrap.Result                     { return r.result }
func (r *App) Context() context.Context                     { return r.ctx }

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

// Component retrieval methods.
func (r *App) Component(name string) (interface{}, error) {
	return r.container.Component(name, r.getMergedOptions(name)...)
}
func (r *App) DefaultRegistrar() (registry.Registrar, error) { return r.container.DefaultRegistrar() }
func (r *App) RegistryProvider() (container.RegistryProvider, error) {
	return r.container.Registry(r.getMergedOptions(string(constant.ComponentRegistries))...)
}
func (r *App) DatabaseProvider() (container.DatabaseProvider, error) {
	return r.container.Database(r.getMergedOptions(string(constant.ComponentDatabases))...)
}
func (r *App) CacheProvider() (container.CacheProvider, error) {
	return r.container.Cache(r.getMergedOptions(string(constant.ComponentCaches))...)
}
func (r *App) ObjectStoreProvider() (container.ObjectStoreProvider, error) {
	return r.container.ObjectStore(r.getMergedOptions(string(constant.ComponentObjectStores))...)
}
func (r *App) MiddlewareProvider() (container.MiddlewareProvider, error) {
	return r.container.Middleware(r.getMergedOptions(string(constant.ComponentMiddlewares))...)
}

// AppInfo returns the current application metadata.
func (r *App) AppInfo() *appv1.App { return r.appInfo }

// ShowAppInfo prints application metadata to stdout.
func (r *App) ShowAppInfo() {
	ai := r.AppInfo()
	if ai == nil {
		return
	}
	ts := time.Now().Format(time.RFC3339)
	host, _ := os.Hostname()
	pid := os.Getpid()
	fmt.Printf("[%s] %s (pid:%d@%s)\n  Version: %s\n  AppId: %s\n  InstanceId: %s\n", ts, ai.Name, pid, host, ai.Version, ai.Id, ai.InstanceId)
}
