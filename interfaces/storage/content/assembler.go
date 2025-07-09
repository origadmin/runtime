package content

import (
	"io"

	"github.com/origadmin/runtime/interfaces/storage/blob"
	"github.com/origadmin/runtime/interfaces/storage/meta"
)

// Assembler is responsible for assembling file content from metadata and blob storage.
type Assembler interface {
	// NewReader creates an io.Reader for the given FileMeta.
	// It uses the blobStore to fetch data chunks if necessary.
	NewReader(fileMeta meta.FileMeta, blobStore blob.Store) (io.Reader, error)
}
