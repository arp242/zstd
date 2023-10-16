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
