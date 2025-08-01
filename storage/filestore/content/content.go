package content

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	blobiface "github.com/origadmin/runtime/interfaces/storage/components/blob"
	contentiface "github.com/origadmin/runtime/interfaces/storage/components/content"
	metaiface "github.com/origadmin/runtime/interfaces/storage/components/meta"
	metav2 "github.com/origadmin/runtime/storage/filestore/meta/v2"
)

// assembler implements the Assembler interface.
type assembler struct {
	blobStore blobiface.Store
}

// New creates a new Assembler instance.
func New(blobStore blobiface.Store) *assembler {
	return &assembler{
		blobStore: blobStore,
	}
}

// NewReader creates an io.Reader for the given FileMeta.
// It uses the blobStore to fetch data chunks if necessary.
func (a *assembler) NewReader(fileMeta metaiface.FileMeta) (io.Reader, error) {
	if embeddedData := fileMeta.GetEmbeddedData(); embeddedData != nil {
		return bytes.NewReader(embeddedData), nil
	}

	shards := fileMeta.GetShards()
	if len(shards) == 0 {
		return nil, fmt.Errorf("file has no content")
	}

	return a.readShards(shards)
}

// WriteContent processes the content from the reader, stores it (either embedded or as sharded blobs),
// and returns the content ID and the generated FileMeta object.
func (a *assembler) WriteContent(r io.Reader, size int64) (contentID string, fileMeta metaiface.FileMeta, err error) {
	var isLargeFile bool

	if size > 0 { // Case 1: Size is known
		if size <= metav2.EmbeddedFileSizeThreshold {
			// Known small file: read exact size and embed.
			contentBytes := make([]byte, size)
			if _, err := io.ReadFull(r, contentBytes); err != nil {
				return "", nil, fmt.Errorf("failed to read content for known small file: %w", err)
			}
			contentID = calculateContentHash(contentBytes)
			fileMeta = &metav2.FileMetaV2{
				FileSize:     size,
				ModifyTime:   time.Now().Unix(),
				MimeType:     "application/octet-stream",
				RefCount:     1,
				EmbeddedData: contentBytes,
			}
		} else {
			// Known large file: chunk directly.
			isLargeFile = true
			chunkID, meta, chunkErr := a.chunkData(r)
			if chunkErr != nil {
				// Note: chunkData does not return partial blob hashes, so no cleanup needed here.
				return "", nil, fmt.Errorf("failed to chunk and store data for known large file: %w", chunkErr)
			}
			// Sanity check to ensure the stream provided the expected amount of data.
			if meta.FileSize != size {
				// Cleanup blobs if size doesn't match
				for _, h := range meta.BlobHashes {
					_ = a.blobStore.Delete(h)
				}
				return "", nil, fmt.Errorf("stream size mismatch: provided size %d, but read %d", size, meta.FileSize)
			}
			contentID = chunkID
			fileMeta = meta
		}
	} else { // Case 2: Size is unknown, fall back to peeking.
		buf := new(bytes.Buffer)
		tee := io.TeeReader(r, buf)
		_, err := io.CopyN(io.Discard, tee, metav2.EmbeddedFileSizeThreshold+1)
		if err != nil && err != io.EOF {
			return "", nil, err
		}

		contentBytes := buf.Bytes()
		if err == io.EOF { // It's a small file
			contentID = calculateContentHash(contentBytes)
			fileMeta = &metav2.FileMetaV2{
				FileSize:     int64(len(contentBytes)),
				ModifyTime:   time.Now().Unix(),
				MimeType:     "application/octet-stream",
				RefCount:     1,
				EmbeddedData: contentBytes,
			}
		} else {
			isLargeFile = true
			fullStream := io.MultiReader(bytes.NewReader(contentBytes), r)
			chunkID, meta, chunkErr := a.chunkData(fullStream)
			if chunkErr != nil {
				return "", nil, fmt.Errorf("failed to chunk and store data for unknown size file: %w", chunkErr)
			}
			contentID = chunkID
			fileMeta = meta
		}
	}

	// This case handles a known-size-0 or unknown-size empty stream.
	if contentID == "" && (fileMeta == nil || fileMeta.Size() == 0) {
		contentID = calculateContentHash([]byte{})
		if fileMeta == nil {
			fileMeta = &metav2.FileMetaV2{
				FileSize:   0,
				ModifyTime: time.Now().Unix(),
				MimeType:   "application/octet-stream",
				RefCount:   1,
			}
		}
	}

	// If persisting metadata fails, we must clean up any blobs we just created.
	// This cleanup logic is now handled by the caller (meta/service.go)
	// as it's responsible for the overall transaction.
	return contentID, fileMeta, nil
}

