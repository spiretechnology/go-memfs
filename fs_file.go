package memfs

import (
	"io"
	"io/fs"
)

var _ fs.File = (*inMemFile)(nil)

// inMemFile is a struct that implements the fs.File interface.
type inMemFile struct {
	reader io.ReadCloser
	entry  fs.DirEntry
}

func (e *inMemFile) Name() string {
	return e.entry.Name()
}

func (e *inMemFile) Stat() (fs.FileInfo, error) {
	return e.entry.Info()
}

func (e *inMemFile) Read(p []byte) (int, error) {
	return e.reader.Read(p)
}

func (e *inMemFile) Close() error {
	return e.reader.Close()
}
