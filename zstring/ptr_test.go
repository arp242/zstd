package zstring

import (
	"fmt"
	"testing"
)

func TestPointer(t *testing.T) {
	test := func(ptr Pointer, want string) {
		func(p *string) {
			got := fmt.Sprintf("%p %s\n", p, *p)
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

	ptr := NewPointer("hello")
	test(ptr, fmt.Sprintf("%p hello\n", ptr.P))

	ptr.Set("X")
	test(ptr, fmt.Sprintf("%p X\n", ptr.P))
	if ptr.String() != "X" || ptr.Value() != "X" {
		t.Error()
	}

	ptr.P = nil
	if ptr.String() != "<nil>" || ptr.Value() != "" {
		t.Error()
	}
}
