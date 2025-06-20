//go:build windows

package ztest

import (
	"testing"
)

// mkfifo
func Mkfifo(t *testing.T, path ...string) {
	t.Helper()
	t.Error("mkfifo not supported on Windows")
}

// mknod
func Mknod(t *testing.T, dev int, path ...string) {
	t.Helper()
	t.Error("mknod not supported on Windows")
}
