/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/goexts/generic/settings"

	"github.com/origadmin/runtime/context"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/toolkits/errors"
	"github.com/origadmin/toolkits/security"
)

// Setting is a function type for setting the Authenticator.
type Setting = func(*Authorizer)

// Authorizer is a struct that implements the Authorizer interface.
type Authorizer struct {
	model        model.Model
	policy       persist.Adapter
	enforcer     *casbin.SyncedEnforcer
	wildcardItem string
}

func (auth Authorizer) Authorized(ctx context.Context, claims security.UserClaims) (bool, error) {
	domain := claims.GetDomain()
	if len(domain) == 0 {
		domain = auth.wildcardItem
	}

	var err error
	var allowed bool
	if allowed, err = auth.enforcer.Enforce(claims.GetSubject(), claims.GetObject(), claims.GetAction(), domain); err != nil {
		return false, err
	} else if allowed {
		return true, nil
	}
	return false, nil
}

func NewAuthorizer(cfg *configv1.Security, ss ...Setting) (security.Authorizer, error) {
	config := cfg.GetAuthz().GetCasbin()
	if config == nil {
		return nil, errors.New("authorizer casbin config is empty")
	}
	auth := settings.Apply(&Authorizer{
		policy: NewAdapter(),
		//projects:                  engine.Projects{},
		wildcardItem: "*",
		//authorizedProjectsMatcher: "g(r.sub, p.sub, p.dom) && (keyMatch(r.dom, p.dom) || p.dom == '*')",
	}, ss)

	var err error
	if auth.model == nil {
		auth.model, err = model.NewModelFromString(DefaultModel())
		if err != nil {
			return nil, err
		}
	}

	auth.enforcer, err = casbin.NewSyncedEnforcer(auth.model, auth.policy)
	if err != nil {
		return nil, err
	}

	return auth, nil
}
