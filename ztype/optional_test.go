package ztype

import (
	"fmt"
	"testing"

	"zgo.at/zstd/zjson"
)

func TestOptional(t *testing.T) {
	type Struct struct {
		A Optional[string]
		B Optional[string] `json:"bbb"`
		C Optional[int]
	}

	s := Struct{A: NewOptional("hello")}

	if h, ok := s.A.Get(); h != "hello" || !ok {
		t.Error()
	}
	if h, ok := s.B.Get(); h != "" || ok {
		t.Error()
	}
	if h, ok := s.C.Get(); h != 0 || ok {
		t.Error()
	}

	t.Run("string", func(t *testing.T) {
		want := "{hello <not set> <not set>}"
		have := fmt.Sprintf("%s", s)
		if have != want {
			t.Errorf("\nhave: %q\nwant: %q", have, want)
		}
	})

	t.Run("json", func(t *testing.T) {
		have := string(zjson.MustMarshal(s))
		want := `{"A":"hello","bbb":null,"C":null}`
		if have != want {
			t.Errorf("\nhave: %q\nwant: %q", have, want)
		}

		var s2 Struct
		zjson.MustUnmarshal([]byte(have), &s2)
		have = fmt.Sprintf("%s", s2)
		want = `{hello <not set> <not set>}`
		if have != want {
			t.Errorf("\nhave: %q\nwant: %q", have, want)
		}

		if h, ok := s2.A.Get(); h != "hello" || !ok {
			t.Error()
		}
		if h, ok := s2.B.Get(); h != "" || ok {
			t.Error()
		}
		if h, ok := s2.C.Get(); h != 0 || ok {
			t.Error()
		}
	})

	t.Run("sql", func(t *testing.T) {
		// TODO: test this too; looks like we need to implement a "fake" DB
		// driver, similar to what's in fakedb_test.go
		// Should probably add something like that to ztest or zsql.
	})

	t.Run("set", func(t *testing.T) {
		s.A.Set("XXX")
		if h, ok := s.A.Get(); h != "XXX" || !ok {
			t.Error()
		}

		s.A.Unset()
		if h, ok := s.A.Get(); h != "" || ok {
			t.Error(s.A.Get())
		}

		s.C.Set(42)
		if h, ok := s.C.Get(); h != 42 || !ok {
			t.Error()
		}

		s.C.Unset()
		if h, ok := s.C.Get(); h != 0 || ok {
			t.Error()
		}
	})
}
