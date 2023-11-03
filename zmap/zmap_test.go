package zmap

import (
	"reflect"
	"testing"
)

func TestKeysOrdered(t *testing.T) {
	var nm map[int]int
	tests := []struct {
		in   map[int]int
		want []int
	}{
		{nm, []int{}},
		{map[int]int{1: 0, 2: 0, 3: 0}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := KeysOrdered(tt.in)
			if !reflect.DeepEqual(have, tt.want) {
				t.Errorf("\nhave: %#v\nwant: %#v", have, tt.want)
			}
		})
	}
}

func TestLongestKey(t *testing.T) {
	tests := []struct {
		in   map[string]int
		want int
	}{
		{nil, 0},
		{map[string]int{"": 0}, 0},
		{map[string]int{"a": 0}, 1},
		{map[string]int{"aa": 0}, 2},
		{map[string]int{"a": 5, "aa": 3}, 2},
		{map[string]int{"aa": 0}, 2},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := LongestKey(tt.in)
			if have != tt.want {
				t.Errorf("\nhave: %v\nwant: %v", have, tt.want)
			}
		})
	}
}
