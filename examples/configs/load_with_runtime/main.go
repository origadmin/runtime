package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/origadmin/runtime"
	appv1 "github.com/origadmin/runtime/api/gen/go/runtime/app/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"

	discoveryv1 "github.com/origadmin/runtime/api/gen/go/runtime/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/runtime/logger/v1"
	// Import the generated Go code from the load_with_runtime proto definition.
	conf "github.com/origadmin/runtime/examples/protos/load_with_runtime"
	"github.com/origadmin/runtime/log" // Import the log package
)

// ProtoConfig is a custom implementation of interfaces.Config that handles Protobuf decoding.
type ProtoConfig struct {
	ifconfig  interfaces.Config // Keep for Raw() and Close()
	bootstrap *conf.Bootstrap   // The fully decoded protobuf config
}

func (d *ProtoConfig) Load() error {
	return d.ifconfig.Load()
}

func (d *ProtoConfig) DecodeApp() (*appv1.App, error) {
	return nil, errors.New("not implemented")
}

func (d *ProtoConfig) DecodeMiddleware() (*middlewarev1.Middlewares, error) {
	return nil, errors.New("not implemented")
}

// NewProtoConfig creates a new ProtoConfig instance and decodes the entire Kratos config into the Bootstrap proto.
func NewProtoConfig(c interfaces.Config) (*ProtoConfig, error) {
	var bc conf.Bootstrap
	if err := c.Decode("", &bc); err != nil {
		return nil, err
	}
	return &ProtoConfig{
		ifconfig:  c,
		bootstrap: &bc,
	}, nil
}

// Decode implements the interfaces.Config interface.
func (d *ProtoConfig) Decode(key string, target interface{}) error {
	var source interface{} // This will hold the Go representation of the config section

	if key == "" || key == "/" {
		source = d.bootstrap
	} else {
		// For specific top-level keys, extract the corresponding field from the bootstrap proto.
		// Note: GetX() methods return nil if the field is not set.
		switch key {
		case "servers":
			source = d.bootstrap.GetServers()
		case "endpoints": // Changed from "clients" to "endpoints"
			source = d.bootstrap.GetEndpoints()
		default:
			// If the key is not a direct top-level field of Bootstrap, and not empty/root,
			// we can try to scan it from the original Kratos config as a fallback.
			// This handles cases where some config might not be explicitly defined in the proto.
			return d.ifconfig.Decode(key, target)
		}
	}

	if source == nil {
		return fmt.Errorf("config key '%s' not found or is nil in bootstrap proto", key)
	}

	// Now, marshal the extracted 'source' (which is a Go struct/slice/map, potentially a proto.Message)
	// to JSON and then unmarshal into the target. This is the most flexible way
	// to handle generic `interface{}` targets, especially when dealing with protobuf messages.
	jsonBytes, err := json.Marshal(source)
	if err != nil {
		return fmt.Errorf("failed to marshal source data for key '%s' to JSON: %w", key, err)
	}

	if pm, ok := target.(proto.Message); ok {
		// If target is a proto.Message, use protojson.Unmarshal
		if err := protojson.Unmarshal(jsonBytes, pm); err != nil {
			return fmt.Errorf("failed to protojson unmarshal JSON to %T for key '%s': %w", pm, key, err)
		}
	} else {
		// If target is not a proto.Message (e.g., map[string]interface{}, struct), use standard json.Unmarshal
		if err := json.Unmarshal(jsonBytes, target); err != nil {
			return fmt.Errorf("failed to unmarshal JSON to %T for key '%s': %w", target, key, err)
		}
	}

	return nil
}

// Raw implements the interfaces.Config interface.
func (d *ProtoConfig) Raw() any {
	return d.ifconfig
}

// Close implements the interfaces.Config interface.
func (d *ProtoConfig) Close() error {
	return d.ifconfig.Close()
}

