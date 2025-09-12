package interfaces

import (
	"os"

	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/log"
	serviceoptions "github.com/origadmin/runtime/service/options" // Re-import service options
)

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
	Load(bs *bootstrap.Bootstrap, opts ...serviceoptions.Option) error
	Run() error
	Stop() error
	WithLoggerAttrs(kvs ...any) Runtime
}
