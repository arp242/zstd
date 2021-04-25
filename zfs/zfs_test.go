package zfs

import (
	"bytes"
	"io/fs"
	"testing"
	"testing/fstest"

	"zgo.at/zstd/zfs/testdata"
)

func TestSubIfExists(t *testing.T) {
	tests := []struct {
		fsys fs.FS
		dir  string
		want string
	}{
		{fstest.MapFS{
			"file": {Data: []byte("XXX")},
		}, "x", "XXX"},
		{fstest.MapFS{
			"x/file": {Data: []byte("XXX")},
		}, "x", "XXX"},
		{fstest.MapFS{
			"x/y/z/file": {Data: []byte("XXX")},
		}, "x/y/z", "XXX"},
		{fstest.MapFS{
			"y/z/file": {Data: []byte("XXX")},
		}, "x/y/z", "XXX"},
		{fstest.MapFS{
			"z/file": {Data: []byte("XXX")},
		}, "x/y/z", "XXX"},
		{fstest.MapFS{
			"file": {Data: []byte("XXX")},
		}, "x/y/z", "XXX"},

		{fstest.MapFS{
			"y/file": {Data: []byte("XXX")},
		}, "x", ""},
		{fstest.MapFS{
			"x/x/y/z/file": {Data: []byte("XXX")},
		}, "x/y/z", ""},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			fsys, err := SubIfExists(tt.fsys, tt.dir)
			if err != nil {
				t.Fatal(err)
			}

			got, err := fs.ReadFile(fsys, "file")
			if tt.want == "" {
				if err == nil {
					t.Error("err is nil")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Errorf("\ngot:  %q\nwant: %q", string(got), tt.want)
			}
		})
	}
}

func TestEmbedOrDir(t *testing.T) {
	tests := []struct {
		dev bool
	}{
		{false},
		{true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			fsys, err := EmbedOrDir(testdata.Embed, "zfs", tt.dev)
			if err != nil {
				t.Fatal(err)
			}

			gotProd, errProd := fs.ReadFile(fsys, "embed.go")
			gotDev, errDev := fs.ReadFile(fsys, "zfs_test.go")

			if tt.dev && (errDev != nil || errProd == nil || gotProd != nil) {
				t.Fatalf("dev failed: errDev: %s; errProd: %s; gotProd: %s", errDev, errProd, gotProd)
			}
			if !tt.dev && (errDev == nil || errProd != nil || gotDev != nil) {
				t.Fatalf("!dev failed: errDev: %s; errProd: %s; gotProd: %s", errDev, errProd, gotProd)
			}

			if tt.dev && !bytes.HasPrefix(gotDev, []byte("package zfs\n")) {
				t.Error(string(gotDev))
			}
			if !tt.dev && !bytes.HasPrefix(gotProd, []byte("package testdata\n")) {
				t.Error(string(gotProd))
			}
		})
	}
}
