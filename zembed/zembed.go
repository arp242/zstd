// +build go1.16

// Package zembed implements functions for working with embeded files.
package zembed

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

// Dir will read from path if it's set, or from the embeded files if it's not.
func Dir(f embed.FS, path string) fs.FS {
	if path != "" {
		return dir(path)
	}
	return f
}

type dir string

func (d dir) Open(p string) (fs.File, error)          { return os.Open(filepath.Join(string(d), p)) }
func (d dir) ReadFile(p string) ([]byte, error)       { return os.ReadFile(filepath.Join(string(d), p)) }
func (d dir) ReadDir(p string) ([]fs.DirEntry, error) { return os.ReadDir(filepath.Join(string(d), p)) }
