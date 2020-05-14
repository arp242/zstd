package intutil

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNonZero(t *testing.T) {
	cases := []struct {
		a, b int64
		c    []int64
		want int64
	}{
		{0, 0, nil, 0},
		{0, 0, []int64{0, 0}, 0},

		{42, 2, nil, 42},
		{0, 43, nil, 43},
		{0, 0, []int64{5, 0}, 5},
		{0, 0, []int64{6, 6}, 6},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%b", i), func(t *testing.T) {
			out := NonZero(tt.a, tt.b, tt.c...)
			if out != tt.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestFilterInt(t *testing.T) {
	cases := []struct {
		fun  func(int64) bool
		in   []int64
		want []int64
	}{
		{
			FilterIntEmpty,
			[]int64(nil),
			[]int64(nil),
		},
		{
			FilterIntEmpty,
			[]int64{},
			[]int64(nil),
		},
		{
			FilterIntEmpty,
			[]int64{1},
			[]int64{1},
		},
		{
			FilterIntEmpty,
			[]int64{0, 1, 0},
			[]int64{1},
		},
		{
			FilterIntEmpty,
			[]int64{0, 1, 0, 2, -1, 0, 0, 0, 42, 666, -666, 0, 0, 0},
			[]int64{1, 2, -1, 42, 666, -666},
		},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := Filter(tt.in, tt.fun)
			if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	cases := []struct {
		in       []int64
		expected string
	}{
		{
			[]int64{1, 2, 3, 4, 4, 5, 6, 6, 6, 6, 7, 8, 8, 8},
			"1, 2, 3, 4, 4, 5, 6, 6, 6, 6, 7, 8, 8, 8",
		},
		{
			[]int64{-1, -2, -3, -4, -4, -5, -6, -6, -6, -6, -7, -8, -8, -8},
			"-1, -2, -3, -4, -4, -5, -6, -6, -6, -6, -7, -8, -8, -8",
		},
		{
			[]int64{},
			"",
		},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := Join(tt.in, ", ")
			if got != tt.expected {
				t.Errorf("want: %q\ngot:  %q", tt.expected, got)
			}
		})
	}
}

func TestUniq(t *testing.T) {
	cases := []struct {
		in       []int64
		expected []int64
	}{
		{
			[]int64{1, 2, 3, 4, 4, 5, 6, 6, 6, 6, 7, 8, 8, 8},
			[]int64{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			[]int64{1, 3, 8, 3, 8},
			[]int64{1, 3, 8},
		},
		{
			[]int64{1, 2, 3},
			[]int64{1, 2, 3},
		},
		{
			[]int64{},
			nil,
		},
		{
			nil,
			nil,
		},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := Uniq(tt.in)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("want: %q\ngot:  %q", tt.expected, got)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		list     []int
		find     int
		expected bool
	}{
		{[]int{42}, 42, true},
		{[]int{42}, 4, false},
		{[]int{42, 666, 14159}, 666, true},
		{[]int{42, 666, 14159}, 0, false},
		{[]int{}, 0, false},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := Contains(tt.list, tt.find)
			if got != tt.expected {
				t.Errorf("want: %#v\ngot:  %#v", tt.expected, got)
			}
		})
	}
}

func TestContains64(t *testing.T) {
	tests := []struct {
		list     []int64
		find     int64
		expected bool
	}{
		{[]int64{42}, 42, true},
		{[]int64{42}, 4, false},
		{[]int64{42, 666, 14159}, 666, true},
		{[]int64{42, 666, 14159}, 0, false},
		{[]int64{}, 0, false},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := Contains64(tt.list, tt.find)
			if got != tt.expected {
				t.Errorf("want: %#v\ngot:  %#v", tt.expected, got)
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

func TestDifference(t *testing.T) {
	tests := []struct {
		inSet    []int64
		inOthers [][]int64
		want     []int64
	}{
		{[]int64{}, [][]int64{}, []int64{}},
		{nil, [][]int64{}, []int64{}},
		{[]int64{}, nil, []int64{}},
		{nil, nil, []int64{}},
		{[]int64{1}, [][]int64{{1}}, []int64{}},
		{[]int64{1, 2, 2, 3}, [][]int64{{1, 2, 2, 3}}, []int64{}},
		{[]int64{1, 2, 2, 3}, [][]int64{{1, 2}, {3}}, []int64{}},
		{[]int64{1, 2}, [][]int64{{1}}, []int64{2}},
		{[]int64{1, 2, 3}, [][]int64{{1}}, []int64{2, 3}},
		{[]int64{1, 2, 3}, [][]int64{{}, {1}}, []int64{2, 3}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := Difference(tt.inSet, tt.inOthers...)
			if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}
