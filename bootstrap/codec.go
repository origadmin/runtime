/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package bootstrap implements the functions, types, and interfaces for the module.
package bootstrap

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/origadmin/toolkits/codec"
	"github.com/origadmin/toolkits/errors"
)

type Decoder interface {
	Decode(string, any) error
}

type Encoder interface {
	Encode(any) (string, error)
}

type Codec interface {
	Decoder
	Encoder
}

// decodeFile loads the config file from the given path
func decodeFile(path string, cfg any, ignores []string) error {
	var ignore string
	for _, ignore = range ignores {
		if strings.HasSuffix(path, ignore) {
			return nil
		}
	}
	codec.IsSupportCodec(path)

	// Decode the file into the config struct
	if err := codec.DecodeFromFile(path, cfg); err != nil {
		return errors.Wrapf(err, "failed to parse config file %s", path)
	}
	return nil
}

func decodeDirWithDepth(path string, cfg any, ignores []string, depth int) error {
	found := false
	err := filepath.WalkDir(path, func(walkpath string, d os.DirEntry, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to get config file %s", walkpath)
		}
		if d.IsDir() && depth > 0 {
			return decodeDirWithDepth(walkpath, cfg, ignores, depth-1)
		}

		// Decode the file into the config struct
		if err := decodeFile(walkpath, cfg, ignores); err != nil {
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
		if err := decodeFile(walkpath, cfg, ignores); err != nil {
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

func decodeConfig(path string, cfg any, ignores []string) error {
	// Check if the path is a directory
	info, err := os.Stat(path)
	if err != nil {
		return errors.Wrapf(err, "failed to get config file %s", path)
	}
	if info.IsDir() {
		return decodeDir(path, cfg, ignores)
	}
	return decodeFile(path, cfg, ignores)
}

type bootstrapCodec struct {
	ignores []string
}

func (c bootstrapCodec) Decode(path string, cfg any) error {
	return decodeConfig(path, cfg, c.ignores)
}

func (c bootstrapCodec) Encode(cfg any) (string, error) {
	return "", nil
}

func NewCodec(ignores []string) Codec {
	return &bootstrapCodec{
		ignores: ignores,
	}
}

var _ Codec = (*bootstrapCodec)(nil)
