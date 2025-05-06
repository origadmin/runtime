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
	ignores   []string
}

func (l *loader) Bootstrap() *Bootstrap {
	return l.bootstrap
}

// Load 加载配置（强制重新加载）
func (l *loader) Load() (*configv1.SourceConfig, error) {
	cfg, err := LoadSourceConfig(l.bootstrap)
	if err != nil {
		return nil, err
	}

	l.config = cfg
	return cfg, nil
}

// NewLoader 创建Loader实例
func NewLoader(bootstrap *Bootstrap, ) Loader {
	return &loader{
		bootstrap: bootstrap,
	}
}

// loadSourceConfig loads the config file from the given path
func loadSourceConfig(si os.FileInfo, path string, ignores []string) (*configv1.SourceConfig, error) {
	// Check if the file or directory exists
	if si == nil {
		return nil, errors.New("load config file target is not exist")
	}
	var cfg configv1.SourceConfig
	err := loadCustomizeConfig(si, path, &cfg, ignores)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// loadCustomizeConfig loads the user config file from the given path
func loadCustomizeConfig(si os.FileInfo, path string, cfg any, ignores []string) error {
	// Check if the path is a directory
	decode := decodeFile
	if si.IsDir() {
		decode = decodeDir
	}
	err := decode(path, cfg, ignores)
	if err != nil {
		return err
	}
	return nil
}

// LoadSourceConfig loads the config file from the given path
func LoadSourceConfig(bootstrap *Bootstrap, ignores []string) (*configv1.SourceConfig, error) {
	// Get the path from the bootstrap
	path := bootstrap.ConfigFilePath()

	// Get the file info from the path
	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "load config stat error")
	}

	// Load the config file
	return loadSourceConfig(stat, path, ignores)
}

// LoadSourceConfigFromPath loads the config file from the given path
func LoadSourceConfigFromPath(path string, ignores []string) (*configv1.SourceConfig, error) {
	// Get the file info from the path
	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "load config stat error")
	}
	// Load the config file
	return loadSourceConfig(stat, path, ignores)
}
