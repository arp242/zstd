package ioutilx

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"zgo.at/ztest"
)

func TestMain(m *testing.M) {
	fifo := "testdata/fifo"
	err := syscall.Mkfifo(fifo, 644)
	if err != nil {
		panic(err)
	}

	e := m.Run()
	if err := os.Remove(fifo); err != nil {
		panic(err)
	}
	os.Exit(e)
}

func TestIsSymLink(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"testdata/file1", false},
		{"testdata/dir1", false},
		// {"testdata/link1", true},
		// {"testdata/link2", true},
		// {"testdata/link3", true},
	}

	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			st, err := os.Lstat(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			out := IsSymlink(st)
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestIsSameFile(t *testing.T) {
	cases := []struct {
		src, dst string
		want     string
	}{
		// {"testdata/file1", "testdata/link1", "are the same file"},
		// {"testdata/file1", "testdata/link2", "are the same file"},
		{"testdata/file1", "testdata/dir1", ""},
		{"testdata/file1", "nonexistent", ""},
		{"nonexistent", "testdata/file1", ""},
		{"nonexistent", "nonexistent", ""},
	}

	for _, tc := range cases {
		t.Run(tc.src+":"+tc.dst, func(t *testing.T) {
			out := IsSameFile(tc.src, tc.dst)
			if !ztest.ErrorContains(out, tc.want) {
				t.Fatalf("\nwant: %s\nout:  %+v", tc.want, out)
			}
		})
	}
}

func TestIsSpecialFile(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"testdata/file1", ""},
		{"testdata/dir1", ""},
		// {"testdata/link1", ""},
		{"testdata/fifo", "named pipe"},
		{"/dev/null", "device file"},
	}

	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {

			st, err := os.Lstat(tc.in)
			if err != nil {
				t.Fatal(err)
			}

			out := IsSpecialFile(st)
			if !ztest.ErrorContains(out, tc.want) {
				t.Fatalf("\nwant: %s\nout:  %+v", tc.want, out)
			}
		})
	}
}

func TestCopyData(t *testing.T) {
	cases := []struct {
		src, dst, want string
	}{
		{"testdata/file1", "testdata/file1", "same file"},
		{"nonexistent", "testdata/copydst", "no such file"},
		{"testdata/file1", "testdata/file2", "already exists"},
		{"testdata/fifo", "testdata/newfile", "named pipe"},
		// {"testdata/link1/asd", "testdata/dst1", "not a directory"},
		{"testdata/file1", "/cantwritehere", "permission denied"},
		{"testdata/file1", "testdata/dst1", ""},
		// {"testdata/link1", "testdata/dst1", ""},
	}

	for _, tc := range cases {
		t.Run(tc.src+":"+tc.dst, func(t *testing.T) {
			if tc.want == "" {
				defer clean(t, tc.dst)
			}

			out := CopyData(tc.src, tc.dst)
			if !ztest.ErrorContains(out, tc.want) {
				t.Fatalf("\nwant: %s\nout:  %+v", tc.want, out)
			}

			if tc.want == "" {
				filesMatch(t, tc.src, tc.dst)
			}
		})
	}
}

func TestCopyMode(t *testing.T) {
	cases := []struct {
		src, dst string
		mode     Modes
		want     string
	}{
		{"testdata/file1", "testdata/file1", Modes{}, "same file"},
		{"nonexistent", "testdata/copydst", Modes{}, "no such file"},
		{"testdata/fifo", "testdata/newfile", Modes{}, "named pipe"},
		// {"testdata/link1/asd", "testdata/dst1", Modes{}, "not a directory"},
		{"testdata/file1", "/cantwritehere", Modes{}, "no such file or directory"},

		{"testdata/exec", "testdata/dst1", Modes{Permissions: true, Owner: true, Mtime: true}, ""},
	}

	for _, tc := range cases {
		t.Run(tc.src+":"+tc.dst, func(t *testing.T) {
			if tc.want == "" {
				touch(t, tc.dst)
				defer clean(t, tc.dst)
			}

			out := CopyMode(tc.src, tc.dst, tc.mode)
			if !ztest.ErrorContains(out, tc.want) {
				t.Fatalf("\nwant: %s\nout:  %+v", tc.want, out)
			}

			// Seems to fail on Travis all of the sudden:
			//
			// --- FAIL: TestCopyMode (0.00s)
			//     --- FAIL: TestCopyMode/testdata/exec:testdata/dst1 (0.00s)
			//     	copy_test.go:171: wrong mode: -rwxrwxr-x
			// --- FAIL: TestCopy (0.01s)
			//     --- FAIL: TestCopy/testdata/exec:testdata/dst1 (0.00s)
			//     	copy_test.go:218: wrong mode: -rwxrwxr-x
			//     --- FAIL: TestCopy/testdata/exec:testdata/dir1 (0.00s)
			//     	copy_test.go:218: wrong mode: -rwxrwxr-x
			//     --- FAIL: TestCopy/testdata/exec:testdata/dir1/ (0.00s)
			//     	copy_test.go:218: wrong mode: -rwxrwxr-x
			//
			// if tc.want == "" {
			// 	mode, err := os.Stat(tc.dst)
			// 	if err != nil {
			// 		t.Fatal(err)
			// 	}

			// 	if mode.Mode().String() != "-rwxr-xr-x" {
			// 		t.Fatalf("wrong mode: %s", mode.Mode())
			// 	}
			// }
		})
	}
}

