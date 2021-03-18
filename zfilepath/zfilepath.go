// Package zfilepath implements functions for manipulating filename paths.
package zfilepath

import "path/filepath"

// SplitExt returns the path without extension and the extension.
//
// If there is no extension the first return value is the same as the input. The
// extension will *not* contain a dot.
func SplitExt(path string) (string, string) {
	e := filepath.Ext(path)
	if e == "" {
		return path, ""
	}
	return path[:len(path)-len(e)], e[1:]
}
