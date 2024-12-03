/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package bootstrap

import (
	"os"
	"path/filepath"

	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/errors"

	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/config"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type decoder = func(string, any) error

// loadSourceConfig loads the config file from the given path
func loadSourceConfig(si os.FileInfo, path string) (*configv1.SourceConfig, error) {
	// Check if the file or directory exists
	if si == nil {
		return nil, errors.New("load config file target is not exist")
	}
	var cfg configv1.SourceConfig
	err := loadCustomizeConfig(si, path, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// loadCustomizeConfig loads the user config file from the given path
func loadCustomizeConfig(si os.FileInfo, path string, cfg any) error {
	// Check if the path is a directory
	decode := decodeFile
	if si.IsDir() {
		decode = decodeDir
	}
	err := decode(path, cfg)
	if err != nil {
		return err
	}
	return nil
}

// decodeFile loads the config file from the given path
func decodeFile(path string, cfg any) error {
	// Decode the file into the config struct
	if err := codec.DecodeFromFile(path, cfg); err != nil {
		return errors.Wrapf(err, "failed to parse config file %s", path)
	}
	return nil
}

// decodeDir loads the config file from the given directory
func decodeDir(path string, cfg any) error {
	found := false
	// Walk through the directory and load each file
	err := filepath.WalkDir(path, func(walkpath string, d os.DirEntry, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to get config file %s", walkpath)
		}
		// Check if the path is a directory
		if d.IsDir() {
			return nil
		}

		// Decode the file into the config struct
		if err := decodeFile(walkpath, cfg); err != nil {
			return err
		}
		found = true
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "load config error")
	}
	if !found {
		return errors.New("no config file found in " + path)
	}
	return nil
}

// LoadSourceConfig loads the config file from the given path
func LoadSourceConfig(bootstrap *Bootstrap) (*configv1.SourceConfig, error) {
	// Get the path from the bootstrap
	path := bootstrap.WorkPath()

	// Get the file info from the path
	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "load config stat error")
	}

	// Load the config file
	return loadSourceConfig(stat, path)
}

// LoadPathSourceConfig loads the config file from the given path
func LoadPathSourceConfig(path string) (*configv1.SourceConfig, error) {
	// Get the file info from the path
	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "load config stat error")
	}
	// Load the config file
	return loadSourceConfig(stat, path)
}

// LoadRemoteConfig loads the config file from the given path
func LoadRemoteConfig(bootstrap *Bootstrap, v any, ss ...config.SourceOptionSetting) error {
	sourceConfig, err := LoadSourceConfig(bootstrap)
	if err != nil {
		return err
	}
	config, err := runtime.NewConfig(sourceConfig, config.WithSourceOption())
	if err != nil {
		return err
	}
	if err := config.Load(); err != nil {
		return err
	}
	if err := config.Scan(v); err != nil {
		return err
	}
	return nil
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
