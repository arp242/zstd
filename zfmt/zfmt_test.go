package zfmt

import "testing"

type typ int8

func TestBinary(t *testing.T) {
	i := int8(7)
	tests := []struct {
		in   any
		want string
	}{
		{int8(0), "0000_0000"},
		{int8(7), "0000_0111"},
		{int8(127), "0111_1111"},
		{int8(-127), "1111_1111"},
		{typ(7), "0000_0111"},
		{&i, "0000_0111"},

		{uint8(0), "0000_0000"},
		{uint8(7), "0000_0111"},
		{uint8(255), "1111_1111"},

		{uint16(255), "0000_0000 1111_1111"},
		{int16(255), "0000_0000 1111_1111"},
		{int16(-255), "1000_0000 1111_1111"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := Binary(tt.in)
			if got != tt.want {
				t.Errorf("\ngot:  %q\nwant: %q", got, tt.want)
			}
		})
	}
}
