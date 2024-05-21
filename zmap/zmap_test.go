package zmap

import (
	"reflect"
	"slices"
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
		in      map[string]int
		wantK   []string
		wantLen int
	}{
		{nil, []string{}, 0},
		{map[string]int{"": 0}, []string{""}, 0},
		{map[string]int{"a": 0}, []string{"a"}, 1},
		{map[string]int{"aa": 0}, []string{"aa"}, 2},
		{map[string]int{"a": 5, "aa": 3}, []string{"a", "aa"}, 2},
		{map[string]int{"aa": 0}, []string{"aa"}, 2},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			haveK, haveLen := LongestKey(tt.in)
			if haveLen != tt.wantLen {
				t.Errorf("\nhave: %v\nwant: %v", haveLen, tt.wantLen)
			}
			slices.Sort(haveK)
			if !reflect.DeepEqual(haveK, tt.wantK) {
				t.Errorf("\nhave: %#v\nwant: %#v", haveK, tt.wantK)
			}
		})
	}
}
