package content

import (
	"bytes"
	"errors"
	"io"

	blobiface "github.com/origadmin/runtime/interfaces/storage/components/blob"
	contentiface "github.com/origadmin/runtime/interfaces/storage/components/content"
	metaiface "github.com/origadmin/runtime/interfaces/storage/components/meta"
)

// New returns a new content assembler.
func New(blobStore blobiface.Store) contentiface.Assembler {
	return &assembler{blobStore: blobStore}
}

type assembler struct {
	blobStore blobiface.Store
}

// NewReader creates a reader for the content of a file.
func (a *assembler) NewReader(fileMeta metaiface.FileMeta) (io.Reader, error) {
	if fileMeta == nil {
		return nil, errors.New("file meta is nil")
	}

	// No longer need type assertion. We can use the interface methods directly.
	embeddedData := fileMeta.GetEmbeddedData()
	if len(embeddedData) > 0 {
		return bytes.NewReader(embeddedData), nil
	}

	shards := fileMeta.GetShards()
	if len(shards) > 0 {
		return &chunkReader{
			storage: a.blobStore,
			hashes:  shards,
		}, nil
	}

	return nil, errors.New("no content available in file meta")
}

// chunkReader is an io.Reader that reads data from multiple chunks in the BlobStore
type chunkReader struct {
	storage blobiface.Store
	hashes  []string
	current int
	reader  io.Reader
}

func (cr *chunkReader) Read(p []byte) (int, error) {
	for cr.current < len(cr.hashes) {
		if cr.reader == nil {
			data, err := cr.storage.Read(cr.hashes[cr.current])
			if err != nil {
				return 0, err
			}
			cr.reader = bytes.NewReader(data)
		}

		n, err := cr.reader.Read(p)
		if err == io.EOF {
			cr.current++
			cr.reader = nil
			continue
		}
		return n, err
	}
	return 0, io.EOF
}
