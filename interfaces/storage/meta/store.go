package meta

// Store defines the interface for managing file content metadata.
type Store interface {
	Create(fileMeta FileMeta) (string, error)
	Get(id string) (FileMeta, error)
	Update(id string, fileMeta FileMeta) error
	Delete(id string) error
	// Add other methods as per storage.md or actual implementation needs
}
