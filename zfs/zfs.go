// Package zfs provides some extensions to the fs package.
package zfs

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

// Exists reports if this file or directory exists in the given fs.
func Exists(fsys fs.FS, name string) bool {
	ls, err := fs.ReadDir(fsys, filepath.Dir(name))
	if err != nil {
		return false
	}

	base := filepath.Base(name)
	for _, l := range ls {
		if l.Name() == base {
			return true
		}
	}

	return false
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

type OverlayFS struct{ base, overlay fs.FS }

// Open implements the fs.FS interface.
func (o OverlayFS) Open(name string) (fs.File, error) {
	f, err := o.overlay.Open(name)
	if err == nil {
		return f, err
	}
	return o.base.Open(name)
}

// ReadDir implements the fs.ReadDirFS interface.
func (o OverlayFS) ReadDir(name string) ([]fs.DirEntry, error) {
	var ls []fs.DirEntry

	if rd, ok := o.base.(fs.ReadDirFS); ok {
		l, err := rd.ReadDir(name)
		if err != nil {
			return nil, fmt.Errorf("OverlayFS: base: %w", err)
		}
		ls = l
	}

	if rd, ok := o.overlay.(fs.ReadDirFS); ok {
		l, err := rd.ReadDir(name)
		if err != nil {
			return nil, fmt.Errorf("OverlayFS: overlay: %w", err)
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

// InOverlay reports if this filename exists in the overlay part. It may or may
// not exist in the base part.
func (o OverlayFS) InOverlay(name string) bool { _, err := o.overlay.Open(name); return err == nil }

// InBase reports if this filename exists in the base part. It may or may
// not exist in the overlay part.
func (o OverlayFS) InBase(name string) bool { _, err := o.base.Open(name); return err == nil }

// OverlayFS returns a filesystem which reads from overlay, falling back to base
// if that fails.
func NewOverlayFS(base, overlay fs.FS) OverlayFS {
	return OverlayFS{base: base, overlay: overlay}
}
