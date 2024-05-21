package zint

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestJoin(t *testing.T) {
	cases := []struct {
		in       []int
		expected string
	}{
		{
			[]int{1, 2, 3, 4, 4, 5, 6, 6, 6, 6, 7, 8, 8, 8},
			"1, 2, 3, 4, 4, 5, 6, 6, 6, 6, 7, 8, 8, 8",
		},
		{
			[]int{-1, -2, -3, -4, -4, -5, -6, -6, -6, -6, -7, -8, -8, -8},
			"-1, -2, -3, -4, -4, -5, -6, -6, -6, -6, -7, -8, -8, -8",
		},
		{
			[]int{},
			"",
		},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := Join(tt.in, ", ")
			if got != tt.expected {
				t.Errorf("\nwant: %q\ngot:  %q", tt.expected, got)
			}
		})
	}
}

func TestRange(t *testing.T) {
	cases := []struct {
		start, end int
		want       []int
	}{
		{1, 5, []int{1, 2, 3, 4, 5}},
		{0, 5, []int{0, 1, 2, 3, 4, 5}},
		{-2, 5, []int{-2, -1, 0, 1, 2, 3, 4, 5}},
		{-5, -1, []int{-5, -4, -3, -2, -1}},
		{100, 105, []int{100, 101, 102, 103, 104, 105}},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("%v-%v", tt.start, tt.end), func(t *testing.T) {
			out := Range(tt.start, tt.end)
			if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestInt(t *testing.T) {
	i := Int(42)

	if v := i.String(); v != "42" {
		t.Errorf("String: %#v", v)
	}
	if v := i.Int(); v != int(42) {
		t.Errorf("Int: %#v", v)
	}
	if v := i.Int64(); v != int64(42) {
		t.Errorf("Int64: %#v", v)
	}
	if v := i.Float32(); v != float32(42) {
		t.Errorf("Float32: %#v", v)
	}
	if v := i.Float64(); v != float64(42) {
		t.Errorf("Float64: %#v", v)
	}

	v, err := json.Marshal(i)
	if err != nil {
		t.Fatal(err)
	}
	if string(v) != "42" {
		t.Errorf("json: %#v", v)
	}
}

func TestToIntSlice(t *testing.T) {
	tests := []struct {
		in   any
		ok   bool
		want []int64
	}{
		{"", false, nil},
		{1, false, nil},
		{nil, false, nil},

		{[]int(nil), true, []int64{}},
		{[]int{}, true, []int64{}},
		{[]int{1}, true, []int64{1}},
		{[]int{-1}, true, []int64{-1}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have, ok := ToIntSlice(tt.in)
			if ok != tt.ok {
				t.Errorf("\nhave: %t\nwant: %t", ok, tt.ok)
			}
			if !reflect.DeepEqual(have, tt.want) {
				t.Errorf("\nhave: %#v\nwant: %#v", have, tt.want)
			}
		})
	}
}

func TestRoundToPowerOf2(t *testing.T) {
	tests := []struct {
		in, want uint64
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 8},
		{math.MaxUint32, math.MaxUint32 + 1},
		{math.MaxUint32 + 2, 8589934592},
		{math.MaxUint64, 0}, // Overflows and wraps to 0
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := RoundToPowerOf2(tt.in)
			if have != tt.want {
				t.Errorf("\nhave: %d\nwant: %d", have, tt.want)
			}
		})
	}
}
