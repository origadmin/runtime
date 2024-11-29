/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package bootstrap implements the functions, types, and interfaces for the module.
package bootstrap

import (
	"encoding/json"
	"os"

	"github.com/goexts/generic/settings"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// SaveOption represents an option for saving configuration data.
type SaveOption = func(*protojson.MarshalOptions)

// SaveConfig saves the configuration data to the specified file path.
func SaveConfig(path string, data any, opts ...SaveOption) error {
	var bytes []byte
	var err error
	if v, ok := data.(proto.Message); ok {
		opt := settings.Apply(&protojson.MarshalOptions{
			Indent: " ",
		}, opts)
		bytes, err = opt.Marshal(v)
	} else {
		bytes, err = json.Marshal(data)
	}
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, bytes, 0644); err != nil {
		return err
	}
	return nil
}
