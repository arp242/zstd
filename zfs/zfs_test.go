package zfs

import (
	"bytes"
	"io/fs"
	"reflect"
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

func TestOverlayFS(t *testing.T) {
	fsys := NewOverlayFS(
		fstest.MapFS{
			"both":      &fstest.MapFile{Data: []byte("both-base")},
			"base-only": &fstest.MapFile{Data: []byte("base-only")},
		},
		fstest.MapFS{
			"both":         &fstest.MapFile{Data: []byte("both-overlay")},
			"overlay-only": &fstest.MapFile{Data: []byte("overlay-only")},
		})

	if both, err := fs.ReadFile(fsys, "both"); err != nil {
		t.Fatal(err)
	} else if string(both) != "both-overlay" {
		t.Errorf("both: %q", string(both))
	}
	if baseOnly, err := fs.ReadFile(fsys, "base-only"); err != nil {
		t.Fatal(err)
	} else if string(baseOnly) != "base-only" {
		t.Errorf("base-only: %q", string(baseOnly))
	}
	if overlayOnly, err := fs.ReadFile(fsys, "overlay-only"); err != nil {
		t.Fatal(err)
	} else if string(overlayOnly) != "overlay-only" {
		t.Errorf("overlay-only: %q", string(overlayOnly))
	}

	ls, err := fs.ReadDir(fsys, ".")
	if err != nil {
		t.Fatal(err)
	}

	files := make([]string, 0, 3)
	for _, l := range ls {
		files = append(files, l.Name())
	}
	if !reflect.DeepEqual(files, []string{"base-only", "both", "overlay-only"}) {
		t.Error()
	}

	if !fsys.InOverlay("both") {
		t.Error()
	}
	if fsys.InOverlay("base-only") {
		t.Error()
	}
	if !fsys.InOverlay("overlay-only") {
		t.Error()
	}

	if !fsys.InBase("both") {
		t.Error()
	}
	if !fsys.InBase("base-only") {
		t.Error()
	}
	if fsys.InBase("overlay-only") {
		t.Error()
	}
}
