package zslice

import (
	"fmt"
	"reflect"
	"testing"
)

func TestChoose(t *testing.T) {
	n := Choose([]int{})
	if n != 0 {
		t.Fatal("not 0")
	}

	n = Choose([]int{1})
	if n != 1 {
		t.Fatal("not 1")
	}

	var (
		l   = []int{1, 2, 3}
		rnd = make([]int, 0, 100)
	)
	for i := 0; i < 100; i++ {
		rnd = append(rnd, Choose(l))
	}

	var one, two, three []int
	for _, r := range rnd {
		switch r {
		case 1:
			one = append(one, r)
		case 2:
			two = append(two, r)
		case 3:
			three = append(three, r)
		}
	}

	if len(one) < 10 {
		t.Error("one", one)
	}
	if len(two) < 10 {
		t.Error("two", two)
	}
	if len(three) < 10 {
		t.Error("three", three)
	}
}

func TestShuffle(t *testing.T) {
	{
		var l []int
		Shuffle(l)
		if l != nil {
			t.Fatal()
		}
	}

	{
		l := []int{1}
		Shuffle(l)
		if !reflect.DeepEqual(l, []int{1}) {
			t.Fatal()
		}
	}

	{
		var (
			rnd = make([]int, 0, 100)
		)
		for i := 0; i < 100; i++ {
			l := []int{1, 2, 3}
			Shuffle(l)
			rnd = append(rnd, l[0])
		}

		var one, two, three []int
		for _, r := range rnd {
			switch r {
			case 1:
				one = append(one, r)
			case 2:
				two = append(two, r)
			case 3:
				three = append(three, r)
			}
		}

		if len(one) < 10 {
			t.Error("one", one)
		}
		if len(two) < 10 {
			t.Error("two", two)
		}
		if len(three) < 10 {
			t.Error("three", three)
		}
	}
}

func TestContainsAny(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			list []string
			find []string
			want bool
		}{
			{[]string{}, []string{""}, false},
			{[]string{"hello"}, []string{"hello"}, true},
			{[]string{"hello"}, []string{"hell"}, false},
			{[]string{"hello", "world", "test"}, []string{"world"}, true},
			{[]string{"hello", "world", "test"}, []string{""}, false},
			{[]string{"hello", "world", "test"}, []string{"asd", "asd", "test"}, true},
		}
		for i, tt := range tests {
			t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
				have := ContainsAny(tt.list, tt.find...)
				if have != tt.want {
					t.Errorf("want: %#v\nhave: %#v", tt.want, have)
				}
			})
		}
	})

	t.Run("float64", func(t *testing.T) {
		tests := []struct {
			list []float64
			find []float64
			want bool
		}{
			{[]float64{}, []float64{0}, false},
			{[]float64{1.123}, []float64{1, 1.123}, true},
		}
		for i, tt := range tests {
			t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
				have := ContainsAny(tt.list, tt.find...)
				if have != tt.want {
					t.Errorf("want: %#v\nhave: %#v", tt.want, have)
				}
			})
		}
	})
}

func TestUniqSort(t *testing.T) {
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
			have := UniqSort(tt.in)
			if !reflect.DeepEqual(have, tt.want) {
				t.Errorf("\nwant: %q\nhave: %q", tt.want, have)
			}
		})
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
			[]string{"a", "b", "c", "n", "aaa", "x"},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			have := Uniq(tt.in)
			if !reflect.DeepEqual(have, tt.want) {
				t.Errorf("\nwant: %q\nhave: %q", tt.want, have)
			}
		})
	}
}

