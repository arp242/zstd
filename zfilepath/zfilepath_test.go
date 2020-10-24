package zfilepath

import "testing"

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
				t.Errorf("\ngot:  %q\nwant: %q", p, tt.path)
			}
			if e != tt.ext {
				t.Errorf("\ngot:  %q\nwant: %q", e, tt.ext)
			}
		})
	}
}
