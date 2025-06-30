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

func MarshalFileMeta(meta any) ([]byte, error) {
	switch v := meta.(type) {
	case *metav1.FileMetaV1:
		v.Version = 1
	case metav1.FileMetaV1:
		v.Version = 1
		meta = &v
	case *metav2.FileMetaV2:
		v.Version = 2
	case metav2.FileMetaV2:
		v.Version = 2
		meta = &v
	default:
		return nil, fmt.Errorf("unknown meta type: %T", v)
	}
	return msgpack.Marshal(meta)
}
