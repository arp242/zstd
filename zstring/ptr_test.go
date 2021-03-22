package zstring

import (
	"fmt"
	"testing"
)

func TestPtr(t *testing.T) {
	test := func(ptr Ptr, want string) {
		func(p *string) {
			got := fmt.Sprintf("%p %s\n", p, *p)
			if got != want {
				t.Error()
			}
		}(ptr.P)
		func(p Ptr) {
			got := fmt.Sprintf("%p %s\n", p.P, p)
			if got != want {
				t.Error()
			}
		}(ptr)
	}

	ptr := NewPtr("hello")
	test(ptr, fmt.Sprintf("%p hello\n", ptr.P))

	ptr.Set("X")
	test(ptr, fmt.Sprintf("%p X\n", ptr.P))
	if ptr.String() != "X" || ptr.Val() != "X" {
		t.Error()
	}

	ptr.P = nil
	if ptr.String() != "<nil>" || ptr.Val() != "" {
		t.Error()
	}
}
