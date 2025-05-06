/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package bootstrap

import (
	"os"

	"github.com/origadmin/toolkits/errors"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type Loader interface {
	Bootstrap() *Bootstrap
	Load() (*configv1.SourceConfig, error)
}

type loader struct {
	config    *configv1.SourceConfig
	bootstrap *Bootstrap
}

func (l *loader) Bootstrap() *Bootstrap {
	return l.bootstrap
}

// Load loading configuration force reload
func (l *loader) Load() (*configv1.SourceConfig, error) {
	cfg, err := LoadSourceConfig(l.bootstrap)
	if err != nil {
		return nil, err
	}
	l.config = cfg
	return cfg, nil
}

// NewLoader creates a new loader instance.
func NewLoader(bootstrap *Bootstrap) Loader {
	return &loader{
		bootstrap: bootstrap,
	}
}

// loadSourceConfig loads the config file from the given path
func loadSourceConfig(path string) (*configv1.SourceConfig, error) {
	var cfg configv1.SourceConfig
	err := decodeFile(path, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// LoadSourceConfig loads the config file from the given path
func LoadSourceConfig(bootstrap *Bootstrap) (*configv1.SourceConfig, error) {
	// Get the path from the bootstrap
	path := bootstrap.ConfigFilePath()

	// Get the file info from the path
	info, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get config file %s", path)
	}
	if info.IsDir() {
		return nil, errors.New("config path is a directory")
	}
	return loadSourceConfig(path)
}

// LoadSourceConfigFromPath loads the config file from the given path
func LoadSourceConfigFromPath(path string) (*configv1.SourceConfig, error) {
	// Get the file info from the path
	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "load config stat error")
	}
	if stat.IsDir() {
		return nil, errors.New("config path is a directory")
	}
	// Load the config file
	return loadSourceConfig(path)
}
