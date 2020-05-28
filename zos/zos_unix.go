// +build unix

package zos

import (
	"fmt"
	"os"
	"syscall"
)

// Readable reports if the file is readable by the current user.
func Readable(s os.FileInfo) (bool, error) {
	stat, ok := s.Sys().(*syscall.Stat_t)
	if !ok {
		return false, fmt.Errorf("zosutil.Readable: assert to syscall.Stat_t failed; platform not supported?")
	}

	perm := ReadPermissions(s.Mode())

	if int(stat.Uid) == os.Geteuid() && perm.User.Read {
		return true, nil
	}

	gids, err := os.Getgroups()
	if err != nil {
		return false, fmt.Errorf("zosutil.Readable: %w", err)
	}
	for _, gid := range gids {
		if int(stat.Gid) == gid {
			return perm.Group.Read, nil
		}
	}

	return perm.Other.Read, nil
}

// Writable reports if the file is writable by the current user.
func Writable(s os.FileInfo) (bool, error) {
	stat, ok := s.Sys().(*syscall.Stat_t)
	if !ok {
		return false, fmt.Errorf("zosutil.Writable: assert to syscall.Stat_t failed; platform not supported?")
	}

	perm := ReadPermissions(s.Mode())

	if int(stat.Uid) == os.Geteuid() && perm.User.Write {
		return true, nil
	}

	gids, err := os.Getgroups()
	if err != nil {
		return false, fmt.Errorf("zosutil.Writable: %w", err)
	}
	for _, gid := range gids {
		if int(stat.Gid) == gid {
			return perm.Group.Read, nil
		}
	}

	return perm.Other.Write, nil
}
