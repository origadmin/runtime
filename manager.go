package runtime

import (
	"path/filepath"

	kratosconfig "github.com/go-kratos/kratos/v2/config" // Add kratosconfig

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1" // Add configv1
	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/config/file" // Add config/file
	"github.com/origadmin/runtime/interfaces"
)

// Manager manages all builders and providers for the runtime.
type Manager struct {
	ConfigBuilder interfaces.ConfigBuilder
	//MiddlewareProvider interfaces.MiddlewareProvider
	//ServerBuilder      interfaces.ServerBuilder   // Renamed from ServiceProvider
	//RegistryBuilder    interfaces.RegistryBuilder // Added RegistryBuilder
	// Add other builders/providers here as needed
}

// defaultManager is the package-level Manager instance.
var defaultManager *Manager

func init() {
	// Initialize the defaultManager
	defaultManager = &Manager{
		ConfigBuilder: config.NewBuilder(),
		//MiddlewareProvider: middleware.NewBuilder(),
		//ServerBuilder:      service.NewBuilder(),
		//RegistryBuilder:    registry.NewBuilder(),
	}

	// Register default config factories
	defaultManager.ConfigBuilder.Register("file", interfaces.FileConfig(func(sourceConfig *configv1.SourceConfig,
		opts *interfaces.Options) (kratosconfig.Source, error) {
		cfg := sourceConfig.GetFile()
		if cfg == nil {
			return nil, config.ErrInvalidConfigType
		}
		var options []file.FileOption
		if len(cfg.Ignores) > 0 {
			options = append(options, file.WithIgnores(cfg.Ignores...))
		}
		//v := new(configs.Bootstrap)
		//options = append(options, file.WithFormatter(fileFormatter(v)))
		path, _ := filepath.Abs(cfg.Path)
		//log.NewHelper(log.GetLogger()).Infof("loading config from %s", path)
		return file.NewSource(path, options...), nil
	}))

	// Register other default factories for middleware, service, etc.
	// No need to register middleware, service, registry builders here as they are directly assigned.
}

// GetManager returns the default runtime manager.
func GetManager() *Manager {
	return defaultManager
}

// RegisterConfigFactory registers a config factory with the default manager's ConfigBuilder.
func RegisterConfigFactory(name string, factory interfaces.ConfigFactory) {
	defaultManager.ConfigBuilder.RegisterConfigFunc(name, factory)
}

// RegisterMiddlewareProvider registers a middleware provider with the default manager.
func RegisterMiddlewareProvider(provider interfaces.MiddlewareProvider) {
	defaultManager.MiddlewareProvider = provider
}

// RegisterServerBuilder registers a server builder with the default manager.
func RegisterServerBuilder(builder interfaces.ServerBuilder) {
	defaultManager.ServerBuilder = builder
}

// RegisterRegistryBuilder registers a registry builder with the default manager.
func RegisterRegistryBuilder(builder interfaces.RegistryBuilder) {
	defaultManager.RegistryBuilder = builder
}

// RegisterConfigSync registers a config syncer with the default manager's ConfigBuilder.
func RegisterConfigSync(name string, syncFunc config.Syncer) {
	defaultManager.ConfigBuilder.RegisterConfigSyncer(name, syncFunc)
}
