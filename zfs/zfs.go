// Package zfs provides some extensions to the fs package.
package zfs

import (
	"embed"
	"io/fs"
	"os"
	"strings"

	"zgo.at/zstd/zgo"
)

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
