package zstring

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestIndent(t *testing.T) {
	tests := []struct {
		n    int
		in   string
		want string
	}{
		{2, "Hello", "  Hello"},
		{2, "Hello\n", "  Hello\n"},
		{2, "Hello\nWorld", "  Hello\n  World"},
		{2, "Hello\nWorld\n", "  Hello\n  World\n"},
		{2, "Hello\nWorld\n\n", "  Hello\n  World\n  \n"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := Indent(tt.in, tt.n)
			if have != tt.want {
				t.Errorf("\nhave:\n%q\n\nwant:\n%q", have, tt.want)
			}
		})
	}
}

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
func BenchmarkReverse(b *testing.B) {
	s := strings.Repeat("H€łø🖖", 20)

	var c string
	for n := 0; n < b.N; n++ {
		c = Reverse(s)
	}
	_ = c
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

func TestElide(t *testing.T) {
	cases := []struct {
		in         string
		n          int
		wantLeft   string
		wantRight  string
		wantCenter string
	}{
		{"Hello", 100, "Hello", "Hello", "Hello"},
		{"Hello", 1, "H…", "…o", "H…"},
		{"Hello", 2, "He…", "…lo", "H…o"},
		{"Hello", 5, "Hello", "Hello", "Hello"},
		{"Hello", 4, "Hell…", "…ello", "He…lo"},
		{"Hello", 0, "…", "…", "…"},
		{"汉语漢語", 1, "汉…", "…語", "汉…"},
		{"汉语漢語", 2, "汉语…", "…漢語", "汉…語"},
		{"汉语漢語", 3, "汉语漢…", "…语漢語", "汉语…語"},
		{"汉语漢語", 4, "汉语漢語", "汉语漢語", "汉语漢語"},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			left := ElideLeft(tt.in, tt.n)
			if left != tt.wantLeft {
				t.Errorf("\nout:  %#v\nwant: %#v\n", left, tt.wantLeft)
			}

			right := ElideRight(tt.in, tt.n)
			if right != tt.wantRight {
				t.Errorf("\nout:  %#v\nwant: %#v\n", right, tt.wantRight)
			}

			center := ElideCenter(tt.in, tt.n)
			if center != tt.wantCenter {
				t.Errorf("\nout:  %#v\nwant: %#v\n", center, tt.wantCenter)
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
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Sub(text, 50, 250)
	}
}

func BenchmarkElideLeft(b *testing.B) {
	text := strings.Repeat("Hello, world, it's a sentences!\n", 200)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ElideLeft(text, 250)
	}
}

func BenchmarkRemoveUnprintable(b *testing.B) {
	text := strings.Repeat("Hello, world, it's a sentences!\n", 20000)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		GetLine(text, 200)
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

		// TODO
		// {
		// 	"2022-10-13 – http://goatcounter.goatcounter.localhost:8081",
		// 	77,
		// 	"2022-10-13 – http://goatcounter.goatcounter.localhost:8081                   ",
		// 	"                   2022-10-13 – http://goatcounter.goatcounter.localhost:8081",
		// 	"          2022-10-13 – http://goatcounter.goatcounter.localhost:8081         ",
		// },

		// {
		// 	"2022-10-13 – http://goatcounter.goatcounter.localhost:8081",
		// 	78,
		// 	"2022-10-13 – http://goatcounter.goatcounter.localhost:8081                    ",
		// 	"                    2022-10-13 – http://goatcounter.goatcounter.localhost:8081",
		// 	"          2022-10-13 – http://goatcounter.goatcounter.localhost:8081          ",
		// },
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%d", tt.in, tt.n), func(t *testing.T) {
			left, right, center := AlignLeft(tt.in, tt.n), AlignRight(tt.in, tt.n), AlignCenter(tt.in, tt.n)
			if left != tt.left {
				t.Errorf("left wrong\nhave: %q\nwant: %q", left, tt.left)
			}
			if right != tt.right {
				t.Errorf("right wrong\nhave: %q\nwant: %q", right, tt.right)
			}
			if center != tt.center {
				t.Errorf("center wrong\nhave: %q\nwant: %q", center, tt.center)
			}
		})
	}
}

