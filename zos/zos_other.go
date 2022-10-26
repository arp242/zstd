//go:build plan9 || windows
// +build plan9 windows

package zos

import "os"

// Readable reports if the file is readable by the current user.
//
// This is a stub for non-POSIX systems which always returns true.
func Readable(s os.FileInfo) (bool, error) { return true, nil }

// Writable reports if the file is writable by the current user.
//
// This is a stub for non-POSIX systems which always returns true.
func Writable(s os.FileInfo) (bool, error) { return true, nil }
