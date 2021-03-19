package zbyte

import (
	"reflect"
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
			got := Binary(tt.in)
			if got != tt.want {
				t.Errorf("want: %t; got: %t", tt.want, got)
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
			got := ElideLeft(tt.in, tt.n)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("want: %s; got: %s", tt.want, got)
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
			got := ElideRight(tt.in, tt.n)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("want: %s; got: %s", tt.want, got)
			}
		})
	}
}