// DecodeLogger implements the interfaces.LoggerConfigDecoder interface.
func (d *ProtoConfig) DecodeLogger() (*loggerv1.Logger, error) {
	// The log config is directly available in d.bootstrap
	if d.bootstrap.GetLogger() == nil {
		return nil, errors.New("logger config not found in bootstrap proto")
	}

	return d.bootstrap.GetLogger(), nil
}

// DecodeDiscoveries implements the interfaces.DiscoveriesConfigDecoder interface.
// It now only decodes the centralized discovery providers from the registries configuration.
func (d *ProtoConfig) DecodeDiscoveries() (map[string]*discoveryv1.Discovery, error) {
	if d.bootstrap.GetRegistries() == nil {
		return make(map[string]*discoveryv1.Discovery), nil
	}
	return d.bootstrap.GetRegistries().GetDiscoveries(), nil
}

// DecodeEndpoints decodes and links endpoint configurations with their corresponding discovery providers.
// It returns a map where the key is the endpoint's name and the value contains all necessary
// information (both behavior and provider) to initialize the client connection.
func (d *ProtoConfig) DecodeEndpoints() (map[string]*discoveryv1.Endpoint, error) {
	// 1. Get all available discovery providers
	providers := d.bootstrap.GetRegistries().GetDiscoveries()
	if providers == nil {
		providers = make(map[string]*discoveryv1.Discovery)
	}

	// 2. Prepare for the final results
	resolvedEndpoints := make(map[string]*discoveryv1.Endpoint)

	// 3. Iterate through all defined EndpointConfig from the bootstrap file
	for endpointName, endpointCfg := range d.bootstrap.GetEndpoints() {
		if endpointCfg == nil {
			continue
		}

		// 4. Find the provider
		var providerCfg *discoveryv1.Discovery
		providerName := endpointCfg.GetDiscoveryName()

		if providerName != "" {
			var found bool
			providerCfg, found = providers[providerName]
			if !found {
				return nil, fmt.Errorf("endpoint '%s' references non-existent discovery provider '%s'", endpointName, providerName)
			}
		}
		_ = providerCfg
		// 5. Assemble the final, "rich" Endpoint object
		resolvedEndpoints[endpointName] = &discoveryv1.Endpoint{
			//Provider: providerCfg, // Link the found provider
			Uri: endpointCfg.GetUri(),
			Selector: &discoveryv1.Selector{ // Manual conversion for Selector
				Type:    endpointCfg.GetSelector().GetType(),
				Version: endpointCfg.GetSelector().GetVersion(),
			},
			//Timeout:   endpointCfg.GetTimeout(),
			//Transport: endpointCfg.GetTransport(), // Direct assignment if types match
			//Middlewares: endpointCfg.GetMiddlewares(),
		}
	}

	return resolvedEndpoints, nil
}

