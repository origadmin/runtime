/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware/selector"

	"github.com/origadmin/runtime/interfaces/options"
	"github.com/origadmin/runtime/optionutil"
)

type Options struct {
	MatchFunc selector.MatchFunc
}

type Option = options.Option

func WithMatchFunc(matchFunc selector.MatchFunc) options.Option {
	return optionutil.Update(func(o *Options) {
		o.MatchFunc = matchFunc
	})
}

func FromOptions(opts ...options.Option) (options.Context, *Options) {
	return optionutil.ApplyNew[Options](opts...)
}
