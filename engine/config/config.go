package config

import (
	"github.com/origadmin/runtime/helpers/configutil"
)

// BlockInfo represents the results of config normalization.
type BlockInfo struct {
	WinnerName string
	Configs    map[string]any
}

// Normalize uses configutil to resolve Active/Default/Configs logic in a type-safe way.
func Normalize[T configutil.Identifiable](active string, def T, configs []T) (*BlockInfo, error) {
	winner, normalized, err := configutil.Normalize(active, def, configs)
	if err != nil {
		return nil, err
	}

	info := &BlockInfo{
		Configs: make(map[string]any),
	}
	if winner != *new(T) {
		info.WinnerName = winner.GetName()
	}

	for _, cfg := range normalized {
		info.Configs[cfg.GetName()] = cfg
	}

	return info, nil
}
