package interfaces

// Bootstrapper defines the interface for the result of the bootstrap process.
// It provides access to the initialized Container, the Config decoder,
// and a cleanup function to release resources.
type Bootstrapper interface {
	// AppInfo returns the application information.
	AppInfo() *AppInfo
	// Container returns the initialized component provider.
	Container() Container
	// Config returns the configuration decoder.
	Config() StructuredConfig
	// Cleanup returns the cleanup function.
	Cleanup() func()
}
