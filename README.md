# go-memfs

An in-memory `fs.FS` implementation for Go.

## Installation

```sh
go get github.com/spiretechnology/go-memfs
```

## Example Usage

```go
fs := memfs.FS{
    "foo.txt": memfs.File("Hello, world!"),
    "bar.txt": memfs.File("Goodbye, world!"),
    "foobar/baz.txt": memfs.File("Hello, again!"),
    "some/empty/dir": memfs.Dir{},
}
```
