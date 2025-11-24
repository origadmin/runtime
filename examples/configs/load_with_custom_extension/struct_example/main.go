package main

import (
	"fmt"
	"log"

	rt "github.com/origadmin/runtime"
	"github.com/origadmin/runtime/config/protoutil"
	conf "github.com/origadmin/runtime/examples/protos/custom_extension"
)

func main() {
	// Use the new functional options pattern for AppInfo creation
	appInfo := rt.NewAppInfo("custom-extension-example", "1.0.0", rt.WithAppInfoEnv("development"))

	// --- 1. Load Configuration ---
	// We use Kratos config to load the YAML file.
	// Initialize runtime with bootstrap configuration
	rtInstance, err := rt.NewFromBootstrap(
		"examples/configs/load_with_custom_extension/config/bootstrap.yaml",
		rt.WithAppInfo(appInfo),
	)
	if err != nil {
		log.Fatalf("Failed to initialize runtime: %v", err)
	}
	defer rtInstance.Config().Close()

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

	// --- 3. Convert Struct to Strongly-Typed Config ---
	// Now, we use our helper function to convert the generic struct
	// into our specific, strongly-typed CustomAuthConfig.
	// Using protoutil.NewFromStruct for a clean conversion from google.protobuf.Struct
	cfg, err := protoutil.NewFromStruct[conf.CustomAuthConfig](mw.Customize)
	if err != nil {
		log.Fatalf("Failed to create typed config from customize struct: %v", err)
	}
	fmt.Println("\n✅ Successfully created typed config using protoutil.NewFromStruct.")
	fmt.Printf("   Policy: %s\n", cfg.Policy)
	fmt.Printf("   Required Scope: %s\n", cfg.RequiredScope)

	// Now you can use the `cfg` object with full type safety.
	if cfg.RequiredScope == "read:users" {
		fmt.Println("\n👍 Logic check passed: Configuration is ready to be used by the auth middleware.")
	}
}
