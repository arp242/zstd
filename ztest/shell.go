package ztest

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

var join = filepath.Join

// mkdir
func Mkdir(t *testing.T, path ...string) {
	t.Helper()
	if len(path) < 1 {
		t.Fatalf("ztest.Mkdir: path must have at least one element: %s", path)
	}
	err := os.Mkdir(join(path...), 0o0755)
	if err != nil {
		t.Fatalf("ztest.Mkdir(%q): %s", join(path...), err)
	}
}

// mkdir -p
func MkdirAll(t *testing.T, path ...string) {
	t.Helper()
	if len(path) < 1 {
		t.Fatalf("ztest.MkdirAll: path must have at least one element: %s", path)
	}
	err := os.MkdirAll(join(path...), 0o0755)
	if err != nil {
		t.Fatalf("ztest.MkdirAll(%q): %s", join(path...), err)
	}
}

// ln
func Hardlink(t *testing.T, target string, link ...string) {
	t.Helper()
	if len(link) < 1 {
		t.Fatalf("ztest.Hardlink: link must have at least one element: %s", link)
	}
	err := os.Link(target, join(link...))
	if err != nil {
		t.Fatalf("ztest.Hardlink(%q, %q): %s", target, join(link...), err)
	}
}

// ln -s
func Symlink(t *testing.T, target string, link ...string) {
	t.Helper()
	if len(link) < 1 {
		t.Fatalf("ztest.Symlink: link must have at least one element: %s", link)
	}
	err := os.Symlink(target, join(link...))
	if err != nil {
		t.Fatalf("ztest.Symlink(%q, %q): %s", target, join(link...), err)
	}
}

// echo > and echo >>
func EchoAppend(t *testing.T, data string, path ...string) { t.Helper(); echo(t, false, data, path...) }
func EchoTrunc(t *testing.T, data string, path ...string)  { t.Helper(); echo(t, true, data, path...) }
func echo(t *testing.T, trunc bool, data string, path ...string) {
	n := "echoAppend"
	if trunc {
		n = "echoTrunc"
	}
	t.Helper()
	if len(path) < 1 {
		t.Fatalf("%s: path must have at least one element: %s", n, path)
	}

	err := func() error {
		var (
			fp  *os.File
			err error
		)
		if trunc {
			fp, err = os.Create(join(path...))
		} else {
			fp, err = os.OpenFile(join(path...), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		}
		if err != nil {
			return err
		}
		if err := fp.Sync(); err != nil {
			return err
		}
		if _, err := fp.WriteString(data); err != nil {
			return err
		}
		if err := fp.Sync(); err != nil {
			return err
		}
		return fp.Close()
	}()
	if err != nil {
		t.Fatalf("%s(%q): %s", n, join(path...), err)
	}
}

// touch
func Touch(t *testing.T, path ...string) {
	t.Helper()
	if len(path) < 1 {
		t.Fatalf("touch: path must have at least one element: %s", path)
	}
	fp, err := os.Create(join(path...))
	if err != nil {
		t.Fatalf("ztest.Touch(%q): %s", join(path...), err)
	}
	err = fp.Close()
	if err != nil {
		t.Fatalf("ztest.Touch(%q): %s", join(path...), err)
	}
}

// mv
func Mv(t *testing.T, src string, dst ...string) {
	t.Helper()
	if len(dst) < 1 {
		t.Fatalf("ztest.Mv: dst must have at least one element: %s", dst)
	}

	err := os.Rename(src, join(dst...))
	if err != nil {
		t.Fatalf("ztest.Mv(%q, %q): %s", src, join(dst...), err)
	}
}

// rm
func Rm(t *testing.T, path ...string) {
	t.Helper()
	if len(path) < 1 {
		t.Fatalf("ztest.Rm: path must have at least one element: %s", path)
	}
	err := os.Remove(join(path...))
	if err != nil {
		t.Fatalf("ztest.Rm(%q): %s", join(path...), err)
	}
}

// rm -r
func RmAll(t *testing.T, path ...string) {
	t.Helper()
	if len(path) < 1 {
		t.Fatalf("ztest.RmAll: path must have at least one element: %s", path)
	}
	err := os.RemoveAll(join(path...))
	if err != nil {
		t.Fatalf("ztest.RmAll(%q): %s", join(path...), err)
	}
}

// cat
func Cat(t *testing.T, path ...string) {
	t.Helper()
	if len(path) < 1 {
		t.Fatalf("ztest.Cat: path must have at least one element: %s", path)
	}
	_, err := os.ReadFile(join(path...))
	if err != nil {
		t.Fatalf("ztest.Cat(%q): %s", join(path...), err)
	}
}

// chmod
func Chmod(t *testing.T, mode fs.FileMode, path ...string) {
	t.Helper()
	if len(path) < 1 {
		t.Fatalf("ztest.Chmod: path must have at least one element: %s", path)
	}
	err := os.Chmod(join(path...), mode)
	if err != nil {
		t.Fatalf("ztest.Chmod(%q): %s", join(path...), err)
	}
}

// truncate
func Truncate(t *testing.T, sz int64, path ...string) {
	t.Helper()
	if len(path) < 1 {
		t.Fatalf("ztest.Truncate: path must have at least one element: %s", path)
	}
	fp, err := os.Create(join(path...))
	if err != nil {
		t.Fatalf("ztest.Truncate(%q): %s", join(path...), err)
	}
	if err := fp.Truncate(sz); err != nil {
		t.Fatalf("ztest.Truncate(%q): %s", join(path...), err)
	}
	if err := fp.Close(); err != nil {
		t.Fatalf("ztest.Truncate(%q): %s", join(path...), err)
	}
}
