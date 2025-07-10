/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/origadmin/runtime/interfaces/storage/meta"
	metav1 "github.com/origadmin/runtime/storage/meta/v1"
	metav2 "github.com/origadmin/runtime/storage/meta/v2"
)

func Unmarshal(data []byte) (*metav2.FileMetaV2, error) {
	var version meta.FileMetaVersion
	if err := msgpack.Unmarshal(data, &version); err != nil {
		return nil, err
	}

	switch version.Version {
	case 1:
		var metadata metav1.FileMetaV1
		if err := msgpack.Unmarshal(data, &metadata); err != nil {
			return nil, err
		}
		return &metav2.FileMetaV2{
			Version:      metav2.Version,
			FileSize:     metadata.FileSize,
			ModifyTime:   metadata.ModifyTime,
			MimeType:     metadata.MimeType,
			RefCount:     metadata.RefCount,
			BlobSize:     0,
			BlobHashes:   nil,
			EmbeddedData: nil,
		}, nil
	case 2:
		var metadata metav2.FileMetaV2
		if err := msgpack.Unmarshal(data, &metadata); err != nil {
			return nil, err
		}
		return &metadata, nil
	default:
		return nil, fmt.Errorf("unsupported metaFile meta version: %d", version.Version)
	}
}

func UnmarshalFileMeta(data []byte) (meta.FileMeta, error) {
	metad, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return metad, nil
}
