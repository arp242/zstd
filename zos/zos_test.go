package zos

import (
	"os"
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

	have, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if string(have) != "X" {
		t.Fatal(string(have))
	}
}
