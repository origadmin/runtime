package factory

// Registry defines the interface for managing factory functions.
// It allows for registering and retrieving factory functions by name.
type Registry[F any] interface {
	Get(name string) (F, bool)
	Register(name string, factory F)
	RegisteredFactories() map[string]F
	Reset()
}
