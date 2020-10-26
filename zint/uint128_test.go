package zint

import (
	"fmt"
	"testing"
)

func TestUint128(t *testing.T) {
	//ztest.MustInline(t, "zgo.at/zstd/zint",
	//	"NewUint128", "Uint128.IsZero", "Uint128.String")

	tests := []struct {
		in               []byte
		wantB10, wantB16 string
		wantZero         bool
	}{
		{
			[]byte{44, 25, 67, 129, 231, 4, 77, 72, 157, 135, 42, 180, 126, 162, 176, 131},
			"3177645237292256584-11351088340517695619", "2c194381e7044d48-9d872ab47ea2b083", false,
		},
		{
			[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			"0-0", "0-0", true,
		},
		{
			[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			"0-1", "0-1", false,
		},
		{
			[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
			"0-256", "0-100", false,
		},
		{
			[]byte{0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255},
			"0-18446744073709551615", "0-ffffffffffffffff", false,
		},
		{
			[]byte{0, 0, 0, 0, 0, 0, 0, 1, 255, 255, 255, 255, 255, 255, 255, 255},
			"1-18446744073709551615", "1-ffffffffffffffff", false,
		},
		{
			[]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			"18446744073709551615-18446744073709551615",
			"ffffffffffffffff-ffffffffffffffff",
			false,
		},
		{
			[]byte{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
			"1-0", "1-0", false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			i, err := NewUint128(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if !tt.wantZero && i.IsZero() {
				t.Fatal("IsZero()")
			}

			b := i.Bytes()
			if string(tt.in) != string(b) {
				t.Errorf("Bytes()\nwant: %v\ngot:  %v", tt.in, b)
			}

			if i.Format(10) != tt.wantB10 {
				t.Errorf("Format(10)\nwant: %s\ngot:  %s", tt.wantB10, i.Format(10))
			}
			if i.String() != tt.wantB16 {
				t.Errorf("String()\nwant: %s\ngot:  %s", tt.wantB16, i.String())
			}

			p10, err := ParseUint128(i.Format(10), 10)
			if err != nil {
				t.Error(err)
			} else if p10.Format(10) != tt.wantB10 {
				t.Errorf("Parse(10)\nwant: %s\ngot:  %s", tt.wantB10, p10.Format(10))
			}

			p16, err := ParseUint128(i.Format(16), 16)
			if err != nil {
				t.Error(err)
			} else if p16.Format(16) != tt.wantB16 {
				t.Errorf("Parse(16)\nwant: %s\ngot:  %s", tt.wantB16, p16.Format(16))
			}
		})
	}

	t.Run("invalid", func(t *testing.T) {
		for _, invalid := range [][]byte{nil, {1}, {1}} {
			i, err := NewUint128(invalid)
			if err == nil {
				t.Fatal(err)
			}
			if i[0] != 0 || i[1] != 0 {
				t.Fatal(err)
			}
		}
	})

	t.Run("zero", func(t *testing.T) {
		i2 := Uint128{}
		if !i2.IsZero() {
			t.Fatal()
		}
	})
}
