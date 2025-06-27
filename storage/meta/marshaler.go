/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package meta implements the functions, types, and interfaces for the module.
package meta

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"

	metav1 "github.com/origadmin/runtime/storage/meta/v1"
)

func MarshalFileMeta(meta interface{}) ([]byte, error) {
	switch v := meta.(type) {
	case *metav1.FileMetaV1:
		v.Version = 1
	case metav1.FileMetaV1:
		v.Version = 1
	//case 2:
	//	var meta FileMetaV2
	//	if err := msgpack.Unmarshal(data, &meta); err != nil {
	//		return nil, err
	//	}
	//	return meta, nil
	//case 3:
	//	var meta FileMetaV3
	//	if err := msgpack.Unmarshal(data, &meta); err != nil {
	//		return nil, err
	//	}
	//	return meta, nil
	default:
		return nil, fmt.Errorf("unknown meta type: %T", v)
	}
	return msgpack.Marshal(meta)
}
