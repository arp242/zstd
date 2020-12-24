package zint

import (
	"encoding"
	"encoding/json"
	"fmt"
	"testing"
)

var (
	_ encoding.TextUnmarshaler = new(Bitflag8)
	_ json.Unmarshaler         = new(Bitflag8)
)

func TestBitflag(t *testing.T) {
	const (
		Foo Bitflag8 = 1 << iota
		Bar
		FooBar
	)

	var d Bitflag8
	test := func(want string) {
		got := fmt.Sprint(d.Has(Foo), d.Has(Bar), d.Has(FooBar))
		if got != want {
			t.Errorf("\ngot:  %s\nwant: %s", got, want)
		}
	}

	test("false false false")

	d.Set(Foo)
	d.Toggle(FooBar)
	test("true false true")

	d.Clear(FooBar)
	test("true false false")

	t.Run("json", func(t *testing.T) {
		f := Foo | Bar
		j, err := json.Marshal(f)
		if err != nil {
			t.Fatal(err)
		}
		if string(j) != "3" {
			t.Errorf(string(j))
		}

		var nf Bitflag8
		err = json.Unmarshal(j, &nf)
		if err != nil {
			t.Fatal(err)
		}

		if uint8(nf) != 3 {
			t.Error(nf)
		}
	})
}
