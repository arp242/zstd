package zint

import (
	"fmt"
	"testing"
)

func TestPointer(t *testing.T) {
	test := func(ptr Pointer, want string) {
		func(p *int) {
			got := fmt.Sprintf("%p %d\n", p, *p)
			if got != want {
				t.Error()
			}
		}(ptr.P)
		func(p Pointer) {
			got := fmt.Sprintf("%p %s\n", p.P, p)
			if got != want {
				t.Error()
			}
		}(ptr)
	}

	ptr := NewPointer(42)
	test(ptr, fmt.Sprintf("%p 42\n", ptr.P))

	ptr.Set(666)
	test(ptr, fmt.Sprintf("%p 666\n", ptr.P))
	if ptr.String() != "666" || ptr.Value() != 666 {
		t.Error()
	}

	ptr.P = nil
	if ptr.String() != "<nil>" || ptr.Value() != 0 {
		t.Error()
	}
}

func TestPointer64(t *testing.T) {
	test := func(ptr Pointer64, want string) {
		func(p *int64) {
			got := fmt.Sprintf("%p %d\n", p, *p)
			if got != want {
				t.Error()
			}
		}(ptr.P)
		func(p Pointer64) {
			got := fmt.Sprintf("%p %s\n", p.P, p)
			if got != want {
				t.Error()
			}
		}(ptr)
	}

	ptr := NewPointer64(42)
	test(ptr, fmt.Sprintf("%p 42\n", ptr.P))

	ptr.Set(666)
	test(ptr, fmt.Sprintf("%p 666\n", ptr.P))
	if ptr.String() != "666" || ptr.Value() != 666 {
		t.Error()
	}

	ptr.P = nil
	if ptr.String() != "<nil>" || ptr.Value() != 0 {
		t.Error()
	}

}
