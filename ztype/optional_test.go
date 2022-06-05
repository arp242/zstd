package zstd

import (
	"fmt"
	"testing"

	"zgo.at/zstd/zjson"
)

func TestPtr(t *testing.T) {
	s := "hello"
	sp := Ptr(s)

	if have := Deref(sp, "NIL"); have != "hello" {
		t.Error(have)
	}

	sp = nil
	if have := Deref(sp, "NIL"); have != "NIL" {
		t.Error(have)
	}
}

func TestOptional(t *testing.T) {
	type Struct struct {
		A Optional[string]
		B Optional[string] `json:"bbb"`
	}

	s := Struct{
		A: NewOptional("hello"),
	}

	fmt.Println(s)
	fmt.Printf("%#v\n", s)

	j := zjson.MustMarshalIndent(s, "", "  ")
	fmt.Println(string(j))

	var s2 Struct
	zjson.MustUnmarshal(j, &s2)
	fmt.Println(s2)

	fmt.Println(s2.A.Get())
	s2.A.Set("XXX")
	fmt.Println(s2.A.Get())
	s2.A.Unset()
	fmt.Println(s2.A.Get())
}
