package zint

import (
	"fmt"
	"testing"
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
}
