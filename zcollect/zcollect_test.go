package zcollect

import (
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