// chunkData reads from the reader and splits the content into blocks of a fixed size.
// It processes a large file stream by chunking it, storing the chunks as blobs,
// and returns a fully populated FileMetaV2 object along with a hash of the entire content stream.
func (a *assembler) chunkData(r io.Reader) (string, *metav2.FileMetaV2, error) {
	var hashes []string
	var totalSize int64
	// Use a default chunk size if not configured, or pass it from the service
	chunkSize := 4 * 1024 * 1024 // DefaultChunkSize, should come from config

	buf := make([]byte, chunkSize)

	// Create a hasher that will calculate the hash of the entire stream.
	streamHasher := sha256.New()
	// TeeReader copies data from r to streamHasher as it's being read.
	teeReader := io.TeeReader(r, streamHasher)

	for {
		// Read from the teeReader to ensure the hasher gets the data.
		n, err := teeReader.Read(buf)
		if err != nil && err != io.EOF {
			return "", nil, err
		}
		if n == 0 {
			break
		}

		data := buf[:n]
		// Write the chunk to the blob storage. The blob store will return the hash of this chunk.
		blobHash, storeErr := a.blobStore.Write(data)
		if storeErr != nil {
			return "", nil, storeErr
		}
		hashes = append(hashes, blobHash)
		totalSize += int64(n)

		if err == io.EOF {
			break
		}
	}

	// Finalize the hash of the entire stream.
	overallHash := hex.EncodeToString(streamHasher.Sum(nil))

	meta := &metav2.FileMetaV2{
		FileSize:   totalSize,
		ModifyTime: time.Now().Unix(),
		MimeType:   "application/octet-stream",
		RefCount:   1,
		BlobHashes: hashes,
		BlobSize:   int32(chunkSize), // Use the actual chunk size
	}

	return overallHash, meta, nil
}

// calculateContentHash is a helper to generate a hash from byte data.
func calculateContentHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// ... (NewReader 保持不变) ...

// WriteContent 成为处理内容的唯一权威实现
func (a *assembler) WriteContent(r io.Reader, size int64) (contentID string, fileMeta metaiface.FileMeta, err error) {
	// ... (这里的逻辑保持不变，它已经是正确的) ...
	// 唯一的改动是在调用 a.chunkData 时，chunkData 内部会使用 a.chunkSize
	// ...
}

// chunkData 使用结构体中的 chunkSize
func (a *assembler) chunkData(r io.Reader) (string, *metav2.FileMetaV2, error) {
	var hashes []string
	var totalSize int64
	// 使用配置的 chunk size，而不是硬编码
	buf := make([]byte, a.chunkSize)

	// ... (其余逻辑不变) ...

	meta := &metav2.FileMetaV2{
		FileSize:   totalSize,
		ModifyTime: time.Now().Unix(),
		MimeType:   "application/octet-stream",
		RefCount:   1,
		BlobHashes: hashes,
		BlobSize:   int32(a.chunkSize), // <-- 使用实际的 chunk size
	}

	return overallHash, meta, nil
}

func (a *assembler) readShards(shards []string) (io.Reader, error) {
	return newChunkReader(a.blobStore, shards), nil
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

func newChunkReader(storage blobiface.Store, hashes []string) io.Reader {
	return &chunkReader{
		storage: storage,
		hashes:  hashes,
	}
}
