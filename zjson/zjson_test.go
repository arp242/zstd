package zjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type ts struct {
	Ts Timestamp `json:"ts"`
}

func TestTimestamp(t *testing.T) {
	var x ts
	MustUnmarshal([]byte(`{"ts": 1234567890}`), &x)

	out := fmt.Sprintf("%v", x)
	want := "{2009-02-13 23:31:30 +0000 UTC}"
	if out != want {
		t.Errorf("Unmarshal\nout:  %q\nwant: %q", out, want)
	}

	out2 := string(MustMarshal(x))
	want2 := `{"ts":1234567890}`
	if out2 != want2 {
		t.Errorf("Marshal\nout:  %q\nwant: %q", out2, want2)
	}

	err := json.Unmarshal([]byte(`{"ts": "NaN"}`), &x)
	if err == nil {
		t.Errorf("no error on NaN")
	}
}

func TestMustMarshal(t *testing.T) {
	cases := []struct {
		in   string
		want []byte
	}{
		{`Hello`, []byte(`"Hello"`)},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := MustMarshal(tc.in)
			if !reflect.DeepEqual(out, tc.want) {
				t.Errorf("\nout:  %s\nwant: %s\n", out, tc.want)
			}
		})
	}
}

func TestMustMarshalIndent(t *testing.T) {
	cases := []struct {
		in   map[string]string
		want []byte
	}{
		{map[string]string{"hello": "world", "a": "b"}, []byte("{\n  \"a\": \"b\",\n  \"hello\": \"world\"\n}")},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := MustMarshalIndent(tc.in, "", "  ")
			if !reflect.DeepEqual(out, tc.want) {
				t.Errorf("\nout:  %s\nwant: %s\n", out, tc.want)
			}
		})
	}
}

func TestMustFormat(t *testing.T) {
	cases := []struct {
		in   []byte
		want []byte
	}{
		{[]byte(`{"hello": "world", "a": "b"}`), []byte("{\n  \"a\": \"b\",\n  \"hello\": \"world\"\n}")},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := MustIndent(tc.in, &map[string]string{}, "", "  ")
			if !reflect.DeepEqual(out, tc.want) {
				t.Errorf("\nout:  %s\nwant: %s\n", out, tc.want)
			}
		})
	}
}

func TestMustUnmarshal(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var out struct {
			Hello string `json:"hello"`
		}
		MustUnmarshal([]byte(`{"hello":"world"}`), &out)
		if out.Hello != "world" {
			t.Errorf("%#v", out)
		}
	})

	t.Run("panic", func(t *testing.T) {
		defer func() {
			rec := recover()
			if rec == nil {
				t.Errorf("no panic?")
			}
		}()

		var out struct {
			Hello time.Time `json:"hello"`
		}
		MustUnmarshal([]byte(`{"hello":"world"}`), &out)
	})

}
