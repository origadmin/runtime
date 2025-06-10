/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime provides functions for loading configurations and registering services.
package runtime

import (
	"os"
	"sync/atomic"
	"syscall"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goexts/generic/settings"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/config/envsetup"
	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/registry"
	"github.com/origadmin/toolkits/errors"
)

const (
	DefaultEnvPrefix = "ORIGADMIN_RUNTIME_SERVICE"
)

// build is a global variable that holds an instance of the builder struct.
var (
	//globalRuntime  = newRuntime()
	runtimeBuilder = NewBuilder()
)

// ErrNotFound is an error that is returned when a ConfigBuilder or RegistryBuilder is not found.
var ErrNotFound = errors.String("not found")

type Logger interface {
	Logger() log.KLogger
	SetLogger(kvs ...any)
	WithLogger(kvs ...any) log.KLogger
}

type SignalHandler interface {
	Signals() []os.Signal
	SetSignals([]os.Signal)
}

type Runtime interface {
	Logger
	SignalHandler
	Client() Runtime
	Builder() Builder
	Context() context.Context
	Load(opts ...config.Option) error
	CreateApp(...transport.Server) *kratos.App
	WithLoggerAttrs(kvs ...any) Runtime
	SetRegistry(registrar registry.KRegistrar)
}

type runtime struct {
	ctx       context.Context
	loaded    *atomic.Bool
	prefix    string
	signals   []os.Signal
	source    *configv1.SourceConfig
	bootstrap *bootstrap.Bootstrap
	logger    log.KLogger
	loader    *config.Loader
	registrar registry.KRegistrar
	client    bool
}

func (r *runtime) SetRegistry(registrar registry.KRegistrar) {
	r.registrar = registrar
}

func (r *runtime) Builder() Builder {
	return runtimeBuilder
}

func (r *runtime) Context() context.Context {
	return r.ctx
}

func (r *runtime) Client() Runtime {
	rr := *r
	rr.client = true
	return &rr
}

func (r *runtime) WithLoggerAttrs(kvs ...any) Runtime {
	rr := *r
	rr.logger = log.With(rr.logger, kvs...)
	return &rr
}

func (r *runtime) Logger() log.KLogger {
	if r.logger == nil {
		r.logger = log.DefaultLogger
	}
	return r.logger
}

func (r *runtime) SetLogger(kvs ...any) {
	r.logger = log.With(r.logger, kvs...)
}

func (r *runtime) WithLogger(kvs ...any) log.KLogger {
	return log.With(r.Logger(), kvs...)
}

func (r *runtime) Signals() []os.Signal {
	return r.signals
}

func (r *runtime) SetSignals(signals []os.Signal) {
	r.signals = signals
}

func (r *runtime) IsLoaded() bool {
	return r.loaded.Load()
}

func (r *runtime) Load(opts ...config.Option) error {
	if r.IsLoaded() {
		return nil
	}
	sourceConfig, err := bootstrap.LoadSourceConfig(r.bootstrap.ConfigFilePath())
	if err != nil {
		return errors.Wrap(err, "load source config")
	}
	log.NewHelper(log.GetLogger()).Infof("loading config: %+v", sourceConfig)
	opts = append(opts, config.WithServiceName(r.bootstrap.ServiceName()))
	if sourceConfig.Env {
		err := envsetup.SetWithPrefix(r.prefix, sourceConfig.EnvArgs)
		if err != nil {
			return errors.Wrap(err, "set env")
		}
		opts = append(opts, config.WithEnvPrefixes(sourceConfig.EnvPrefixes...))
	}

	if err := r.loader.Load(sourceConfig, opts...); err != nil {
		return err
	}
	resolved, err := r.loader.GetResolved()
	if err != nil {
		return err
	}
	// Initialize the logs
	if r.logger == nil {
		if err := r.initLogger(resolved.Logger()); err != nil {
			return errors.Wrap(err, "init logger")
		}
	}

	r.loaded.Store(true)
	return nil
}

func (r *runtime) buildRegistrar() (registry.KRegistrar, error) {
	if r.client {
		return nil, nil
	}
	resolved, err := r.loader.GetResolved()
	if err != nil {
		return nil, err
	}
	registrar, err := runtimeBuilder.NewRegistrar(resolved.Discovery())
	if err != nil {
		return nil, err
	}
	return registrar, nil
}

func (r *runtime) Resolve(fn func(kConfig config.KConfig) error) error {
	if fn == nil {
		return errors.New("resolve function is nil")
	}
	if err := r.Load(); err != nil {
		return err
	}
	source, err := r.loader.GetSource()
	if err != nil {
		return err
	}
	return fn(source)
}

func (r *runtime) CreateApp(ss ...transport.Server) *kratos.App {
	opts := buildServiceOptions(r.bootstrap.ServiceInfo())
	opts = append(opts,
		kratos.Context(r.ctx),
		kratos.Logger(r.WithLogger("module", "server")),
		kratos.Signal(r.signals...),
	)
	rr, err := r.buildRegistrar()
	if err != nil {
		_ = r.WithLogger("module", "runtime").Log(log.LevelError, "create registrar failed", err)
	} else if rr != nil {
		opts = append(opts, kratos.Registrar(rr))
	}

	if len(ss) > 0 {
		opts = append(opts, kratos.Server(ss...))
	}

	return kratos.New(opts...)
}

func buildServiceOptions(info bootstrap.ServiceInfo) []kratos.Option {
	return []kratos.Option{
		kratos.ID(info.ID),
		kratos.Name(info.Name),
		kratos.Version(info.Version),
		kratos.Metadata(info.Metadata),
	}
}

func (r *runtime) initLogger(loggingCfg *configv1.Logger) error {
	if loggingCfg == nil {
		return errors.New("logger config is nil")
	}

	r.logger = log.New(loggingCfg)
	return nil
}

func (r *runtime) reload(bs *bootstrap.Bootstrap, opts []config.Option) error {
	r.loaded.Store(false)
	r.bootstrap = bs

	if err := r.Load(opts...); err != nil {
		return err
	}

	return nil
}

// Global function returns the interface type
func Global() Runtime {
	return newRuntime()
}

// GlobalBuilder returns the global Builder instance.
func GlobalBuilder() Builder {
	return runtimeBuilder
}

func newRuntime() *runtime {
	return &runtime{
		ctx:     context.Background(),
		loaded:  new(atomic.Bool),
		prefix:  DefaultEnvPrefix,
		signals: defaultSignals(),
	}
}

// Load uses the global Runtime instance to load configurations and other resources
// with the provided bootstrap settings. It returns an error if the loading process fails.
func Load(bs *bootstrap.Bootstrap, opts ...Option) (Runtime, error) {
	r := newRuntime()
	options := settings.ApplyZero(opts)

	if options.Context != nil {
		r.ctx = options.Context
	}
	if options.Prefix != "" {
		r.prefix = options.Prefix
	}
	if options.Logger != nil {
		r.logger = options.Logger
	}
	if len(options.Signals) > 0 {
		r.signals = options.Signals
	}
	if options.Resolver != nil {
		//r.resolver = options.Resolver
		r.loader = config.NewWithBuilder(runtimeBuilder.Config())
		if err := r.loader.SetResolver(options.Resolver); err != nil {
			return nil, err
		}
	}

	if err := r.reload(bs, options.ConfigOptions); err != nil {
		return nil, err
	}
	return r, nil
}

func defaultSignals() []os.Signal {
	return []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	}
}
