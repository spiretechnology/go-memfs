package memfs

import "io/fs"

// Entry is the interface that describes a file or directory.
type Entry interface {
	ToEntry(name string) fs.DirEntry
}
