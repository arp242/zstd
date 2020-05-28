package zstring

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestFields(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"", nil},
		{"a", []string{"a"}},
		{"a;b", []string{"a", "b"}},
		{"  a  ;  b  ", []string{"a", "b"}},
		{"  a  ;  b  ; ", []string{"a", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			out := Fields(tt.in, ";")
			if !reflect.DeepEqual(out, tt.want) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestLeft(t *testing.T) {
	cases := []struct {
		in   string
		n    int
		want string
	}{
		{"Hello", 100, "Hello"},
		{"Hello", 1, "H…"},
		{"Hello", 5, "Hello"},
		{"Hello", 4, "Hell…"},
		{"Hello", 0, "…"},
		{"Hello", -2, "…"},
		{"汉语漢語", 1, "汉…"},
		{"汉语漢語", 3, "汉语漢…"},
		{"汉语漢語", 4, "汉语漢語"},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := Left(tt.in, tt.n)
			if out != tt.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestUpperFirst(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"hello", "Hello"},
		{"helloWorld", "HelloWorld"},
		{"h", "H"},
		{"hh", "Hh"},
		{"ëllo", "Ëllo"},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			out := UpperFirst(tt.in)
			if out != tt.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestLowerFirst(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Hello", "hello"},
		{"HelloWorld", "helloWorld"},
		{"H", "h"},
		{"HH", "hH"},
		{"Ëllo", "ëllo"},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			out := LowerFirst(tt.in)
			if out != tt.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestRemoveUnprintable(t *testing.T) {
	cases := []struct {
		in      string
		lenLost int
		want    string
	}{
		{"Hello, 世界", 0, "Hello, 世界"},
		{"m", 1, "m"},
		{"m", 0, "m"},
		{" ", 3, " "},
		{"a‎b‏c", 6, "abc"}, // only 2 removed but count as 3 each
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := RemoveUnprintable(tt.in)
			charsRemoved := len(tt.in) - len(out)
			if tt.lenLost != charsRemoved {
				t.Errorf("\ncharsRemoved:  %#v\nwant: %#v\n", charsRemoved, tt.lenLost)
			}
			if out != tt.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestGetLine(t *testing.T) {
	cases := []struct {
		in   string
		line int
		want string
	}{
		{"Hello", 1, "Hello"},
		{"Hello", 2, ""},
		{"Hello\nworld", 1, "Hello"},
		{"Hello\nworld", 2, "world"},
		{"Hello\nworld", 3, ""},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			out := GetLine(tt.in, tt.line)
			if out != tt.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func BenchmarkLeft(b *testing.B) {
	text := strings.Repeat("Hello, world, it's a sentences!\n", 200)
	for n := 0; n < b.N; n++ {
		Left(text, 250)
	}
}

func BenchmarkRemoveUnprintable(b *testing.B) {
	text := strings.Repeat("Hello, world, it's a sentences!\n", 20000)
	for n := 0; n < b.N; n++ {
		GetLine(text, 200)
	}
}

func TestUniq(t *testing.T) {
	tests := []struct {
		in   []string
		want []string
	}{
		{
			[]string{"a", "b", "c"},
			[]string{"a", "b", "c"},
		},
		{
			[]string{"a", "b", "c", "a", "b", "n", "a", "aaa", "n", "x"},
			[]string{"a", "aaa", "b", "c", "n", "x"},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := Uniq(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nwant: %q\ngot:  %q", tt.want, got)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		list []string
		find string
		want bool
	}{
		{[]string{"hello"}, "hello", true},
		{[]string{"hello"}, "hell", false},
		{[]string{"hello", "world", "test"}, "world", true},
		{[]string{"hello", "world", "test"}, "", false},
		{[]string{}, "", false},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := Contains(tt.list, tt.find)
			if got != tt.want {
				t.Errorf("want: %#v\ngot:  %#v", tt.want, got)
			}
		})
	}
}

func TestChoose(t *testing.T) {
	tests := []struct {
		in   []string
		want string
	}{
		{nil, ""},
		{[]string{}, ""},
		{[]string{"a"}, "a"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := Choose(tt.in)
			if out != tt.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	cases := []struct {
		fun  func(string) bool
		in   []string
		want []string
	}{
		{
			FilterEmpty,
			[]string(nil),
			[]string(nil),
		},
		{
			FilterEmpty,
			[]string{},
			[]string(nil),
		},
		{
			FilterEmpty,
			[]string{"1"},
			[]string{"1"},
		},
		{
			FilterEmpty,
			[]string{"", "1", ""},
			[]string{"1"},
		},
		{
			FilterEmpty,
			[]string{"", "1", "", "2", "asd", "", "", "", "zx", "", "a"},
			[]string{"1", "2", "asd", "zx", "a"},
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

func TestDifference(t *testing.T) {
	tests := []struct {
		inSet    []string
		inOthers [][]string
		want     []string
	}{
		{[]string{}, [][]string{}, []string{}},
		{nil, [][]string{}, []string{}},
		{[]string{}, nil, []string{}},
		{nil, nil, []string{}},
		{[]string{"1"}, [][]string{{"1"}}, []string{}},
		{[]string{"1", "2", "2", "3"}, [][]string{{"1", "2", "2", "3"}}, []string{}},
		{[]string{"1", "2", "2", "3"}, [][]string{{"1", "2"}, {"3"}}, []string{}},
		{[]string{"1", "2"}, [][]string{{"1"}}, []string{"2"}},
		{[]string{"1", "2", "3"}, [][]string{{"1"}}, []string{"2", "3"}},
		{[]string{"1", "2", "3"}, [][]string{{}, {"1"}}, []string{"2", "3"}},
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
