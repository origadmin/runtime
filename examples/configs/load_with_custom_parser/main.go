package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/config"
)

type CustomConfig struct {
	App struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"app"`
}

func main() {
	logger := log.NewStdLogger(os.Stdout)
	helper := log.NewHelper(logger)

	// Use custom transformer
	res, err := bootstrap.New("bootstrap.yaml",
		bootstrap.WithConfigTransformer(bootstrap.ConfigTransformFunc(TransformConfig)),
	)
	if err != nil {
		helper.Errorf("failed to bootstrap: %v", err)
		os.Exit(1)
	}

	cfg := res.Config().(*CustomConfig)
	fmt.Printf("App Name: %s, Version: %s\n", cfg.App.Name, cfg.App.Version)

	// Use custom loader operations
	loader := res.Loader()
	val, _ := loader.Value("app.name").String()
	fmt.Printf("Raw App Name from Loader: %s\n", val)

	ctx := context.Background()
	<-ctx.Done()
}

func TransformConfig(cfg config.KConfig) (any, error) {
	var c CustomConfig
	if err := cfg.Scan(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
