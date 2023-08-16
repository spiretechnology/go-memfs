package memfs

import "io/fs"

// Dir is a type of Entry that represents a directory with no contents.
type Dir struct{}

func (f Dir) ToEntry(name string) fs.DirEntry {
	return &inMemDirEntry{
		name:     name,
		isdir:    true,
		size:     0,
		filemode: fs.ModeDir,
	}
}
