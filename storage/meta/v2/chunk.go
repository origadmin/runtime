// Package metav2 implements the functions, types, and interfaces for the module.
package metav2

import (
	"bytes"
	"errors"
	"io"

	blob_interface "github.com/origadmin/runtime/interfaces/storage/blob"
)

const (
	// DefaultBlockSize is the default size for content-defined blocks.
	// 4MB is a common choice.
	DefaultBlockSize = 4 * 1024 * 1024
)

type chunkReader struct {
	storage blob_interface.BlobStore
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

// chunkData reads from the reader and splits the content into blocks of a fixed size.
// For each block, it calls the provided store function.
func chunkData(r io.Reader, store func([]byte) (string, error)) ([]string, int64, error) {
	var hashes []string
	var totalSize int64
	buf := make([]byte, DefaultBlockSize)

	for {
		n, err := io.ReadFull(r, buf)
		if err != nil && err != io.EOF && !errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, 0, err
		}
		if n == 0 {
			break
		}

		data := buf[:n]
		hash, storeErr := store(data)
		if storeErr != nil {
			return nil, 0, storeErr
		}
		hashes = append(hashes, hash)
		totalSize += int64(n)

		if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
			break
		}
	}

	return hashes, totalSize, nil
}
