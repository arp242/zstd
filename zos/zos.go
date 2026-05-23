// Package zos implements functions for interfacing with the operating system.
package zos

import (
	"errors"
	"io/fs"
	"os"
	"strings"
)

const (
	Setuid uint32 = 1 << (12 - 1 - iota)
	Setgid
	Sticky
	UserRead
	UserWrite
	UserExecute
	GroupRead
	GroupWrite
	GroupExecute
	OtherRead
	OtherWrite
	OtherExecute
)

type Permissions struct {
	User  Permission
	Group Permission
	Other Permission
}

type Permission struct {
	Read    bool
	Write   bool
	Execute bool
}

// ReadPermissions reads all Unix permissions.
func ReadPermissions(mode os.FileMode) Permissions {
	m := uint32(mode)
	return Permissions{
		User: Permission{
			Read:    m&UserRead != 0,
			Write:   m&UserWrite != 0,
			Execute: m&UserExecute != 0,
		},
		Group: Permission{
			Read:    m&GroupRead != 0,
			Write:   m&GroupWrite != 0,
			Execute: m&GroupExecute != 0,
		},
		Other: Permission{
			Read:    m&OtherRead != 0,
			Write:   m&OtherWrite != 0,
			Execute: m&OtherExecute != 0,
		},
	}
}

// Arg gets the nth argument from os.Args, or an empty string if os.Args is too
// short.
func Arg(n int) string {
	if n > len(os.Args)-1 {
		return ""
	}
	return os.Args[n]
}

// CreateNew creates a new file.
//
// An error is returned if the file already exists, unless zeroByte is set it's
// a zero-byte file. Note that this is subject to race conditions, and not
// suitable if multiple processes may create the same file.
func CreateNew(path string, zeroByte bool) (*os.File, error) {
	out, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o666)
	if zeroByte && errors.Is(err, fs.ErrExist) {
		if st, stErr := os.Stat(path); stErr == nil && st.Size() == 0 {
			out, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o666)
		}
	}
	return out, err
}

type createAtomic struct{ *os.File }

// Close the [os.File], rendering it unusable for I/O. On files that support
// [File.SetDeadline], any pending I/O operations will be cancelled and return
// immediately with an [os.ErrClosed] error.
//
// The temporary file will be renamed to the destination.
//
// Close will return an error if it has already been called.
func (c createAtomic) Close() error {
	err := c.File.Close()
	if err != nil {
		return err
	}
	return os.Rename(c.File.Name(), c.Name())
}

// Name returns the name of the file as presented to Open.
//
// It is safe to call Name after [Close].
func (c createAtomic) Name() string {
	return strings.TrimSuffix(c.File.Name(), "-tmp")
}

// RemoveTmp removes the temporary file, if any.
func (c createAtomic) RemoveTemp() error {
	return os.Remove(c.File.Name())
}

// CreateAtomic is like [os.Create], but creates a temporary file (name +
// "-tmp"), which is moved to name on Close().
//
// The tmp file is not removed on errors, but the RemoveTemp() method can be
// used to clean up the -tmp file.
func CreateAtomic(name string) (createAtomic, error) {
	fp, err := os.Create(name + "-tmp")
	if err != nil {
		return createAtomic{}, err
	}
	return createAtomic{fp}, err
}
