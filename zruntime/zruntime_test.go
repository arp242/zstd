package zruntime

import (
	"fmt"
	"testing"
)

func TestGoroutineID(t *testing.T) {
	id := GoroutineID()
	if id <= 0 {
		t.Errorf("lower than 0: %d", id)
	}
}

func TestSizeOf(t *testing.T) {
	n := int64(123)
	tests := []struct {
		in   any
		want int64
	}{
		{"", 16},
		{"abcd", 20},
		{int32(1), 32 / 8},
		{float64(1), 64 / 8},
		{true, 1},

		{[]byte{}, 24},
		{new([]byte), 32}, // 24 for byte, 8 for ptr
		{[]string{}, 24},
		{[]string{""}, 40},
		{[]string{"aa"}, 42},
		{[]string{"aa"}[:0], 48}, // TODO: is this correct?
		{map[int8]int8{}, 32},
		{map[int8]int8{1: 2}, 32},
		{map[int8]int8{1: 2, 3: 4, 5: 6, 7: 8, 9: 10, 11: 12, 13: 14, 15: 16}, 64},
		{struct{}{}, 0},
		{struct {
			x string
			y int64
			z []string
		}{"h", 1, []string{"aa"}}, 90},

		{struct {
			x string
			y *int64
			z []string
		}{"h", &n, []string{"aa"}}, 98},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			got := SizeOf(tt.in)
			if got != tt.want {
				t.Errorf("\ngot:  %d\nwant: %d", got, tt.want)
			}
		})
	}
}
func TestSizeOfCycles(t *testing.T) {
	type V struct {
		Z int
		E *V
	}

	v := &V{Z: 25}
	want := SizeOf(v)
	v.E = v // induce a cycle
	got := SizeOf(v)
	if got != want {
		t.Errorf("Cyclic size: got %d, want %d", got, want)
	}
}

func TestCallers(t *testing.T) {
	func() {
		for _, c := range Callers() {
			fmt.Println(c)
		}
	}()
}
