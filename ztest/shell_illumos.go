//go:build solaris

package ztest

import (
	"errors"
	"syscall"
)

// TODO: there's no entry for this in the syscall package, don't want to depend
// on x/sys/unix here, and not so easy to copy the relevant bits to here. Just
// error out for now as it's a little-used feature for a little-used system.
func mkfifo(path string, mode uint32) error         { return errors.ErrUnsupported }
func mknod(path string, mode uint32, dev int) error { return syscall.Mknod(path, mode, dev) }