func TestIndexPairs(t *testing.T) {
	tests := []struct {
		in, start, end string
		want           [][]int
	}{
		{"", "{", "}", nil},
		{"xx", "{", "}", nil},
		{"{xx}", "{", "}", [][]int{{0, 3}}},
		{"{xx} {yyy}", "{", "}", [][]int{{5, 9}, {0, 3}}},

		// No end mark
		{" {xx} {yyy", "{", "}", [][]int{{1, 4}}},
		{" {xx} {yyy {x}", "{", "}", [][]int{{6, 13}, {1, 4}}},

		// Nested
		{"{xx {yy} }", "{", "}", [][]int{{0, 7}}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := IndexPairs(tt.in, tt.start, tt.end)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:  %#v\nwant: %#v", got, tt.want)
			}
		})
	}
}
func BenchmarkIndexPairs(b *testing.B) {
	text := "Hello {world}, {asc}\n"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		IndexPairs(text, "{", "}")
	}
}

func TestReplacePairs(t *testing.T) {
	tests := []struct {
		in, start, end string
		f              func(int, string) string
		want           string
	}{
		{"", "{", "}", func(int, string) string { return "X" }, ""},
		{"xx", "{", "}", func(int, string) string { return "X" }, "xx"},

		{"{xx}", "{", "}", func(int, string) string { return "X" }, "X"},
		{"{xx} {yyy}", "{", "}", func(i int, s string) string { return strings.Repeat("X", i+1) }, "XX X"},

		// No end mark
		{"{xx} → {yyy", "{", "}", func(i int, s string) string { return strings.Repeat("X", i+1) }, "X → {yyy"},
		{"{xx} → {yyy {x}", "{", "}", func(i int, s string) string { return strings.Repeat("X", i+1) }, "XX → X"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			have := ReplacePairs(tt.in, tt.start, tt.end, tt.f)
			if have != tt.want {
				t.Errorf("\nhave: %q\nwant: %q", have, tt.want)
			}
		})
	}
}

func TestUpto(t *testing.T) {
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

func TestIndexAll(t *testing.T) {
	tests := []struct {
		in, find string
		want     []int
	}{
		{"", "", nil},
		{"", ".", nil},
		{".", "", nil},
		{"a.b", ".", []int{1}},
		{"a.b.", ".", []int{1, 3}},
		{".a.b.", ".", []int{0, 2, 4}},
		{".a.b.", "b.", []int{3}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := IndexAll(tt.in, tt.find)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:  %#v\nwant: %#v", got, tt.want)
			}
		})
	}
}

func TestIndexN(t *testing.T) {
	tests := []struct {
		in, find string
		n        uint
		want     int
	}{
		{"", "", 1, -1},
		{"", ".", 1, -1},
		{".", "", 1, -1},

		{"a.b.c.d", ".", 0, 1},
		{"a.b.c.d", ".", 1, 1},
		{"a.b.c.d", ".", 2, 3},
		{"a.b.c.d", ".", 3, 5},
		{"a.b.c.d", ".", 4, -1},

		{"aa ... bb ... cc ... dd", "..", 0, 3},
		{"aa ... bb ... cc ... dd", "..", 1, 3},
		{"aa ... bb ... cc ... dd", "..", 2, 10},
		{"aa ... bb ... cc ... dd", "..", 3, 17},
		{"aa ... bb ... cc ... dd", "..", 4, -1},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := IndexN(tt.in, tt.find, tt.n)
			if have != tt.want {
				t.Errorf("\nhave: %#v\nwant: %#v", have, tt.want)
			}
		})
	}
}

func TestIsASCII(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"asd", true},
		{"asó", false},
		{"as€", false},

		{"\x00", true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := IsASCII(tt.in)
			if have != tt.want {
				t.Errorf("\nhave: %t\nwant: %t", have, tt.want)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", ""},
		{"Hello", "Hello"},
		{"Hello\nworld", "Hello world"},
		{"Hello\nworld\n!", "Hello world !"},

		{"Hello\nworld!\n\nXXX", "Hello world!\n\nXXX"},

		{"\nHello\nworld!\n\nXXX\n", "Hello world!\n\nXXX"},
		{"\n\n\nHello\nworld!\n\nXXX\n", "Hello world!\n\nXXX"},
		{"\n\n\nHello\n\n\n\n\n\nworld!\n\n\n\n\n\n\nXXX\n\n\n", "Hello\n\nworld!\n\nXXX"},

		{"\tHello\n\n\tworld!", "\tHello\n\n\tworld!"},

		{" Hello\n world!", " Hello  world!"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := Unwrap(tt.in)
			if have != tt.want {
				t.Errorf("\nhave: %q\nwant: %q", have, tt.want)
			}
		})
	}
}

func BenchmarkUnwrap(b *testing.B) {
	var (
		l = strings.Repeat("Hello, world, test. ", 4) + "\n"
		s string
	)
	for i := 0; i < 10; i++ {
		if i > 0 {
			s += "\n"
		}
		s += strings.Repeat(l, 3)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Unwrap(s)
	}
}
