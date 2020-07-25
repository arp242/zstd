package zstring

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"zgo.at/zstd/ztest"
)

func TestReverse(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", ""},
		{"ab", "ba"},
		{"Hello, world", "dlrow ,olleH"},
		{"H€łø🖖", "🖖øł€H"},

		// This is broken, as combining marks. That's probably okay.
		//{"🤦‍♂️", "🤦‍♂️"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := Reverse(tt.in)
			if got != tt.want {
				t.Errorf("\ngot:  %q\nwant: %q", got, tt.want)
			}
		})
	}
}

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

func TestSub(t *testing.T) {
	cases := []struct {
		in   string
		i, j int
		want string
	}{
		{"", 0, 0, ""},
		{"Hello", 0, 1, "H"},
		{"Hello", 1, 5, "ello"},
		{"汉语漢語", 0, 1, "汉"},
		{"汉语漢語", 0, 3, "汉语漢"},
		{"汉语漢語", 2, 4, "漢語"},
		{"汉语漢語", 0, 4, "汉语漢語"},

		// Length longer than string.
		{"He", 0, 100, "He"},
		{"汉语漢語", 0, 100, "汉语漢語"},

		// Start longer than string.
		{"He", 100, 100, ""},
		{"汉语漢語", 100, 100, ""},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := Sub(tt.in, tt.i, tt.j)
			if out != tt.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestLeft(t *testing.T) {
	ztest.MustInline(t, "zgo.at/zstd/zstring Left")

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

func BenchmarkSub(b *testing.B) {
	text := strings.Repeat("Hello, world, it's a sentences!\n", 200)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Sub(text, 50, 250)
	}
}

func BenchmarkLeft(b *testing.B) {
	text := strings.Repeat("Hello, world, it's a sentences!\n", 200)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Left(text, 250)
	}
}

func BenchmarkRemoveUnprintable(b *testing.B) {
	text := strings.Repeat("Hello, world, it's a sentences!\n", 20000)
	b.ReportAllocs()
	b.ResetTimer()
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
	ztest.MustInline(t, "zgo.at/zstd/zstring FilterEmpty")

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

func TestAlign(t *testing.T) {
	tests := []struct {
		in                  string
		n                   int
		left, right, center string
	}{
		{"", 4, "    ", "    ", "    "},
		{"a", 4, "a   ", "   a", " a  "},

		{"Hello", 4, "Hello", "Hello", "Hello"},
		{"Hello", -2, "Hello", "Hello", "Hello"},

		{"Hello", 6, "Hello ", " Hello", "Hello "},
		{"Hello", 7, "Hello  ", "  Hello", " Hello "},
		{"Hello", 8, "Hello   ", "   Hello", " Hello  "},
		{"Hello", 9, "Hello    ", "    Hello", "  Hello  "},
		{"Hello", 10, "Hello     ", "     Hello", "  Hello   "},
		{"Hello", 11, "Hello      ", "      Hello", "   Hello   "},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%d", tt.in, tt.n), func(t *testing.T) {
			left := AlignLeft(tt.in, tt.n)
			right := AlignRight(tt.in, tt.n)
			center := AlignCenter(tt.in, tt.n)

			if left != tt.left {
				t.Errorf("left wrong\ngot:  %q\nwant: %q", left, tt.left)
			}
			if right != tt.right {
				t.Errorf("right wrong\ngot:  %q\nwant: %q", right, tt.right)
			}
			if center != tt.center {
				t.Errorf("center wrong\ngot:  %q\nwant: %q", center, tt.center)
			}
		})
	}
}

func BenchmarkReverse(b *testing.B) {
	s := strings.Repeat("H€łø🖖", 20)
	b.ReportAllocs()

	var c string
	for n := 0; n < b.N; n++ {
		c = Reverse(s)
	}
	_ = c
}

func TestUpto(t *testing.T) {
	ztest.MustInline(t, "zgo.at/zstd/zstring", "Upto", "From")

	tests := []struct {
		in, sep, upto, from string
	}{
		{"", "", "", ""},
		{"ab", ":", "ab", "ab"},
		{"a:b", ":", "a", "b"},
		{"a:b", "::", "a:b", "a:b"},
		{"a::b", "::", "a", "b"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			a := Upto(tt.in, tt.sep)
			b := From(tt.in, tt.sep)

			if a != tt.upto {
				t.Errorf("upto\ngot:  %q\nwant: %q", a, tt.upto)
			}
			if b != tt.from {
				t.Errorf("from\ngot:  %q\nwant: %q", b, tt.from)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		in, sep string
		want2   []string
		want3   []string
		want4   []string
	}{
		{"a", ":",
			[]string{"a", ""},
			[]string{"a", "", ""},
			[]string{"a", "", "", ""},
		},
		{"a:b", ":",
			[]string{"a", "b"},
			[]string{"a", "b", ""},
			[]string{"a", "b", "", ""},
		},

		{"a:b:c:d:e", ":",
			[]string{"a", "b:c:d:e"},
			[]string{"a", "b", "c:d:e"},
			[]string{"a", "b", "c", "d:e"},
		},

		{"a:::b:c:d:e", ":",
			[]string{"a", "::b:c:d:e"},
			[]string{"a", "", ":b:c:d:e"},
			[]string{"a", "", "", "b:c:d:e"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got1, got2 := Split2(tt.in, tt.sep)
			got := []string{got1, got2}
			if !reflect.DeepEqual(got, tt.want2) {
				t.Errorf("\ngot:  %q\nwant: %q", got, tt.want2)
			}

			got1, got2, got3 := Split3(tt.in, tt.sep)
			got = []string{got1, got2, got3}
			if !reflect.DeepEqual(got, tt.want3) {
				t.Errorf("\ngot:  %q\nwant: %q", got, tt.want3)
			}

			got1, got2, got3, got4 := Split4(tt.in, tt.sep)
			got = []string{got1, got2, got3, got4}
			if !reflect.DeepEqual(got, tt.want4) {
				t.Errorf("\ngot:  %q\nwant: %q", got, tt.want4)
			}
		})
	}
}
