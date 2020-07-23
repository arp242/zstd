package zint

import (
	"testing"

	"zgo.at/zstd/ztest"
)

func TestUint128(t *testing.T) {
	ztest.MustInline(t, "zgo.at/zstd/zint",
		"NewUint128",
		"Uint128.IsZero",
		"Uint128.String",
		// Too complex:
		// "Uint128.Bytes",
		// "(*Uint128).New",
	)

	n := []byte{44, 25, 67, 129, 231, 4, 77, 72, 157, 135, 42, 180, 126, 162, 176, 131}

	i, err := NewUint128(n)
	if err != nil {
		t.Fatal(err)
	}

	b := i.Bytes()
	if string(n) != string(b) {
		t.Errorf("Bytes()\nwant: %v\ngot:  %v", n, b)
	}

	_, err = NewUint128(nil)
	if err == nil {
		t.Fatal(err)
	}

	if i.IsZero() {
		t.Fatal()
	}
	i2 := Uint128{}
	if !i2.IsZero() {
		t.Fatal()
	}
}
