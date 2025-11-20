package mapper

import (
	"github.com/origadmin/runtime/interfaces/mapper"
)

// ListMapper maps a slice of source types to a slice of target types.
// It uses the provided Mapper for individual elements.
func ListMapper[S any, T any](sources []S, m mapper.Mapper[S, T]) ([]T, error) {
	if sources == nil {
		return nil, nil
	}
	targets := make([]T, len(sources))
	for i, source := range sources {
		target, err := m.Map(source)
		if err != nil {
			return nil, err
		}
		targets[i] = target
	}
	return targets, nil
}
