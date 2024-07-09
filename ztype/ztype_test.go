package ztype

import "testing"

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

func TestPtrOrNil(t *testing.T) {
	if have := Deref(PtrOrNil("hello"), "<NIL>"); have != "hello" {
		t.Error(have)
	}
	if have := Deref(PtrOrNil(""), "<NIL>"); have != "<NIL>" {
		t.Error(have)
	}
}
