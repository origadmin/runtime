package storage

// Provider is the top-level storage provider.
type Provider interface {
	// FileStore returns the configured file storage service by name.
	FileStore(name string) (FileStore, error)
	// DefaultFileStore returns the default file storage service.
	DefaultFileStore() (FileStore, error)

	// Cache returns the configured cache service by name.
	Cache(name string) (Cache, error)
	// DefaultCache returns the default cache service.
	DefaultCache() (Cache, error)

	// Database returns the configured database service by name.
	Database(name string) (Database, error)
	// DefaultDatabase returns the default database service.
	DefaultDatabase() (Database, error)
}
