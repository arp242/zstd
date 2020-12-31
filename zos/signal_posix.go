// +build !windows

package zos

import (
	"os"
	"syscall"
)

const (
	SIGUSR1 os.Signal = syscall.SIGUSR1
	SIGUSR2 os.Signal = syscall.SIGUSR2
)
