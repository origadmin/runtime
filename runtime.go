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

	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/config/envsetup"
	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/registry"
	"github.com/origadmin/toolkits/errors"
)

const (
	DefaultEnvPrefix = "ORIGADMIN_RUNTIME_SERVICE"
)

// build is a global variable that holds an instance of the builder struct.
var (
	globalRuntime = newRuntime()
)

// ErrNotFound is an error that is returned when a ConfigBuilder or RegistryBuilder is not found.
var ErrNotFound = errors.String("not found")

type Runtime interface {
	Signals() []os.Signal
	SetSignals([]os.Signal)
	Load(opts ...config.Option) error
	CreateApp(context.Context, registry.Registry, ...transport.Server) *kratos.App
}

type runtime struct {
	loaded    *atomic.Bool
	builder   Builder
	prefix    string
	signals   []os.Signal
	source    *configv1.SourceConfig
	bootstrap *bootstrap.Bootstrap
	logger    log.KLogger
	Logging   log.Logging
	options   *Options
	loader    *config.Loader
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

	opts = append(opts, config.WithServiceName(r.bootstrap.ServiceName()))
	if sourceConfig.Env {
		err := envsetup.SetWithPrefix(r.prefix, sourceConfig.EnvArgs)
		if err != nil {
			return errors.Wrap(err, "set env")
		}
		opts = append(opts, config.WithEnvPrefixes(sourceConfig.EnvPrefixes...))
	}

	r.source = sourceConfig
	r.loader = config.NewWithBuilder(r.builder.Config())
	if err := r.loader.Load(r.source, opts...); err != nil {
		return err
	}
	resolved, err := r.loader.GetResolved()
	if err != nil {
		return err
	}
	// Initialize the logs
	if err := r.initLogger(resolved.Logger()); err != nil {
		return errors.Wrap(err, "init logger")
	}

	r.loaded.Store(true)
	return nil
}

func (r *runtime) CreateApp(ctx context.Context, rr registry.Registry, ss ...transport.Server) *kratos.App {
	opts := buildServiceOptions(r.bootstrap.ServiceInfo())
	opts = append(opts,
		kratos.Context(ctx),
		kratos.Logger(r.logger),
		kratos.Signal(r.signals...),
	)

	if rr != nil {
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

func (r *runtime) reload(bs *bootstrap.Bootstrap, opts ...config.Option) error {
	r.loaded.Store(false)
	r.bootstrap = bs

	if err := r.Load(opts...); err != nil {
		return err
	}

	return nil
}

// Global function returns the interface type
func Global() Runtime {
	return globalRuntime
}

// New creates a new Runtime instance with default settings.
func New() Runtime {
	return newRuntime()
}

func newRuntime() *runtime {
	return &runtime{
		loaded:  new(atomic.Bool),
		builder: NewBuilder(),
		prefix:  DefaultEnvPrefix,
		signals: defaultSignals(),
	}
}

// Load uses the global Runtime instance to load configurations and other resources
// with the provided bootstrap settings. It returns an error if the loading process fails.
func Load(bs *bootstrap.Bootstrap, opts ...Option) (Runtime, error) {
	r := newRuntime()
	options := settings.ApplyZero(opts)
	if r.builder == nil {
		r.builder = NewBuilder()
	}
	if options.Prefix != "" {
		r.prefix = r.options.Prefix
	}
	if options.Logger != nil {
		r.logger = r.options.Logger
	}
	if len(options.Signals) > 0 {
		r.signals = r.options.Signals
	}

	if err := r.reload(bs, r.options.ConfigOptions...); err != nil {
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
