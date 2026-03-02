package main

import (
	"fmt"
	"log"

	rt "github.com/origadmin/runtime"
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
)

func main() {
	appInfo := rt.NewAppInfo("load-with-runtime-example", "1.0.0")
	rtInstance := rt.NewWithAppInfo(appInfo)
	err := rtInstance.Load("examples/configs/load_with_runtime/config/bootstrap.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	defer rtInstance.Config().Close()

	var appCfg appv1.App
	// Use KConfig native Value().Scan() instead of Decode
	if err := rtInstance.Config().Value("app").Scan(&appCfg); err != nil {
		log.Fatalf("Failed to scan app config: %v", err)
	}

	fmt.Printf("�?Successfully loaded config via runtime App.\n")
	fmt.Printf("   App Name: %s\n", appCfg.Name)
	fmt.Printf("   App Version: %s\n", appCfg.Version)
}