func TestIsUniq(t *testing.T) {
	tests := []struct {
		in   []string
		want bool
	}{
		{[]string{}, true},
		{[]string{"a", "b", "c"}, true},
		{[]string{"a", "b", "c", "a", "b", "n", "a", "aaa", "n", "x"}, false},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			have := IsUniq(tt.in)
			if have != tt.want {
				t.Errorf("\nwant: %v\nhave: %v", tt.want, have)
			}
		})
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		n    int
		want []string
	}{
		{0, []string{}},
		{1, []string{"X"}},
		{3, []string{"X", "X", "X"}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			have := Repeat("X", tt.n)
			if !reflect.DeepEqual(have, tt.want) {
				t.Errorf("\nwant: %q\nhave: %q", tt.want, have)
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

func TestMax(t *testing.T) {
	tests := []struct {
		in   []int
		want int
	}{
		{[]int{0}, 0},
		{[]int{0, 5, 6}, 6},
		{[]int{0, 6, 5}, 6},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := Max(tt.in)
			if have != tt.want {
				t.Errorf("\nhave: %q\nwant: %q", have, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		in   []int
		want int
	}{
		{[]int{0}, 0},
		{[]int{0, 5, 6}, 0},
		{[]int{0, 6, 5}, 0},
		{[]int{0, 6, 5, -5}, -5},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := Min(tt.in)
			if have != tt.want {
				t.Errorf("\nhave: %q\nwant: %q", have, tt.want)
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

func TestIntersect(t *testing.T) {
	tests := []struct {
		inA  []string
		inB  []string
		want []string
	}{
		{[]string{}, []string{}, []string{}},
		{[]string{"X"}, []string{"X"}, []string{"X"}},
		{[]string{"X", "a"}, []string{"X"}, []string{"X"}},
		{[]string{"X", "a"}, []string{"X", "b"}, []string{"X"}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := Intersect(tt.inA, tt.inB)
			if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestSameElements(t *testing.T) {
	tests := []struct {
		a    []string
		b    []string
		want bool
	}{
		{[]string{}, []string{}, true},
		{[]string{"a", "b"}, []string{"b", "a"}, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := SameElements(tt.a, tt.b)
			if have != tt.want {
				t.Errorf("\nhave: %v\nwant: %v", have, tt.want)
			}
		})
	}
}

func TestRemoveIndexes(t *testing.T) {
	tests := []struct {
		in   []string
		rm   []int
		want []string
	}{
		{[]string{}, []int{}, []string{}},
		{[]string{"a"}, []int{0}, []string{}},
		{[]string{"a", "b", "c"}, []int{1}, []string{"a", "c"}},
		{[]string{"a", "b", "c"}, []int{1, 2}, []string{"a"}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			RemoveIndexes(&tt.in, tt.rm...)
			if !reflect.DeepEqual(tt.in, tt.want) {
				t.Errorf("\nhave: %#v\nwant: %#v", tt.in, tt.want)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	tests := []struct {
		in           []string
		inLen, inCap int
		want         string
	}{
		{[]string{}, 0, 0, `0 0 []`},
		{[]string{}, 0, 1, `0 1 []`},
		{[]string{}, 1, 1, `1 1 []`},
		{[]string{"x"}, 1, 1, `1 1 [x]`},

		{[]string{"x"}, 2, 2, `2 2 [x ]`},
		{[]string{"x"}, 2, 8, `2 8 [x ]`},

		{[]string{"a", "b", "c"}, 1, 1, `1 1 [a]`},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			h := Copy(tt.in, tt.inLen, tt.inCap)
			have := fmt.Sprintf("%d %d %s", len(h), cap(h), h)
			if have != tt.want {
				t.Errorf("\nhave: %s\nwant: %s", have, tt.want)
			}
		})
	}
}

func TestAppendCopy(t *testing.T) {
	tests := []struct {
		in   []string
		app  string
		more []string
		want []string
	}{
		{nil, "X", nil, []string{"X"}},
		{[]string{}, "X", nil, []string{"X"}},
		{[]string{}, "X", []string{"Y"}, []string{"X", "Y"}},
		{[]string{}, "X", []string{"Y", "Z"}, []string{"X", "Y", "Z"}},
		{[]string{"a", "b"}, "X", nil, []string{"a", "b", "X"}},
		{[]string{"a", "b"}, "X", []string{"Y"}, []string{"a", "b", "X", "Y"}},
		{[]string{"a", "b"}, "X", []string{"Y", "Z"}, []string{"a", "b", "X", "Y", "Z"}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			before := fmt.Sprintf("%v", tt.in)
			have := AppendCopy(tt.in, tt.app, tt.more...)
			if !reflect.DeepEqual(have, tt.want) {
				t.Errorf("\nhave: %s\nwant: %s", have, tt.want)
			}
			pin, phave := fmt.Sprintf("%p", tt.in), fmt.Sprintf("%p", have)
			if pin == phave {
				t.Errorf("same array; wasn't copied\nhave: %s\nwant: %s", pin, phave)
			}
			if a := fmt.Sprintf("%v", tt.in); a != before {
				t.Errorf("tt.in changed\nbefore: %s\nafter:  %s", before, a)
			}
		})
	}
}
