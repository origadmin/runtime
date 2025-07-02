/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package meta

import (
	"errors"
	"io"
)

const (
	// DefaultBlockSize is the default size for content-defined blocks.
	// 4MB is a common choice.
	DefaultBlockSize = 4 * 1024 * 1024
)

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