func main() {
	// Define the ConfigTransformFunc to create our custom ProtoConfig.
	configTransformer := bootstrap.ConfigTransformFunc(func(kc interfaces.Config) (interfaces.StructuredConfig, error) {
		protoCfg, err := NewProtoConfig(kc)
		if err != nil {
			return nil, fmt.Errorf("failed to create ProtoConfig: %w", err)
		}
		return protoCfg, nil
	})

	// 1. Create a new Runtime instance from the new bootstrap config.
	//    Path is now relative to the CWD (runtime directory), pointing to the bootstrap.yaml.
	rt, cleanup, err := runtime.NewFromBootstrap(
		"examples/configs/load_with_runtime/config/bootstrap.yaml", // Correctly load bootstrap.yaml
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "rich-config-runtime-example",
			Name:    "RichConfigRuntimeExample",
			Version: "1.0.0",
			Env:     "dev",
		}),
		bootstrap.WithConfigTransformer(configTransformer), // Inject custom transformer
	)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// Get the configured logger from the runtime instance
	appLogger := log.NewHelper(rt.Logger()) // Use log.NewHelper for convenience

	appLogger.Info("Application started successfully!") // Log a message using the configured logger

	// 2. Get the configuration decoder from the runtime instance and assert it to our ProtoConfig type
	decoder := rt.Config()

	// Type assert the decoder to our ProtoConfig
	protoCfg, ok := decoder.(*ProtoConfig)
	if !ok {
		appLogger.Error("Failed to assert config to ProtoConfig type")
		panic("config decoder is not of type *ProtoConfig")
	}

	// 3. Get the bootstrap config directly from our ProtoConfig
	bc := protoCfg.bootstrap
	if bc == nil {
		appLogger.Error("Bootstrap config is nil")
		panic("bootstrap config is nil")
	}

	// 4. Print the loaded configuration to verify
	appLogger.Info("--- Loaded Rich Config via runtime interface ---")

	// Verify logger config
	if loggerCfg := bc.GetLogger(); loggerCfg != nil {
		appLogger.Infof("Logger Level: %s, Format: %s, Stdout: %v",
			loggerCfg.GetLevel(),
			loggerCfg.GetFormat(),
			loggerCfg.GetStdout())

		if fileCfg := loggerCfg.GetFile(); fileCfg != nil {
			appLogger.Infof("Log File: %s, MaxSize: %dMB, MaxBackups: %d, MaxAge: %dd, Compress: %v",
				fileCfg.GetPath(),
				fileCfg.GetMaxSize(),
				fileCfg.GetMaxBackups(), fileCfg.GetMaxAge(),
				fileCfg.GetCompress())
		}
	} else {
		appLogger.Info("No logger configuration found")
	}

	// Verify server configs
	servers := bc.GetServers()
	if len(servers) > 0 {
		appLogger.Infof("Found %d server configurations", len(servers))
		for i, srv := range servers {
			if httpSrv := srv.GetHttp(); httpSrv != nil {
				appLogger.Infof("Server[%d] HTTP: Addr=%s, Network=%s, Timeout=%s",
					i, httpSrv.GetAddr(), httpSrv.GetNetwork(), httpSrv.GetTimeout().AsDuration())
			}
			if grpcSrv := srv.GetGrpc(); grpcSrv != nil {
				appLogger.Infof("Server[%d] gRPC: Addr=%s, Network=%s, Timeout=%s",
					i, grpcSrv.GetAddr(), grpcSrv.GetNetwork(), grpcSrv.GetTimeout().AsDuration())
			}
		}
	} else {
		appLogger.Info("No server configurations found")
	}

	// Verify resolved endpoints
	resolvedEndpoints, err := protoCfg.DecodeEndpoints()
	if err != nil {
		appLogger.Errorf("Failed to decode endpoints: %v", err)
		panic(err)
	}

	if len(resolvedEndpoints) > 0 {
		appLogger.Infof("Found %d resolved endpoints", len(resolvedEndpoints))
		for name, endpoint := range resolvedEndpoints {
			appLogger.Infof("Endpoint '%s':", name)

			if endpoint != nil {
				appLogger.Infof("  Provider: Type=%s, ServiceName=%s",
					endpoint.GetName(),
					endpoint.GetDiscoveryName())
			} else {
				appLogger.Info("  Provider: None (static endpoint or missing discovery_name)")
			}

			appLogger.Infof("  URI: %s", endpoint.GetUri())

			if selector := endpoint.GetSelector(); selector != nil {
				appLogger.Infof("  Selector: Type=%s, Version=%s",
					selector.GetType(),
					selector.GetVersion())
			}

		}
	} else {
		appLogger.Info("No resolved endpoints found")
	}

	// Verify registries (discovery providers)
	discoveries := bc.GetRegistries().GetDiscoveries()
	if len(discoveries) > 0 {
		appLogger.Infof("Found %d raw discovery configurations", len(discoveries))
		for name, disc := range discoveries {
			appLogger.Infof("Discovery '%s': Type=%s, ServiceName=%s",
				name, disc.GetType(), disc.GetName())
		}
	} else {
		appLogger.Info("No raw discovery configurations found in registries")
	}

	appLogger.Info("Application finished.") // Log a final message
}
