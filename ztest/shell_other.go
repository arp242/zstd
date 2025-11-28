//go:build !unix

package ztest

import "errors"

func mkfifo(path string, mode uint32) error         { return errors.ErrUnsupported }
func mknod(path string, mode uint32, dev int) error { return errors.ErrUnsupported }
