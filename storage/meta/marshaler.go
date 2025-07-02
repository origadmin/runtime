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
	case *metav1.FileMeta:
		v.Version = 1
	case metav1.FileMeta:
		v.Version = 1
		meta = &v
	case *metav2.FileMeta:
		v.Version = 2
	case metav2.FileMeta:
		v.Version = 2
		meta = &v
	case *DirectoryIndex:
		v.Version = 1 // Or the appropriate version
	default:
		return nil, fmt.Errorf("unknown meta type: %T", v)
	}
	return msgpack.Marshal(meta)
}
