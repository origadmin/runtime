/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package bootstrap

import (
	"os"
	"path/filepath"

	"github.com/origadmin/runtime"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/toolkits/codec"

	"github.com/origadmin/toolkits/errors"
)

// loadSource loads the config file from the given path
func loadSource(si os.FileInfo, path string) (*configv1.SourceConfig, error) {
	// Check if the file or directory exists
	if si == nil {
		return nil, errors.New("load config file target is not exist")
	}
	var cfg configv1.SourceConfig
	err := loadCustomize(si, path, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// loadCustomize loads the user config file from the given path
func loadCustomize(si os.FileInfo, path string, cfg any) error {
	// Check if the path is a directory
	if si.IsDir() {
		err := loadDir(path, cfg)
		if err != nil {
			return err
		}
		return nil
	}
	// Get the file type from the extension
	typo := codec.TypeFromExt(filepath.Ext(path))
	// Check if the file type is unknown
	if typo == codec.UNKNOWN {
		return errors.New("unknown file type: " + path)
	}
	// Load the config file
	err := loadFile(path, cfg)
	if err != nil {
		return err
	}
	return nil
}

// loadFile loads the config file from the given path
func loadFile(path string, cfg any) error {
	// Decode the file into the config struct
	if err := codec.DecodeFromFile(path, &cfg); err != nil {
		return errors.Wrapf(err, "failed to parse config file %s", path)
	}
	return nil
}

// loadDir loads the config file from the given directory
func loadDir(path string, cfg any) error {
	// Walk through the directory and load each file
	err := filepath.WalkDir(path, func(walkpath string, d os.DirEntry, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to get config file %s", walkpath)
		}
		// Check if the path is a directory
		if d.IsDir() {
			return nil
		}
		// Get the file type from the extension
		typo := codec.TypeFromString(filepath.Ext(walkpath))
		// Check if the file type is unknown
		if typo == codec.UNKNOWN {
			return nil
		}

		// Decode the file into the config struct
		if err := codec.DecodeFromFile(walkpath, &cfg); err != nil {
			return errors.Wrapf(err, "failed to parse config file %s", walkpath)
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "load config error")
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
	return loadSource(stat, path)
}

// LoadSourceConfigFromPath loads the config file from the given path
func LoadSourceConfigFromPath(path string) (*configv1.SourceConfig, error) {
	// Get the file info from the path
	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "load config stat error")
	}
	// Load the config file
	return loadSource(stat, path)
}

// LoadRemoteConfig loads the config file from the given path
func LoadRemoteConfig(bootstrap *Bootstrap, v any) error {
	sourceConfig, err := LoadSourceConfig(bootstrap)
	if err != nil {
		return err
	}
	config, err := runtime.NewConfig(sourceConfig)
	if err != nil {
		return err
	}
	if err := config.Load(); err != nil {
		return err
	}
	return config.Scan(v)
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

	return loadCustomize(stat, path, v)
}
