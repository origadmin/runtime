package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/origadmin/runtime"
	appv1 "github.com/origadmin/runtime/api/gen/go/runtime/app/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/runtime/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/runtime/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/log"

	// Import the generated Go code from the api_gateway proto definition.
	conf "github.com/origadmin/runtime/examples/protos/api_gateway"
)

// Endpoint represents a single API endpoint configuration.
type Endpoint struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// CustomSettings represents the structure of our custom configuration section.
type CustomSettings struct {
	config         interfaces.Config
	FeatureEnabled bool       `json:"feature_enabled"`
	APIKey         string     `json:"api_key"`
	RateLimit      int        `json:"rate_limit"`
	Endpoints      []Endpoint `json:"endpoints"`
}

func (c *CustomSettings) Load() error {
	return errors.New("already loaded")
}

func (c *CustomSettings) Decode(key string, value any) error {
	switch key {
	case "feature_enabled":
		if value, ok := value.(*bool); ok {
			*value = c.FeatureEnabled
			return nil
		}
	case "api_key":
		if value, ok := value.(*string); ok {
			*value = c.APIKey
			return nil
		}
	case "rate_limit":
		if value, ok := value.(*int); ok {
			*value = c.RateLimit
			return nil
		}
	case "endpoints":
		if value, ok := value.(*[]Endpoint); ok {
			*value = c.Endpoints
			return nil
		}
	}
	return c.config.Decode(key, value)
}

func (c *CustomSettings) Raw() any {
	return c.config
}

func (c *CustomSettings) Close() error {
	return c.config.Close()
}

func (c *CustomSettings) DecodeApp() (*appv1.App, error) {
	return nil, errors.New("not implemented")
}

func (c *CustomSettings) DecodeLogger() (*loggerv1.Logger, error) {
	return nil, errors.New("not implemented")
}

func (c *CustomSettings) DecodeDiscoveries() (map[string]*discoveryv1.Discovery, error) {
	return nil, errors.New("not implemented")
}

func (c *CustomSettings) DecodeMiddleware() (*middlewarev1.Middlewares, error) {
	return nil, errors.New("not implemented")
}

func TransformConfig(cfg interfaces.Config) (interfaces.StructuredConfig, error) {
	//if err := cfg.Load(); err != nil {
	//	return nil, err
	//}
	// Create a new instance of CustomSettings
	log.Infof("Loaded config: %+v", cfg)
	var settings CustomSettings
	//log.Infof("Decoded config: %v", settingMap)
	settings.config = cfg
	var settingMap map[string]any
	if err := cfg.Decode("", &settingMap); err != nil {
		log.Errorf("Failed to decode config: %v", err)
		return nil, err
	}
	log.Infof("Decoded config: %v", settingMap)

	return &settings, nil
}

func main() {
	// Call DummyInit to ensure the local_registry package's init() function is executed.
	DummyInit()

	// Create AppInfo struct
	appInfo := interfaces.AppInfo{
		ID:        "api_gateway_custom_parser_example",
		Name:      "ApiGatewayCustomParserExample",
		Version:   "1.0.0",
		Env:       "dev",
		StartTime: time.Now(),
		Metadata:  make(map[string]string),
	}

	// 1. Create a new Runtime instance from the bootstrap config
	rt, cleanup, err := runtime.NewFromBootstrap(
		"examples/configs/load_with_custom_parser/config/bootstrap.yaml",
		bootstrap.WithAppInfo(&appInfo), // Pass the AppInfo struct
		bootstrap.WithConfigTransformer(bootstrap.ConfigTransformFunc(TransformConfig)),
		bootstrap.WithComponent("my-custom-settings", func(cfg interfaces.StructuredConfig, container interfaces.Container) (interface{}, error) {
			cfg, ok := cfg.(*CustomSettings)
			if ok {
				fmt.Printf("Custom Settings: %+v\n", cfg)
			}
			return cfg, nil
		}),
	)
	if err != nil {
		panic(fmt.Errorf("failed to initialize runtime: %w", err))
	}
	defer cleanup()

	// Get logger from runtime
	logger := rt.Logger()
	appLogger := log.NewHelper(logger)
	log.SetLogger(logger)
	appLogger.Info("Application started successfully!")

	// 2. Get the custom settings component
	comp, ok := rt.Component("my-custom-settings") // Updated to use the instance name
	if !ok {
		appLogger.Error("Custom settings component not found")
		return
	}

	// Type assert to our custom settings type
	customSettings, ok := comp.(*CustomSettings)
	if !ok {
		appLogger.Error("Failed to type assert custom settings component")
		return
	}

	// Log the custom settings
	appLogger.Infof("Custom Settings: %+v", customSettings)
	appLogger.Infof("Feature Enabled: %t", customSettings.FeatureEnabled)
	appLogger.Infof("API Key: %s", customSettings.APIKey)
	appLogger.Infof("Rate Limit: %d", customSettings.RateLimit)
	for i, ep := range customSettings.Endpoints {
		appLogger.Infof("Endpoint %d: Name=%s, Path=%s", i, ep.Name, ep.Path)
	}

	// 3. Get the config interface to decode other configurations
	config := rt.Config()

	// Decode servers and clients directly using config.Decode
	var bc conf.Bootstrap
	if err := config.Decode("servers", &bc.Servers); err != nil { // Direct decode
		appLogger.Errorf("Failed to decode servers config: %v", err)
		return
	}
	if err := config.Decode("clients", &bc.Clients); err != nil { // Direct decode
		appLogger.Errorf("Failed to decode clients config: %v", err)
		return
	}

	// Log the server configuration
	if len(bc.Servers) > 0 && bc.Servers[0].GetHttp() != nil {
		appLogger.Infof("Server HTTP Addr: %s", bc.Servers[0].GetHttp().GetAddr())
	} else {
		appLogger.Info("No HTTP server configuration found.")
	}

	// Log client configurations
	for name, client := range bc.Clients {
		if client.GetEndpoint() != nil {
			appLogger.Infof("Client '%s' Endpoint: %s", name, client.GetEndpoint())
		}
	}

	// Keep the application running
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Wait for interrupt signal
	<-ctx.Done()
	appLogger.Info("Application finished.")
}
