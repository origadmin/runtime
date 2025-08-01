/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package storage implements the functions, types, and interfaces for the module.
package storage

import (
	"github.com/origadmin/toolkits/errors"

	storagev1 "github.com/origadmin/runtime/api/gen/go/storage/v1"
	metaiface "github.com/origadmin/runtime/interfaces/storage/components/meta"
)

const ErrUnknownFileMetaType = errors.String("storage: unknown file meta type")

func FromFileMeta(meta metaiface.FileMeta) (*storagev1.FileMeta, error) {
	return &storagev1.FileMeta{
		Name:     "", // Name is not directly available from FileMeta interface
		Hash:     "", // Hash is not directly available from FileMeta interface
		Size:     meta.Size(),
		MimeType: "", // MimeType is not directly available from FileMeta interface
		ModTime:  meta.ModTime().Unix(),
	}, nil
}
