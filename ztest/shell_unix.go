//go:build unix && !freebsd && !solaris

package ztest

import "syscall"

func mkfifo(path string, mode uint32) error         { return syscall.Mkfifo(path, mode) }
func mknod(path string, mode uint32, dev int) error { return syscall.Mknod(path, mode, dev) }
