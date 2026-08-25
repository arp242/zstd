package zos

import (
	"errors"
	"io/fs"
	"os"
	"runtime"
	"slices"
	"testing"
)

func TestCreateNew(t *testing.T) {
	tmp := t.TempDir() + "/new"

	// Create new
	fp, err := CreateNew(tmp, false)
	if err != nil {
		t.Fatal(err)
	}
	fp.Close()

	// 0 byte, don't re-open
	fp, err = CreateNew(tmp, false)
	if err == nil {
		t.Fatal("error is nil")
	}
	fp.Close()

	// 0 byte, re-open
	fp, err = CreateNew(tmp, true)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fp.WriteString("X"); err != nil {
		t.Fatal(err)
	}
	fp.Close()

	// Should error, as it has more than 0 bytes.
	fp, err = CreateNew(tmp, true)
	if err == nil {
		t.Fatal("error is nil")
	}
	fp.Close()

	have, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if string(have) != "X" {
		t.Fatal(string(have))
	}
}

func TestCreateAtomic(t *testing.T) {
	list := func(t *testing.T, tmp string, want ...string) {
		t.Helper()
		ls, err := os.ReadDir(tmp)
		if err != nil {
			t.Fatal(err)
		}

		have := make([]string, 0, len(ls))
		for _, f := range ls {
			have = append(have, f.Name())
		}
		slices.Sort(want)
		slices.Sort(have)
		if !slices.Equal(have, want) {
			t.Errorf("\nhave: %#v\nwant: %#v", have, want)
		}
	}

	t.Run("works", func(t *testing.T) {
		tmp := t.TempDir()
		fp, err := CreateAtomic(tmp + "/file")
		if err != nil {
			t.Fatal(err)
		}
		defer fp.Close()
		if _, err := fp.WriteString("abc"); err != nil {
			t.Fatal(err)
		}
		list(t, tmp, "file-tmp")
		if err := fp.Close(); err != nil {
			t.Fatal(err)
		}
		list(t, tmp, "file")

		if h := fp.Name(); h != tmp+"/file" {
			t.Error(h)
		}
	})

	t.Run("removetemp", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			// remove C:\Users\RUNNER~1\AppData\Local\Temp\TestCreateAtomicremovetemp501911070\001/file-tmp:
			// The process cannot access the file because it is being used by another process.
			t.Skip()
		}

		tmp := t.TempDir()
		fp, err := CreateAtomic(tmp + "/file")
		if err != nil {
			t.Fatal(err)
		}
		defer fp.Close()

		if _, err := fp.WriteString("abc"); err != nil {
			t.Fatal(err)
		}
		list(t, tmp, "file-tmp")

		if err := fp.RemoveTemp(); err != nil {
			t.Fatal(err)
		}
		list(t, tmp)
		if err := fp.RemoveTemp(); !errors.Is(err, fs.ErrNotExist) {
			t.Fatal(err)
		}
	})
}
