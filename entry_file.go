package memfs

import "io/fs"

var _ Entry = File(nil)

// File is a type of Entry that represents a file.
type File []byte

func (f File) ToEntry(name string) fs.DirEntry {
	return &inMemDirEntry{
		name:     name,
		isdir:    false,
		size:     int64(len(f)),
		filemode: 0,
	}
}
