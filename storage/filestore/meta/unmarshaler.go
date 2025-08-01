/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/origadmin/runtime/interfaces/storage/components/meta"
	metav1 "github.com/origadmin/runtime/storage/filestore/meta/v1"
	metav2 "github.com/origadmin/runtime/storage/filestore/meta/v2"
)

func Unmarshal(data []byte) (meta.FileMeta, error) {
	var versionOnly meta.FileMetaVersion
	if err := msgpack.Unmarshal(data, &versionOnly); err != nil {
		return nil, err
	}

	switch versionOnly.Version {
	case metav1.Version:
		var fileMetaData meta.StoreMeta[metav1.FileMetaV1]
		if err := msgpack.Unmarshal(data, &fileMetaData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal FileMetaV1: %w", err)
		}
		return fileMetaData.Data, nil
	case metav2.Version:
		var fileMetaData meta.StoreMeta[metav2.FileMetaV2]
		if err := msgpack.Unmarshal(data, &fileMetaData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal FileMetaV2: %w", err)
		}
		return fileMetaData.Data, nil
	default:
		return nil, fmt.Errorf("unsupported metaFile meta version: %d", versionOnly.Version)
	}
}

func UnmarshalFileMeta(data []byte) (meta.FileMeta, error) {
	return Unmarshal(data) // Simply call the main Unmarshal function
}
