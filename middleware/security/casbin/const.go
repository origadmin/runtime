/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package casbin

import (
	"github.com/origadmin/runtime/middleware/security/internal/model"
	"github.com/origadmin/runtime/middleware/security/internal/policy"
)

func DefaultModel() string {
	return model.DefaultRestfullWithRoleModel
}

func DefaultPolicy() []byte {
	return policy.MustPolicy("keymatch_with_rbac_in_domain.csv")
}