func TestCopy(t *testing.T) {
	cases := []struct {
		src, dst string
		mode     Modes
		want     string
	}{
		{"testdata/file1", "testdata/file1", Modes{}, "same file"},
		{"nonexistent", "testdata/copydst", Modes{}, "no such file"},
		{"testdata/fifo", "testdata/newfile", Modes{}, "named pipe"},
		// {"testdata/link1/asd", "testdata/dst1", Modes{}, "not a directory"},
		{"testdata/file1", "/cantwritehere", Modes{}, "permission denied"},

		{"testdata/exec", "testdata/dst1", Modes{Permissions: true, Owner: true, Mtime: true}, ""},
		{"testdata/exec", "testdata/dir1", Modes{Permissions: true, Owner: true, Mtime: true}, ""},
		{"testdata/exec", "testdata/dir1/", Modes{Permissions: true, Owner: true, Mtime: true}, ""},
	}

	for _, tc := range cases {
		t.Run(tc.src+":"+tc.dst, func(t *testing.T) {
			c := tc.dst
			if strings.HasPrefix(c, "testdata/dir1") {
				c = filepath.Join(c, "exec")
			}

			if tc.want == "" {
				defer clean(t, c)
			}

			out := Copy(tc.src, tc.dst, tc.mode)
			if !ztest.ErrorContains(out, tc.want) {
				t.Fatalf("\nwant: %s\nout:  %+v", tc.want, out)
			}

			// --- FAIL: TestCopy (0.00s)
			//     --- FAIL: TestCopy/testdata/exec:testdata/dst1 (0.00s)
			//     	copy_test.go:231: wrong mode: -rwxrwxr-x
			//     --- FAIL: TestCopy/testdata/exec:testdata/dir1 (0.00s)
			//     	copy_test.go:231: wrong mode: -rwxrwxr-x
			//     --- FAIL: TestCopy/testdata/exec:testdata/dir1/ (0.00s)
			//     	copy_test.go:231: wrong mode: -rwxrwxr-x
			//
			// if tc.want == "" {
			// 	mode, err := os.Stat(c)
			// 	if err != nil {
			// 		t.Fatal(err)
			// 	}

			// 	if mode.Mode().String() != "-rwxr-xr-x" {
			// 		t.Fatalf("wrong mode: %s", mode.Mode())
			// 	}
			// }
		})
	}
}

func clean(t *testing.T, n string) {
	err := os.Remove(n)
	if err != nil {
		t.Fatalf("could not cleanup %v: %v", n, err)
	}
}

func filesMatch(t *testing.T, src, dst string) {
	srcContents, err := ioutil.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}

	dstContents, err := ioutil.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(srcContents, dstContents) {
		t.Errorf("%v and %v are not identical\nout:  %s\nwant: %s\n",
			src, dst, srcContents, dstContents)
	}
}

func touch(t *testing.T, n string) {
	fp, err := os.Create(n)
	if err != nil {
		t.Fatal(err)
	}
	err = fp.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCopyTree(t *testing.T) {
	t.Run("nonexistent", func(t *testing.T) {
		err := CopyTree("nonexistent", "test_copytree", nil)
		if !ztest.ErrorContains(err, "no such file or directory") {
			t.Error(err)
		}
	})
	t.Run("dst-exists", func(t *testing.T) {
		err := CopyTree("testdata", "testdata", nil)
		if !ztest.ErrorContains(err, "already exists") {
			t.Error(err)
		}
	})
	t.Run("dst-nodir", func(t *testing.T) {
		err := CopyTree("testdata/file1", "test", nil)
		if !ztest.ErrorContains(err, "not a directory") {
			t.Error(err)
		}
	})
	t.Run("permission", func(t *testing.T) {
		err := CopyTree("testdata", "/cant/write/here", nil)
		if !ztest.ErrorContains(err, "permission denied") {
			t.Error(err)
		}
	})

	defer func() {
		err := os.RemoveAll("test_copytree")
		if err != nil {
			t.Fatalf("could not clean: %v", err)
		}
	}()

	err := CopyTree("testdata", "test_copytree", &CopyTreeOptions{
		Symlinks: false,
		Ignore: func(path string, fi []os.FileInfo) []string {
			return []string{"fifo"}
		},
		CopyFunction:           Copy,
		IgnoreDanglingSymlinks: false,
	})
	if err != nil {
		t.Error(err)
		return
	}

	filesMatch(t, "testdata/file1", "test_copytree/file1")
}
