/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements the functions, types, and interfaces for the module.
package registry

type Registry struct {
}

func (r *Registry) Registrar(serviceName string) Registrar {
	panic("implement me")
}

func (r *Registry) Discovery(serviceName string) Discovery {
	panic("implement me")
}
