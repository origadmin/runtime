/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"

	metaiface "github.com/origadmin/runtime/interfaces/storage/components/meta"
	metav1 "github.com/origadmin/runtime/storage/filestore/meta/v1"
	metav2 "github.com/origadmin/runtime/storage/filestore/meta/v2"
)

func MarshalFileMeta(meta metaiface.FileMeta) ([]byte, error) {
	var fileMetaData interface{} // Use interface{} to hold the concrete StoreMeta[T] type

	switch v := meta.(type) {
	case *metav1.FileMetaV1:
		fileMetaData = &metaiface.StoreMeta[*metav1.FileMetaV1]{
			Version: metav1.Version,
			Data:    v,
		}
	case metav1.FileMetaV1:
		fileMetaData = &metaiface.StoreMeta[metav1.FileMetaV1]{
			Version: metav1.Version,
			Data:    v,
		}
	case *metav2.FileMetaV2:
		fileMetaData = &metaiface.StoreMeta[*metav2.FileMetaV2]{
			Version: metav2.Version,
			Data:    v,
		}
	case metav2.FileMetaV2:
		fileMetaData = &metaiface.StoreMeta[metav2.FileMetaV2]{
			Version: metav2.Version,
			Data:    v,
		}
	default:
		return nil, fmt.Errorf("unknown meta type: %T", meta)
	}
	return msgpack.Marshal(fileMetaData)
}
