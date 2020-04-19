package osutil

import "os"

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
