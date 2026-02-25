package mapper

// Mapper is a generic interface for mapping a source type (S) to a target type (T).
type Mapper[S any, T any] interface {
	Map(source S) (T, error)
}
