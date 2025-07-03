/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package storage implements the functions, types, and interfaces for the module.
package storage

import (
	"github.com/origadmin/toolkits/errors"

	storagev1 "github.com/origadmin/runtime/api/gen/go/storage/v1"
	metav1 "github.com/origadmin/runtime/storage/meta/v1"
	metav2 "github.com/origadmin/runtime/storage/meta/v2"
)

const ErrUnknownFileMetaType = errors.String("storage: unknown file meta type")

func FromFileMeta(meta interface{}) (*storagev1.FileMeta, error) {
	switch v := meta.(type) {
	case *metav1.FileMeta:
		return &storagev1.FileMeta{
			Name:     v.Info.Name,
			Hash:     v.Info.Hash,
			Size:     v.Data.Size,
			MimeType: v.Data.MimeType,
			ModTime:  v.Data.ModTime,
		}, nil
	case metav1.FileMeta:
		return &storagev1.FileMeta{
			Name:     v.Info.Name,
			Hash:     v.Info.Hash,
			Size:     v.Data.Size,
			MimeType: v.Data.MimeType,
			ModTime:  v.Data.ModTime,
		}, nil
	case *metav2.FileMeta:
		return &storagev1.FileMeta{
			Name:     v.Info.Name,
			Hash:     v.Info.Hash,
			Size:     v.Data.Size,
			MimeType: v.Data.MimeType,
			ModTime:  v.Data.ModTime,
		}, nil
	case metav2.FileMeta:
		return &storagev1.FileMeta{
			Name:     v.Info.Name,
			Hash:     v.Info.Hash,
			Size:     v.Data.Size,
			MimeType: v.Data.MimeType,
			ModTime:  v.Data.ModTime,
		}, nil
	default:
		return nil, ErrUnknownFileMetaType
	}
}
