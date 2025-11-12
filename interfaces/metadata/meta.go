// Package metadata implements the functions, types, and interfaces for the module.
package metadata

import (
	securityv1 "github.com/origadmin/runtime/api/gen/go/config/security/v1"
)

type Meta interface {
	Append(key string, values ...string)
	Values(key string) []string
	Get(key string) string
	Set(key string, value string)
	Clone() Meta
	ToProto() map[string]*securityv1.MetaValue
	GetAll() map[string][]string
}
