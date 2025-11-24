//go:build freebsd

package ztest

import (
	"syscall"
	"testing"
)

// mkfifo
func Mkfifo(t *testing.T, path ...string) {
	t.Helper()
	if len(path) < 1 {
		t.Fatalf("ztest.Mkfifo: path must have at least one element: %s", path)
	}
	err := syscall.Mkfifo(join(path...), 0o644)
	if err != nil {
		t.Fatalf("ztest.Mkfifo(%q): %s", join(path...), err)
	}
}

// mknod
func Mknod(t *testing.T, dev int, path ...string) {
	t.Helper()
	if len(path) < 1 {
		t.Fatalf("ztest.Mknod: path must have at least one element: %s", path)
	}
	err := syscall.Mknod(join(path...), 0o644, uint64(dev))
	if err != nil {
		t.Fatalf("ztest.Mknod(%d, %q): %s", dev, join(path...), err)
	}
}
