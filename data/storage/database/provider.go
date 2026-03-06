/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package database

import (
	"context"

	databasev1 "github.com/origadmin/runtime/api/gen/go/config/data/database/v1"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/helpers/comp"
)

// DefaultProvider is the default provider for database components.
var DefaultProvider component.Provider = func(ctx context.Context, h component.Handle) (any, error) {
	cfg, err := comp.AsConfig[databasev1.DatabaseConfig](h)
	if err != nil {
		return nil, err
	}
	return New(cfg)
}

