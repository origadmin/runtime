/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package objectstore

import (
	"context"

	ossv1 "github.com/origadmin/runtime/api/gen/go/config/data/oss/v1"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/engine"
)

// DefaultProvider is the default provider for object store components.
var DefaultProvider component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
	cfg, err := engine.AsConfig[ossv1.ObjectStoreConfig](h)
	if err != nil {
		return nil, err
	}
	return New(cfg, opts...)
}
