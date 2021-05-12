package zstring

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestWordWrap(t *testing.T) {
	tests := []struct {
		n      int
		prefix string
		in     string
		want   string
	}{
		{10, "",
			"Hello",
			"Hello"},
		{3, "", "Hello",
			"Hello"},

		{10, "",
			"Hello, world!",
			"Hello,\nworld!"},
		{10, "",
			"Hello, world! it's a test",
			"Hello,\nworld!\nit's a\ntest"},
		{30, "",
			"Click this link yo: https://github.com/zgoat/zstd/blob/master/README.md",
			"Click this link yo:\nhttps://github.com/zgoat/zstd/blob/master/README.md"},

		{10, "> ",
			"> Hello, world! it's a test",
			"> Hello,\n> world!\n> it's a\n> test"},
		{10, "> ",
			"> H€łłø, ωørłð¿ it’s a ţësţ",
			"> H€łłø,\n> ωørłð¿\n> it’s a\n> ţësţ"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := WordWrap(tt.in, tt.prefix, tt.n)
			if got != tt.want {
				t.Errorf("\ngot:\n%s\n\nwant:\n%s", got, tt.want)
			}
		})
	}
}

func TestTabWidth(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"a", 1},

		// Tabs.
		{"\t", 8},
		{"\ta", 9},
		{"a\t", 8},
		{"aaaa\tx", 9},
		{"aaaaaaa\tx", 9},

		{"\t\t", 16},
		{"a\ta\t", 16},
		{"a\ta\ta", 17},

		// Emojis.
		// {"🧑\u200d🚒", 1},
		// {"🧑🏽\u200d🚒", 1},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := TabWidth(tt.in)
			if got != tt.want {
				t.Errorf("\ngot:  %d\nwant: %d", got, tt.want)
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

func TestContainsAny(t *testing.T) {
	tests := []struct {
		list []string
		find []string
		want bool
	}{
		{[]string{"hello"}, []string{"hello"}, true},
		{[]string{"hello"}, []string{"hell"}, false},
		{[]string{"hello", "world", "test"}, []string{"world"}, true},
		{[]string{"hello", "world", "test"}, []string{""}, false},
		{[]string{"hello", "world", "test"}, []string{"asd", "asd", "test"}, true},
		{[]string{}, []string{""}, false},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			got := ContainsAny(tt.list, tt.find...)
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

func TestRemove(t *testing.T) {
	tests := []struct {
		in        []string
		rm        string
		want      bool
		wantSlice []string
	}{
		{nil, "xx", false, nil},
		{[]string{}, "xx", false, []string{}},

		{[]string{"xx"}, "xx", true, []string{}},
		{[]string{"xx", "a"}, "xx", true, []string{"a"}},
		{[]string{"a", "xx"}, "xx", true, []string{"a"}},
		{[]string{"xx", "a", "xx"}, "xx", true, []string{"a"}},
		{[]string{"xx", "a", "xx", "b"}, "xx", true, []string{"a", "b"}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := Remove(&tt.in, tt.rm)
			if have != tt.want {
				t.Errorf("\nhave: %t\nwant: %t", have, tt.want)
			}
			if !reflect.DeepEqual(tt.in, tt.wantSlice) {
				fmt.Println(len(tt.in), cap(tt.in))
				t.Errorf("\nhave: %v\nwant: %v", tt.in, tt.wantSlice)
			}
		})
	}
}
