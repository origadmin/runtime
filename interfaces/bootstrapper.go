package interfaces

// Bootstrapper defines the interface for the result of the bootstrap process.
// It provides access to the initialized Container, the Config decoder,
// and a cleanup function to release resources.
type Bootstrapper interface {
	Provider() Container
	Config() Config
	Cleanup() func()
}
