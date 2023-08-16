package memfs

import (
	"io/fs"
	"time"
)

var _ fs.DirEntry = (*inMemDirEntry)(nil)
var _ fs.FileInfo = (*inMemDirEntry)(nil)

// inMemDirEntry is an in-memory implementation of fs.DirEntry.
type inMemDirEntry struct {
	name     string
	isdir    bool
	size     int64
	filemode fs.FileMode
}

func (e *inMemDirEntry) Name() string {
	return e.name
}

func (e *inMemDirEntry) IsDir() bool {
	return e.isdir
}

func (e *inMemDirEntry) Type() fs.FileMode {
	return e.filemode
}

func (e *inMemDirEntry) Info() (fs.FileInfo, error) {
	return e, nil
}

// underlying data source (can return nil)
func (e *inMemDirEntry) Sys() any {
	return nil
}

// length in bytes for regular files; system-dependent for others
func (e *inMemDirEntry) Size() int64 {
	return e.size
}

// file mode bits
func (e *inMemDirEntry) Mode() fs.FileMode {
	return e.filemode
}

// modification time
func (e *inMemDirEntry) ModTime() time.Time {
	return time.Now()
}
