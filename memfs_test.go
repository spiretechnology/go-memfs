package memfs_test

import (
	"io/fs"
	"testing"

	"github.com/spiretechnology/go-memfs"
	"github.com/stretchr/testify/require"
)

func TestMemFS(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		fsys := memfs.FS{
			"hello/world.txt":          memfs.File("hello"),
			"hello/golang.txt":         memfs.File("foo"),
			"hello/world/testfile.txt": memfs.File("bar"),
			"emptydir":                 memfs.Dir{},
		}

		b, err := fs.ReadFile(fsys, "hello/world.txt")
		require.NoError(t, err, "reading file")
		require.Equal(t, []byte("hello"), b, "reading file")

		b, err = fs.ReadFile(fsys, "hello/golang.txt")
		require.NoError(t, err, "reading file")
		require.Equal(t, []byte("foo"), b, "reading file")

		b, err = fs.ReadFile(fsys, "hello/world/testfile.txt")
		require.NoError(t, err, "reading file")
		require.Equal(t, []byte("bar"), b, "reading file")

		entries, err := fsys.ReadDir("hello")
		require.NoError(t, err, "reading dir")
		require.Len(t, entries, 3, "reading dir")

		info, err := entries[0].Info()
		require.NoError(t, err, "reading entry info")
		require.Equal(t, "golang.txt", entries[0].Name(), "entry name")
		require.Equal(t, "golang.txt", info.Name(), "entry info name")
		require.Equal(t, false, entries[0].IsDir(), "entry is not dir")
		require.Equal(t, false, info.IsDir(), "entry info is not dir")
		require.Equal(t, int64(3), info.Size(), "entry info size")

		info, err = entries[1].Info()
		require.NoError(t, err, "reading entry info")
		require.Equal(t, "world", entries[1].Name(), "entry name")
		require.Equal(t, "world", info.Name(), "entry info name")
		require.Equal(t, true, entries[1].IsDir(), "entry is dir")
		require.Equal(t, true, info.IsDir(), "entry info is dir")

		info, err = entries[2].Info()
		require.NoError(t, err, "reading entry info")
		require.Equal(t, "world.txt", entries[2].Name(), "entry name")
		require.Equal(t, "world.txt", info.Name(), "entry info name")
		require.Equal(t, false, entries[2].IsDir(), "entry is not dir")
		require.Equal(t, false, info.IsDir(), "entry info is not dir")
		require.Equal(t, int64(5), info.Size(), "entry info size")

		// Read the root directory
		for _, path := range []string{".", "/", "", "./"} {
			entries, err = fsys.ReadDir(path)
			require.NoError(t, err, "reading dir with %s", path)
			require.Len(t, entries, 2, "reading dir with %s", path)
			require.Equal(t, "emptydir", entries[0].Name(), "entry name")
			require.Equal(t, "hello", entries[1].Name(), "entry name")
		}
	})
	t.Run("attempting Open() on a directory", func(t *testing.T) {
		fsys := memfs.FS{
			"hello/world/testfile.txt": memfs.File("bar"),
		}
		_, err := fsys.Open("hello/world")
		require.Error(t, err, "opening directory")
	})
	t.Run("attempting ReadDir() on a file", func(t *testing.T) {
		fsys := memfs.FS{
			"hello/world/testfile.txt": memfs.File("bar"),
		}
		_, err := fsys.ReadDir("hello/world/testfile.txt")
		require.Error(t, err, "readdir a file")
	})
	t.Run("attempting ReadDir() on a directory that doesn't exist", func(t *testing.T) {
		fsys := memfs.FS{
			"hello/world/testfile.txt": memfs.File("bar"),
		}
		_, err := fsys.ReadDir("hello/doesntexist")
		require.Error(t, err, "readdir a non-existent directory")
	})
	t.Run("attempting ReadDir() on empty root directory", func(t *testing.T) {
		fsys := memfs.FS{}
		_, err := fsys.ReadDir(".")
		require.NoError(t, err, "readdir empty root directory")
	})
	t.Run("adding files", func(t *testing.T) {
		fsys := memfs.FS{
			"hello/foo.txt": memfs.File("hello"),
		}

		entries, err := fsys.ReadDir("hello")
		require.NoError(t, err, "reading dir")
		require.Len(t, entries, 1, "reading dir")
		require.Equal(t, "foo.txt", entries[0].Name(), "entry name")

		fsys["hello/bar.txt"] = memfs.File("world")

		entries, err = fsys.ReadDir("hello")
		require.NoError(t, err, "reading dir")
		require.Len(t, entries, 2, "reading dir")
		require.Equal(t, "bar.txt", entries[0].Name(), "entry name")
		require.Equal(t, "foo.txt", entries[1].Name(), "entry name")
	})
	t.Run("removing files", func(t *testing.T) {
		fsys := memfs.FS{
			"hello/foo.txt": memfs.File("hello"),
			"hello/bar.txt": memfs.File("world"),
		}

		entries, err := fsys.ReadDir("hello")
		require.NoError(t, err, "reading dir")
		require.Len(t, entries, 2, "reading dir")
		require.Equal(t, "bar.txt", entries[0].Name(), "entry name")
		require.Equal(t, "foo.txt", entries[1].Name(), "entry name")

		delete(fsys, "hello/bar.txt")

		entries, err = fsys.ReadDir("hello")
		require.NoError(t, err, "reading dir")
		require.Len(t, entries, 1, "reading dir")
		require.Equal(t, "foo.txt", entries[0].Name(), "entry name")
	})
	t.Run("dirty path names", func(t *testing.T) {
		fsys := memfs.FS{
			"hello/world.txt":          memfs.File("hello"),
			"hello/golang.txt":         memfs.File("foo"),
			"hello/world/testfile.txt": memfs.File("bar"),
			"emptydir":                 memfs.Dir{},
		}

		b, err := fs.ReadFile(fsys, "./hello/../hello/world.txt")
		require.NoError(t, err, "reading file")
		require.Equal(t, []byte("hello"), b, "reading file")

		b, err = fs.ReadFile(fsys, "/hello//golang.txt")
		require.NoError(t, err, "reading file")
		require.Equal(t, []byte("foo"), b, "reading file")

		b, err = fs.ReadFile(fsys, "hello/world/../../hello/../hello//world/testfile.txt")
		require.NoError(t, err, "reading file")
		require.Equal(t, []byte("bar"), b, "reading file")

		entries, err := fsys.ReadDir("./hello/world/..")
		require.NoError(t, err, "reading dir")
		require.Len(t, entries, 3, "reading dir")
	})
	t.Run("reading nested directories", func(t *testing.T) {
		fsys := memfs.FS{
			"hello/foo/a":       memfs.File(""),
			"hello/foo/bar/a":   memfs.File(""),
			"hello/foo/bar/baz": memfs.Dir{},
			"hello/bar/a":       memfs.File(""),
		}

		entries, err := fsys.ReadDir(".")
		require.NoError(t, err, "reading dir")
		require.ElementsMatch(t, []string{"hello"}, entryNames(entries), "reading dir")

		entries, err = fsys.ReadDir("hello")
		require.NoError(t, err, "reading dir")
		require.ElementsMatch(t, []string{"bar", "foo"}, entryNames(entries), "reading dir")

		entries, err = fsys.ReadDir("hello/foo")
		require.NoError(t, err, "reading dir")
		require.ElementsMatch(t, []string{"a", "bar"}, entryNames(entries), "reading dir")

		entries, err = fsys.ReadDir("hello/foo/bar")
		require.NoError(t, err, "reading dir")
		require.ElementsMatch(t, []string{"a", "baz"}, entryNames(entries), "reading dir")

		entries, err = fsys.ReadDir("hello/foo/bar/baz")
		require.NoError(t, err, "reading dir")
		require.ElementsMatch(t, []string{}, entryNames(entries), "reading dir")

		entries, err = fsys.ReadDir("hello/bar")
		require.NoError(t, err, "reading dir")
		require.ElementsMatch(t, []string{"a"}, entryNames(entries), "reading dir")
	})
	t.Run("stat files and directories", func(t *testing.T) {
		fsys := memfs.FS{
			"hello/foo/a":       memfs.File("helloworld"),
			"hello/foo/bar/baz": memfs.Dir{},
		}

		stat, err := fsys.Stat("hello/foo/a")
		require.NoError(t, err, "stat file")
		require.Equal(t, "a", stat.Name(), "stat file name")
		require.Equal(t, false, stat.IsDir(), "stat file is not dir")
		require.Equal(t, int64(10), stat.Size(), "stat file size")

		stat, err = fsys.Stat("hello/foo/bar/baz")
		require.NoError(t, err, "stat dir")
		require.Equal(t, "baz", stat.Name(), "stat dir name")
		require.Equal(t, true, stat.IsDir(), "stat dir is dir")

		stat, err = fsys.Stat("hello/foo/doesntexist")
		require.Error(t, err, "stat non-existent file")
		require.Nil(t, stat, "stat non-existent file")
	})
}

func entryNames(entries []fs.DirEntry) []string {
	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name()
	}
	return names
}
