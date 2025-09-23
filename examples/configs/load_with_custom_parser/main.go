package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/origadmin/runtime"
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
	FeatureEnabled bool   `json:"feature_enabled"`
	APIKey         string `json:"api_key"`
	RateLimit      int    `json:"rate_limit"`
	Endpoints      []Endpoint `json:"endpoints"`
}

// Register the component factory function for our custom settings
func init() {
	bootstrap.RegisterComponentFactory("custom_settings", func(cfg interfaces.Config, componentConfig map[string]interface{}) (interface{}, error) {
		// Create a new instance of CustomSettings
		settings := &CustomSettings{}

		// If config is provided, we can unmarshal it into our settings
		if componentConfig != nil {
			// Convert config to JSON and back to handle different config formats
			configBytes, err := json.Marshal(componentConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal config: %w", err)
			}

			if err := json.Unmarshal(configBytes, settings); err != nil {
				return nil, fmt.Errorf("failed to unmarshal config into CustomSettings: %w", err)
			}
		}

		return settings, nil
	})
}

func main() {
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
		bootstrap.WithAppInfo(appInfo), // Pass the AppInfo struct
	)
	if err != nil {
		panic(fmt.Errorf("failed to initialize runtime: %w", err))
	}
	defer cleanup()

	// Get logger from runtime
	logger := rt.Logger()
	appLogger := log.NewHelper(logger)
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

	// Decode servers and clients individually
	var bc conf.Bootstrap
	if err := config.Decode("servers", &bc.Servers); err != nil {
		appLogger.Errorf("Failed to decode servers config: %v", err)
		return
	}
	if err := config.Decode("clients", &bc.Clients); err != nil {
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
		if client.GetDiscovery() != nil {
			appLogger.Infof("Client '%s' Endpoint: %s", name, client.GetDiscovery().GetEndpoint())
		}
	}

	// Keep the application running
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Wait for interrupt signal
	<-ctx.Done()
	appLogger.Info("Application finished.")
}
