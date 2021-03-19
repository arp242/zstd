package zint

import (
	"fmt"
	"testing"
)

func TestPtr(t *testing.T) {
	test := func(ptr Ptr, want string) {
		func(p *int) {
			got := fmt.Sprintf("%p %d\n", p, *p)
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

	ptr := NewPtr(42)
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

func TestPtr64(t *testing.T) {
	test := func(ptr Ptr64, want string) {
		func(p *int64) {
			got := fmt.Sprintf("%p %d\n", p, *p)
			if got != want {
				t.Error()
			}
		}(ptr.P)
		func(p Ptr64) {
			got := fmt.Sprintf("%p %s\n", p.P, p)
			if got != want {
				t.Error()
			}
		}(ptr)
	}

	ptr := NewPtr64(42)
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
