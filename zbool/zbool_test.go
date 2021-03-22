package zbool

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"testing"
	"time"

	"zgo.at/zstd/ztest"
)

func TestBool(t *testing.T) {
	t.Run("Bool", func(t *testing.T) {
		b := interface{ Bool() bool }(Bool(false))
		if b.Bool() {
			t.Error()
		}

		b = interface{ Bool() bool }(Bool(true))
		if !b.Bool() {
			t.Error()
		}
	})

	t.Run("value", func(t *testing.T) {
		cases := []struct {
			in   Bool
			want driver.Value
		}{
			{false, int64(0)},
			{true, int64(1)},
		}

		for _, tc := range cases {
			t.Run(fmt.Sprintf("%t", tc.in), func(t *testing.T) {
				out, err := tc.in.Value()
				if err != nil {
					t.Fatal(err)
				}
				if out != tc.want {
					t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
				}
			})
		}
	})

	t.Run("scan", func(t *testing.T) {
		cases := []struct {
			in      interface{}
			want    Bool
			wantErr string
		}{
			{[]byte("true"), true, ""},
			{float64(1.0), true, ""},
			{[]byte{1}, true, ""},
			{int64(1), true, ""},
			{"true", true, ""},
			{true, true, ""},
			{"1", true, ""},

			{[]byte("false"), false, ""},
			{float64(0), false, ""},
			{[]byte{0}, false, ""},
			{int64(0), false, ""},
			{"false", false, ""},
			{false, false, ""},
			{"0", false, ""},
			{nil, false, ""},

			{"not a valid bool", false, "invalid value \"not a valid bool\""},
			{time.Time{}, false, "unsupported type time.Time"},
		}

		for _, tc := range cases {
			t.Run(fmt.Sprintf("%s", tc.in), func(t *testing.T) {
				var out Bool
				err := out.Scan(tc.in)
				if !ztest.ErrorContains(err, tc.wantErr) {
					t.Errorf("\nout:  %#v\nwant: %#v\n", err, tc.wantErr)
				}
				if !reflect.DeepEqual(out, tc.want) {
					t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
				}
			})
		}
	})

	t.Run("marshalText", func(t *testing.T) {
		cases := []struct {
			in      Bool
			want    []byte
			wantErr string
		}{
			{false, []byte("false"), ""},
			{true, []byte("true"), ""},
		}

		for _, tc := range cases {
			t.Run(fmt.Sprintf("%t", tc.in), func(t *testing.T) {
				out, err := tc.in.MarshalText()
				if !ztest.ErrorContains(err, tc.wantErr) {
					t.Errorf("\nout:  %#v\nwant: %#v\n", err, tc.wantErr)
				}
				if !reflect.DeepEqual(out, tc.want) {
					t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
				}
			})
		}
	})

	t.Run("unmarshalText", func(t *testing.T) {
		cases := []struct {
			in      []byte
			want    Bool
			wantErr string
		}{
			{[]byte("  true  "), true, ""},
			{[]byte(` "true"`), true, ""},
			{[]byte(`  1 `), true, ""},
			{[]byte("false  "), false, ""},
			{[]byte(`"false" `), false, ""},
			{[]byte(` 0 `), false, ""},
			{[]byte(`not a valid bool`), false, "invalid value \"not a valid bool\""},
		}

		for _, tc := range cases {
			t.Run(fmt.Sprintf("%s", tc.in), func(t *testing.T) {
				var out Bool
				err := out.UnmarshalText(tc.in)
				if !ztest.ErrorContains(err, tc.wantErr) {
					t.Errorf("\nout:  %#v\nwant: %#v\n", err, tc.wantErr)
				}
				if !reflect.DeepEqual(out, tc.want) {
					t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
				}
			})
		}
	})
}
