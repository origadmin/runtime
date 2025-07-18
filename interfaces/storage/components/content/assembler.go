package content

import (
	"io"

	metaiface "github.com/origadmin/runtime/interfaces/storage/components/meta"
)

// Assembler is responsible for assembling file content from metadata and blob storage.
type Assembler interface {
	// NewReader creates an io.Reader for the given FileMeta.
	// It uses the blobStore to fetch data chunks if necessary.
	NewReader(fileMeta metaiface.FileMeta) (io.Reader, error)

	// WriteContent processes the content from the reader, stores it (either embedded or as sharded blobs),
	// and returns the content ID and the generated FileMeta object.
	WriteContent(r io.Reader, size int64) (contentID string, fileMeta metaiface.FileMeta, err error)
}
