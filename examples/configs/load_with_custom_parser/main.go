package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/origadmin/runtime"
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/container"
	"github.com/origadmin/runtime/contracts"
	"github.com/origadmin/runtime/contracts/options"
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
	config         contracts.ConfigLoader
	FeatureEnabled bool       `json:"feature_enabled"`
	APIKey         string     `json:"api_key"`
	RateLimit      int        `json:"rate_limit"`
	Endpoints      []Endpoint `json:"endpoints"`
}

func (c *CustomSettings) DecodedConfig() any {
	return c
}

func (c *CustomSettings) DecodeCaches() (*datav1.Caches, error) {
	return nil, errors.New("not implemented")
}

func (c *CustomSettings) DecodeDatabases() (*datav1.Databases, error) {
	return nil, errors.New("not implemented")
}

func (c *CustomSettings) DecodeObjectStores() (*datav1.ObjectStores, error) {
	return nil, errors.New("not implemented")
}

func (c *CustomSettings) DecodeData() (*datav1.Data, error) {
	return nil, errors.New("not implemented")
}

func (c *CustomSettings) DecodeDefaultDiscovery() (string, error) {
	return "", errors.New("not implemented")
}

func (c *CustomSettings) DecodeDiscoveries() (*discoveryv1.Discoveries, error) {
	return nil, errors.New("not implemented")
}

func (c *CustomSettings) DecodeServers() (*transportv1.Servers, error) {
	return nil, errors.New("not implemented")
}

func (c *CustomSettings) DecodeClients() (*transportv1.Clients, error) {
	return nil, errors.New("not implemented")
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

func (c *CustomSettings) DecodeMiddlewares() (*middlewarev1.Middlewares, error) {
	return nil, errors.New("not implemented")
}

func TransformConfig(cfg contracts.ConfigLoader) (any, error) {
	log.Infof("Loaded config: %+v", cfg)
	var settings CustomSettings
	settings.config = cfg
	if err := cfg.Decode("components.my-custom-settings", &settings); err != nil {
		log.Errorf("Failed to decode config: %v", err)
		return nil, err
	}
	return &settings, nil
}

func main() {
	DummyInit()

	rtInstance := runtime.New(
		"ApiGatewayCustomParserExample",
		"1.0.0",
		runtime.WithID("api_gateway_custom_parser_example"),
		runtime.WithEnv("dev"),
		runtime.WithStartTime(time.Now()),
		runtime.WithContainerOptions(
			container.WithComponentFactory("my--settings", container.ComponentFunc(
				func(cfg contracts.StructuredConfig, ctn container.Container, opts ...options.Option) (contracts.Component, error) {
					customCfg, ok := cfg.(*CustomSettings)
					if !ok {
						return nil, fmt.Errorf("expected *CustomSettings, but got %T", cfg)
					}
					return customCfg, nil
				})),
		),
	)

	err := rtInstance.Load("examples/configs/load_with_custom_parser/config/bootstrap.yaml",
		bootstrap.WithConfigTransformer(bootstrap.ConfigTransformFunc(TransformConfig)))
	if err != nil {
		panic(fmt.Errorf("failed to initialize runtime: %w", err))
	}

	logger := rtInstance.Logger()
	appLogger := log.NewHelper(logger)
	log.SetLogger(logger)
	appLogger.Info("Application started successfully!")

	comp, err := rtInstance.Component("my-custom-settings")
	if err != nil {
		appLogger.Error("Custom settings component not found")
		return
	}

	customSettings, ok := comp.(*CustomSettings)
	if !ok {
		appLogger.Error("Failed to type assert custom settings component")
		return
	}

	appLogger.Infof("Custom Settings: %+v", customSettings)
	appLogger.Infof("Feature Enabled: %t", customSettings.FeatureEnabled)
	appLogger.Infof("API Key: %s", customSettings.APIKey)
	appLogger.Infof("Rate Limit: %d", customSettings.RateLimit)
	for i, ep := range customSettings.Endpoints {
		appLogger.Infof("Endpoint %d: Name=%s, Path=%s", i, ep.Name, ep.Path)
	}

	config := rtInstance.Config()

	var bc conf.Bootstrap
	if err := config.Decode("servers", &bc.Servers); err != nil {
		appLogger.Errorf("Failed to decode servers config: %v", err)
		return
	}
	if err := config.Decode("clients", &bc.Clients); err != nil {
		appLogger.Errorf("Failed to decode clients config: %v", err)
		return
	}

	if len(bc.Servers) > 0 && bc.Servers[0].GetHttp() != nil {
		appLogger.Infof("Server HTTP Addr: %s", bc.Servers[0].GetHttp().GetAddr())
	} else {
		appLogger.Info("No HTTP server configuration found.")
	}

	for name, client := range bc.Clients {
		if client.GetClient() != nil {
			appLogger.Infof("Client '%s' Endpoint: %s", name, client.GetClient())
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	<-ctx.Done()
	appLogger.Info("Application finished.")
}
