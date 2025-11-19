/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package objectstore

import (
	filev1 "github.com/origadmin/runtime/api/gen/go/config/data/file/v1" // Changed import
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	internalfactory "github.com/origadmin/runtime/internal/factory"
)

const Module = "storage.objectstore"

// Factory is the interface for creating new ObjectStore components.
type Factory interface {
	New(cfg *filev1.FilestoreConfig) (storageiface.ObjectStore, error) // Changed cfg type
}

// defaultFactory is the default, package-level instance of the object store factory registry.
var defaultFactory = internalfactory.New[Factory]()
