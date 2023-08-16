package memfs

import (
	"bytes"
	"io"
	"io/fs"
	"path"
	"strings"

	"golang.org/x/exp/slices"
)

var _ fs.FS = FS(nil)

type FS map[string]Entry

func (f FS) Open(name string) (fs.File, error) {
	// If the map is nil, return an error
	if f == nil {
		return nil, fs.ErrNotExist
	}

	// Normalize the path
	trimmed := path.Clean(name)
	trimmed = strings.Trim(trimmed, "/")
	if trimmed == "." {
		trimmed = ""
	}

	// Find the entry with the pathname
	entry, ok := f[trimmed]
	if !ok {
		return nil, fs.ErrNotExist
	}

	// If it's not a file (ie. a directory) then return an error
	fileEntry, ok := entry.(File)
	if !ok {
		return nil, fs.ErrInvalid
	}

	// Wrap the file in a struct that implements the fs.File interface
	basename := path.Base(trimmed)
	return &inMemFile{
		reader: io.NopCloser(bytes.NewReader(fileEntry)),
		entry:  fileEntry.ToEntry(basename),
	}, nil
}

func (f FS) ReadDir(name string) ([]fs.DirEntry, error) {
	// If the map is nil, return an error
	if f == nil {
		return nil, fs.ErrNotExist
	}

	// Normalize the path
	trimmed := path.Clean(name)
	trimmed = strings.Trim(trimmed, "/")
	if trimmed == "." {
		trimmed = ""
	}

	// Check if there is an empty directory entry at the path
	entryAtPath, hasEntry := f[trimmed]
	if hasEntry {
		// If the entry at that path is actually a file, return an error
		if _, ok := entryAtPath.(Dir); !ok {
			return nil, fs.ErrInvalid
		}
	}

	// Look for all the entries that are children of the dir path
	var entries []fs.DirEntry
	intermediateDirs := make(map[string]struct{})
	for path, entry := range f {
		// If the entry isn't in the directory, skip
		if trimmed != "" && !strings.HasPrefix(path, trimmed+"/") {
			continue
		}

		// If the entry isn't a direct child of the directory, add an entry for the
		// intermediate directory
		suffix := strings.TrimPrefix(path, trimmed+"/")
		if strings.Contains(suffix, "/") {
			dirname := strings.Split(suffix, "/")[0]
			if _, ok := intermediateDirs[dirname]; !ok {
				intermediateDirs[dirname] = struct{}{}
				entries = append(entries, Dir{}.ToEntry(dirname))
			}
			continue
		}

		// Add the entry to the list of entries
		entries = append(entries, entry.ToEntry(suffix))
	}

	// Sort the entries by name
	slices.SortFunc(entries, func(a, b fs.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})

	return entries, nil
}
