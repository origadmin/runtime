/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package objectstore

import (
	ossv1 "github.com/origadmin/runtime/api/gen/go/config/data/oss/v1" // Changed import
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	internalfactory "github.com/origadmin/runtime/internal/factory"
)

const Module = "storage.objectstore"

// Factory is the interface for creating new ObjectStore components.
type Factory interface {
	New(cfg *ossv1.ObjectStoreConfig) (storageiface.ObjectStore, error) // Changed cfg type
}

// defaultFactory is the default, package-level instance of the object store factory registry.
var defaultFactory = internalfactory.New[Factory]()
