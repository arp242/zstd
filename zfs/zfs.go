// Package zfs provides some extensions to the fs package.
package zfs

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"

	"zgo.at/zstd/zgo"
)

// MustReadFile is like fs.ReadFile, but will panic() on errors.
func MustReadFile(fsys fs.FS, name string) []byte {
	b, err := fs.ReadFile(fsys, name)
	if err != nil {
		panic("zfs.MustReadFile: " + err.Error())
	}
	return b
}

// EmbedOrDir returns e if dev is false, or a os.DirFS at the module root.
//
// This is intended to serve files from the local FS when developing and
// requires to be run in the source code directory.
//
// In both cases it will try to fs.Sub() to dir, but it's a non-fatal error if
// this directory doesn't exist.
func EmbedOrDir(e embed.FS, dir string, dev bool) (fs.FS, error) {
	var fsys fs.FS = e
	if dev {
		if r := zgo.ModuleRoot(); r != "" {
			fsys = os.DirFS(r)
		}
	}
	if dir == "" {
		return fsys, nil
	}
	return SubIfExists(fsys, dir)
}

// SubIfExists will fs.Sub() to a directory only if it exists.
//
// This will continue with the next directory if it doesn't exist. For example
// with "db/query", it will try to Sub to "db" first, and will continue to Sub
// to "query" regardless of whether the Sub to "db" succeeded or not, so
// "db/query" and "query" will end up beiing Sub()'d.
func SubIfExists(fsys fs.FS, dir string) (fs.FS, error) {
	for _, p := range strings.Split(dir, "/") {
		ls, err := fs.ReadDir(fsys, ".")
		if err != nil {
			return nil, err
		}

		for _, e := range ls {
			if e.IsDir() && e.Name() == p {
				fsys, err = fs.Sub(fsys, p)
				if err != nil {
					return nil, err
				}
				break
			}
		}
	}
	return fsys, nil
}

type overlayFS struct{ base, overlay fs.FS }

func (o overlayFS) Open(name string) (fs.File, error) {
	f, err := o.overlay.Open(name)
	if err == nil {
		return f, err
	}
	return o.base.Open(name)
}

func (o overlayFS) ReadDir(name string) ([]fs.DirEntry, error) {
	var ls []fs.DirEntry

	if rd, ok := o.base.(fs.ReadDirFS); ok {
		l, err := rd.ReadDir(name)
		if err != nil {
			return nil, fmt.Errorf("overlayFS: base: %w", err)
		}
		ls = l
	}

	if rd, ok := o.overlay.(fs.ReadDirFS); ok {
		l, err := rd.ReadDir(name)
		if err != nil {
			return nil, fmt.Errorf("overlayFS: overlay: %w", err)
		}

		merge := make(map[string]fs.DirEntry)
		for _, f := range ls {
			merge[f.Name()] = f
		}
		for _, f := range l {
			merge[f.Name()] = f
		}
		ls = make([]fs.DirEntry, 0, len(merge))
		for _, f := range merge {
			ls = append(ls, f)
		}
		sort.Slice(ls, func(i, j int) bool { return ls[i].Name() < ls[j].Name() })
	}

	return ls, nil
}

// OverlayFS returns a filesystem which reads from overlay, falling back to base
// if that fails.
func OverlayFS(base, overlay fs.FS) fs.FS {
	return overlayFS{base: base, overlay: overlay}
}
