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
			have := Binary(tt.in)
			if have != tt.want {
				t.Errorf("\nhave: %q\nwant: %q", have, tt.want)
			}
		})
	}
}

func TestNumber(t *testing.T) {
	tt := func(t *testing.T, have, want string) {
		t.Run("", func(t *testing.T) {
			if have != want {
				t.Errorf("\nhave: %q\nwant: %q", have, want)
			}
		})
	}

	tt(t, Number(32, ','), "32")
	tt(t, Number(3200, ','), "3,200")
	tt(t, Number(233200, ','), "233,200")
	tt(t, Number(1233200, ','), "1,233,200")

	tt(t, Number(123.0, ','), "123")
	tt(t, Number(123.1, ','), "123.1")
	tt(t, Number(123456.994, ','), "123,456.994")
	tt(t, Number(9123456.994, ','), "9,123,456.994")

	tt(t, Number(1233200, '.'), "1.233.200")
	tt(t, Number(9123456.994, '.'), "9.123.456,994")
	tt(t, Number(9123456.994, '_'), "9_123_456.994")
	tt(t, Number(9123456.994, '\''), "9'123'456.994")
	tt(t, Number(9123456.994, ' '), "9 123 456.994")
}
