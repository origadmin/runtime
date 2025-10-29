package main

import (
	"fmt"
	"log"

	rt "github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
	conf "github.com/origadmin/runtime/examples/protos/custom_extension"
	"github.com/origadmin/runtime/extension/customize"
	"github.com/origadmin/runtime/interfaces"
)

func main() {
	// --- 1. Load Configuration ---
	// We use Kratos config to load the YAML file.
	// Initialize runtime with bootstrap configuration
	rtInstance, err := rt.NewFromBootstrap(
		"examples/configs/load_with_custom_extension/config/bootstrap.yaml",
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "custom-extension-example",
			Name:    "Custom Extension Example",
			Version: "1.0.0",
			Env:     "development",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to initialize runtime: %v", err)
	}
	defer rtInstance.Cleanup()

	// Get config decoder
	decoder := rtInstance.Config()

	// --- 2. Scan into Middleware Proto ---
	// Kratos automatically handles unmarshaling the 'customize' YAML block
	// into the google.protobuf.Struct field.
	// Process extensions
	var mw conf.MiddlewareStruct
	if err := decoder.Decode("", &mw); err != nil {
		log.Fatalf("Failed to decode configuration: %v", err)
	}

	fmt.Println("✅ Successfully loaded config into Middleware with generic Struct.")
	fmt.Printf("   Middleware Name: %s\n", mw.Name)
	fmt.Printf("   Raw Customize Struct: %v\n", mw.Customize)

	var authCfgDirect conf.CustomAuthConfig
	if err := decoder.Decode("customize", &authCfgDirect); err != nil {
		log.Fatalf("Failed to decode customize config: %v", err)
	}
	fmt.Printf("   Direct Customize Config.")
	fmt.Printf("   Policy: %s\n", authCfgDirect.Policy)
	fmt.Printf("   Required Scope: %s\n", authCfgDirect.RequiredScope)

	// --- 3. Convert Struct to Strongly-Typed Config ---
	// Now, we use our helper function to convert the generic struct
	// into our specific, strongly-typed CustomAuthConfig.
	var authCfg conf.CustomAuthConfig
	if err := customize.UnmarshalTo(mw.Customize, &authCfg); err != nil {
		log.Fatalf("Failed to get typed config from struct: %v", err)
	}

	fmt.Println("\n✅ Successfully converted generic Struct to strongly-typed CustomAuthConfig.")
	fmt.Printf("   Policy: %s\n", authCfg.Policy)
	fmt.Printf("   Required Scope: %s\n", authCfg.RequiredScope)

	cfg, err := customize.NewFromStruct[conf.CustomAuthConfig](mw.Customize)
	if err != nil {
		log.Fatalf("Failed to create typed config: %v", err)
	}
	fmt.Println("\n✅ Successfully created typed config using NewTypedConfig.")
	fmt.Printf("   Policy: %s\n", cfg.Policy)
	fmt.Printf("   Required Scope: %s\n", cfg.RequiredScope)

	// Now you can use the `authCfg` object with full type safety.
	if authCfg.RequiredScope == "read:users" {
		fmt.Println("\n👍 Logic check passed: Configuration is ready to be used by the auth middleware.")
	}
}
