/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package runtime

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
)

// ProvideLogger is a Wire provider function that extracts the logger
// from the App interface. It is intended to be used by the application's
// own Wire injector.
func ProvideLogger(rt *App) log.Logger {
	return rt.Logger()
}

// ProvideDefaultRegistrar is a Wire provider function that extracts the default registrar
// from the App interface. It is intended to be used by the application's
// own Wire injector.
func ProvideDefaultRegistrar(rt *App) registry.Registrar {
	return rt.DefaultRegistrar()
}
