package zfilepath

import "testing"

func TestTrimPrefix(t *testing.T) {
	tests := []struct {
		path, prefix, want string
	}{
		{"", "", ""},
		{"/", "/", ""},
		{"/etc/passwd", "/etc", "passwd"},
		{"/etc/passwd", "/etc/", "passwd"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := TrimPrefix(tt.path, tt.prefix)

			if have != tt.want {
				t.Errorf("\nhave: %q\nwant: %q", have, tt.want)
			}
		})
	}
}

func TestSplitExt(t *testing.T) {
	tests := []struct {
		in, path, ext string
	}{
		{"", "", ""},
		{"a", "a", ""},
		{"a.b", "a", "b"},
		{"/a/b/c.foo.bar", "/a/b/c.foo", "bar"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			p, e := SplitExt(tt.in)

			if p != tt.path {
				t.Errorf("\nhave: %q\nwant: %q", p, tt.path)
			}
			if e != tt.ext {
				t.Errorf("\nhave: %q\nwant: %q", e, tt.ext)
			}
		})
	}
}
