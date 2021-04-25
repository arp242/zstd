package zbyte

import (
	"testing"
)

func TestBinary(t *testing.T) {
	tests := []struct {
		in   []byte
		want bool
	}{
		{[]byte(""), false},
		{[]byte("€"), false},
		{[]byte("helllo\x00"), true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := Binary(tt.in)
			if have != tt.want {
				t.Errorf("want: %t; have: %t", tt.want, have)
			}
		})
	}
}

func TestElideLeft(t *testing.T) {
	tests := []struct {
		in   []byte
		n    int
		want []byte
	}{
		{[]byte("abcdef"), 6, []byte("abcdef")},
		{[]byte("abcdef"), 2, []byte("ab")},
		{[]byte("abcdef"), 0, []byte("")},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := ElideLeft(tt.in, tt.n)
			if string(have) != string(have) {
				t.Errorf("want: %s; have: %s", tt.want, have)
			}
		})
	}
}

func TestElideRight(t *testing.T) {
	tests := []struct {
		in   []byte
		n    int
		want []byte
	}{
		{[]byte("abcdef"), 6, []byte("abcdef")},
		{[]byte("abcdef"), 2, []byte("ef")},
		{[]byte("abcdef"), 0, []byte("")},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := ElideRight(tt.in, tt.n)
			if string(have) != string(have) {
				t.Errorf("want: %s; have: %s", tt.want, have)
			}
		})
	}
}
