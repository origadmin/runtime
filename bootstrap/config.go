/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package bootstrap implements the functions, types, and interfaces for the module.
package bootstrap

import (
	"os"
	"path/filepath"
	"strings"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/errors"
)

// decodeFile loads the config file from the given path
func decodeFile(path string, cfg any) error {
	c := codec.TypeFromPath(path)
	if !c.IsSupported() {
		return errors.New("unsupported config file type: " + path)
	}

	// Decode the file into the config struct
	if err := codec.DecodeFromFile(path, cfg); err != nil {
		return errors.Wrapf(err, "failed to parse config file %s", path)
	}
	return nil
}

func decodeDirWithDepth(path string, cfg any, ignores []string, depth int) error {
	found := false
	var ignore string
	err := filepath.WalkDir(path, func(walkpath string, d os.DirEntry, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to get config file %s", walkpath)
		}
		if d.IsDir() && depth > 0 {
			return decodeDirWithDepth(walkpath, cfg, ignores, depth-1)
		}

		for _, ignore = range ignores {
			if strings.HasSuffix(walkpath, ignore) {
				return nil
			}
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

// decodeDir loads the config file from the given directory
func decodeDir(path string, cfg any, ignores []string) error {
	found := false
	err := filepath.WalkDir(path, func(walkpath string, d os.DirEntry, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to get config file %s", walkpath)
		}
		if d.IsDir() {
			return decodeDirWithDepth(walkpath, cfg, ignores, 3)
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

func decodeConfig(path string, cfg any) error {
	// Check if the path is a directory
	info, err := os.Stat(path)
	if err != nil {
		return errors.Wrapf(err, "failed to get config file %s", path)
	}
	if info.IsDir() {
		return errors.New("config path is a directory")
	}
	return decodeFile(path, cfg)
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

// LoadSourceConfigFromBootstrap loads the config file from the given path
func LoadSourceConfigFromBootstrap(bootstrap *Bootstrap) (*configv1.SourceConfig, error) {
	// Get the path from the bootstrap
	return loadSourceConfig(bootstrap.ConfigFilePath())
}

// LoadSourceConfig loads the config file from the given path
func LoadSourceConfig(path string) (*configv1.SourceConfig, error) {
	// Get the file info from the path
	info, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed locate to config path %s", path)
	}
	if info.IsDir() {
		return nil, errors.New("config path is a directory")
	}
	return loadSourceConfig(path)
}
