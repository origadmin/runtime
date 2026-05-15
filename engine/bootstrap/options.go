package bootstrap

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"

	sourcev1 "github.com/origadmin/runtime/api/gen/go/config/source/v1"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/helpers/optionutil"
)

type ProviderOptions struct {
	configTransformer ConfigTransformer
	configTarget      any
	directly          bool
	extraSources      []*sourcev1.SourceConfig
	providerOptions   []Option
	rawOptions        []kratosconfig.Option
	frameworkOptions  []options.Option
	config            any
	pathResolver      func(string) string
	prefixes          []string
}

type Option = options.Option

func WithConfigTransformer(transformer ConfigTransformer) Option {
	return optionutil.Update(func(opt *ProviderOptions) {
		opt.configTransformer = transformer
	})
}

func WithConfigTarget(target any) Option {
	return optionutil.Update(func(opt *ProviderOptions) {
		opt.configTarget = target
	})
}

func WithDirectly(directly bool) Option {
	return optionutil.Update(func(opt *ProviderOptions) {
		opt.directly = directly
	})
}

func WithProviderOptions(opts ...Option) Option {
	return optionutil.Update(func(opt *ProviderOptions) {
		opt.providerOptions = append(opt.providerOptions, opts...)
	})
}

func WithRawOptions(opts ...kratosconfig.Option) Option {
	return optionutil.Update(func(opt *ProviderOptions) {
		opt.rawOptions = append(opt.rawOptions, opts...)
	})
}

func WithFrameworkOptions(opts ...options.Option) Option {
	return optionutil.Update(func(opt *ProviderOptions) {
		opt.frameworkOptions = append(opt.frameworkOptions, opts...)
	})
}

func WithConfig(cfg any) Option {
	return optionutil.Update(func(opt *ProviderOptions) {
		opt.config = cfg
	})
}

func WithPathResolver(fn func(string) string) Option {
	return optionutil.Update(func(opt *ProviderOptions) {
		opt.pathResolver = fn
	})
}

func WithPrefixes(prefixes ...string) Option {
	return optionutil.Update(func(opt *ProviderOptions) {
		opt.prefixes = prefixes
	})
}

// WithEnvSource appends an env source to the extra sources list.
// It is a shortcut for WithExtraSources({Type: "env"}).
// Typically used together with WithDirectly to enable environment variable injection.
func WithEnvSource() Option {
	return optionutil.Update(func(opt *ProviderOptions) {
		opt.extraSources = append(opt.extraSources, &sourcev1.SourceConfig{Type: "env"})
	})
}

func FromOptions(opts ...Option) *ProviderOptions {
	return optionutil.NewT[ProviderOptions](opts...)
}
