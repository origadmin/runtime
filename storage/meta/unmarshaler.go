/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"

	metav1 "github.com/origadmin/runtime/storage/meta/v1"
	metav2 "github.com/origadmin/runtime/storage/meta/v2"
)

func UnmarshalFileMeta(data []byte) (interface{}, error) {
	var version FileMetaVersion
	if err := msgpack.Unmarshal(data, &version); err != nil {
		return nil, err
	}

	switch version.Version {
	case 1:
		var meta metav1.FileMetaV1
		if err := msgpack.Unmarshal(data, &meta); err != nil {
			return nil, err
		}
		return &meta, nil
	case 2:
		var meta metav2.FileMetaV2
		if err := msgpack.Unmarshal(data, &meta); err != nil {
			return nil, err
		}
		return &meta, nil
	default:
		return nil, fmt.Errorf("unsupported file meta version: %d", version.Version)
	}
}
