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

func Unmarshal(data []byte) (any, error) {
	var version meta.FileMetaVersion
	if err := msgpack.Unmarshal(data, &version); err != nil {
		return nil, err
	}

	switch version.Version {
	case 1:
		var metadata metav1.FileMeta
		if err := msgpack.Unmarshal(data, &metadata); err != nil {
			// try to unmarshal as DirectoryIndex
			var dir DirectoryIndex
			if err2 := msgpack.Unmarshal(data, &dir); err2 == nil {
				return &dir, nil
			}
			return nil, err
		}
		return &metav2.FileMeta{
			Data: &metav2.FileMetaV2{
				FileSize:   metadata.Data.Size,
				ModifyTime: metadata.Data.ModTime,
				MimeType:   metadata.Data.MimeType,
			},
		}, nil
	case 2:
		var metadata metav2.FileMeta
		if err := msgpack.Unmarshal(data, &metadata); err != nil {
			return nil, err
		}
		return &metadata, nil
	default:
		return nil, fmt.Errorf("unsupported file meta version: %d", version.Version)
	}
}

func UnmarshalFileMeta(data []byte) (*metav2.FileMeta, error) {
	meta, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}
	if fileMeta, ok := meta.(*metav2.FileMeta); ok {
		return fileMeta, nil
	}
	return nil, fmt.Errorf("unexpected meta type: %T", meta)
}
