/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package objectstore

import (
	ossv1 "github.com/origadmin/runtime/api/gen/go/config/data/oss/v1" // Changed import
	"github.com/origadmin/runtime/contracts/options"
	storageiface "github.com/origadmin/runtime/contracts/storage"
	internalfactory "github.com/origadmin/runtime/helpers/builderutil"
)

const Module = "storage.objectstore"

// Factory is the interface for creating new ObjectStore components.
type Factory interface {
	New(cfg *ossv1.ObjectStoreConfig, opts ...options.Option) (storageiface.ObjectStore, error) // Changed cfg type
}

// defaultFactory is the default, package-level instance of the object store factory registry.
var defaultFactory = internalfactory.New[Factory]()
