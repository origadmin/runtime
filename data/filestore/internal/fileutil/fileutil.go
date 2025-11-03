package fileutil

import (
	"os"
)

// AtomicWrite writes data to a file atomically.
// It writes to a temporary file and then renames it to the final path.
func AtomicWrite(path string, content []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, content, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
