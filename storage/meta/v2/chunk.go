// Package metav2 implements the functions, types, and interfaces for the module.
package metav2

import (
	"bytes"
	"io"

	"github.com/origadmin/runtime/interfaces/storage/meta"
)

type chunkReader struct {
	storage meta.BlobStorage
	hashes  []string
	current int
	reader  io.Reader
}

func (cr *chunkReader) Read(p []byte) (int, error) {
	for cr.current < len(cr.hashes) {
		if cr.reader == nil {
			data, err := cr.storage.Retrieve(cr.hashes[cr.current])
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
