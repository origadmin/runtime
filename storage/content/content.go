package content

import (
	"bytes"
	"errors"
	"io"

	"github.com/origadmin/runtime/interfaces/storage/blob"
	contentiface "github.com/origadmin/runtime/interfaces/storage/content"
	"github.com/origadmin/runtime/interfaces/storage/meta"
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

	metaV2, ok := fileMeta.(interface {
		GetEmbeddedData() []byte
		GetShards() []string
	})
	if !ok {
		return nil, errors.New("unsupported FileMeta type")
	}

	embeddedData := metaV2.GetEmbeddedData()
	if len(embeddedData) > 0 {
		return bytes.NewReader(embeddedData), nil
	}

	shards := metaV2.GetShards()
	if len(shards) > 0 {
		return &chunkReader{
			storage: blobStore,
			hashes:  shards,
		}, nil
	}

	return nil, errors.New("no content available in file meta")
}

// chunkReader is an io.Reader that reads data from multiple chunks in the BlobStore
type chunkReader struct {
	storage blob.Store
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
