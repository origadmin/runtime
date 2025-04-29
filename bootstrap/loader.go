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
	Load() (*configv1.SourceConfig, error)
	Reload() error
}

type decoder = func(string, any) error

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
func LoadSourceConfig(bootstrap *Bootstrap) (*configv1.SourceConfig, error) {
	// Get the path from the bootstrap
	path := bootstrap.ConfigFilePath()

	// Get the file info from the path
	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "load config stat error")
	}

	// Load the config file
	return loadSourceConfig(stat, path)
}

// LoadSourceConfigFromPath loads the config file from the given path
func LoadSourceConfigFromPath(path string) (*configv1.SourceConfig, error) {
	// Get the file info from the path
	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "load config stat error")
	}
	// Load the config file
	return loadSourceConfig(stat, path)
}

// LoadLocalConfig loads the config file from the given path
func LoadLocalConfig(bs *Bootstrap, v any) error {
	source, err := LoadSourceConfig(bs)
	if err != nil {
		return err
	}
	if source.GetType() != "file" {
		return errors.New("local config type must be file")
	}

	path := source.GetFile().GetPath()
	// Get the file info from the path
	stat, err := os.Stat(path)
	if err != nil {
		return errors.Wrap(err, "load config stat error")
	}

	return loadCustomizeConfig(stat, path, v)
}
