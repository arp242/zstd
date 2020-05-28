package zbyte

import (
	"fmt"
	"testing"
)

func TestBinary(t *testing.T) {
	tests := []struct {
		in   []byte
		want bool
	}{
		{[]byte(""), false},
		{[]byte("€"), false},
		{[]byte{0x12}, true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%x", tt.in), func(t *testing.T) {
			out := Binary(tt.in)
			if out != tt.want {
				t.Errorf("want: %t; out: %t", tt.want, out)
			}
		})
	}
}
