/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package storage implements the functions, types, and interfaces for the module.
package storage

import (
	storagev1 "github.com/origadmin/runtime/api/gen/go/storage/v1"
	metav1 "github.com/origadmin/runtime/storage/meta/v1"
	metav2 "github.com/origadmin/runtime/storage/meta/v2"
	"github.com/origadmin/toolkits/errors"
)

const ErrUnknownFileMetaType = errors.String("storage: unknown file meta type")

func FromFileMeta(meta interface{}) (*storagev1.FileMeta, error) {
	switch v := meta.(type) {
	case *metav1.FileMetaV1:
		return &storagev1.FileMeta{
			Name:     v.Name,
			Hash:     v.Hash,
			Size:     v.Size,
			MimeType: v.MimeType,
			ModTime:  v.ModTime,
		}, nil
	case metav1.FileMetaV1:
		return &storagev1.FileMeta{
			Name:     v.Name,
			Hash:     v.Hash,
			Size:     v.Size,
			MimeType: v.MimeType,
			ModTime:  v.ModTime,
		}, nil
	case *metav2.FileMetaV2:
		return &storagev1.FileMeta{
			Name:     v.Name,
			Hash:     v.Hash,
			Size:     v.Size,
			MimeType: v.MimeType,
			ModTime:  v.ModTime,
		}, nil
	case metav2.FileMetaV2:
		return &storagev1.FileMeta{
			Name:     v.Name,
			Hash:     v.Hash,
			Size:     v.Size,
			MimeType: v.MimeType,
			ModTime:  v.ModTime,
		}, nil
	default:
		return nil, ErrUnknownFileMetaType
	}
}
