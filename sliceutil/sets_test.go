package sliceutil

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDifferenceInt(t *testing.T) {
	tests := []struct {
		inSet    []int64
		inOthers [][]int64
		want     []int64
	}{
		{[]int64{}, [][]int64{}, []int64{}},
		{nil, [][]int64{}, []int64{}},
		{[]int64{}, nil, []int64{}},
		{nil, nil, []int64{}},
		{[]int64{1}, [][]int64{{1}}, []int64{}},
		{[]int64{1, 2, 2, 3}, [][]int64{{1, 2, 2, 3}}, []int64{}},
		{[]int64{1, 2, 2, 3}, [][]int64{{1, 2}, {3}}, []int64{}},
		{[]int64{1, 2}, [][]int64{{1}}, []int64{2}},
		{[]int64{1, 2, 3}, [][]int64{{1}}, []int64{2, 3}},
		{[]int64{1, 2, 3}, [][]int64{{}, {1}}, []int64{2, 3}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := DifferenceInt(tt.inSet, tt.inOthers...)
			if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestDifferenceString(t *testing.T) {
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
			out := DifferenceString(tt.inSet, tt.inOthers...)
			if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestComplementInt(t *testing.T) {
	tests := []struct {
		name         string
		inA, inB     []int64
		wantA, wantB []int64
	}{
		{
			name: "EmptyLists",
		},
		{
			name:  "AOnly",
			inA:   []int64{1, 2, 3},
			wantA: []int64{1, 2, 3},
		},
		{
			name:  "BOnly",
			inB:   []int64{1, 2, 3},
			wantB: []int64{1, 2, 3},
		},
		{
			name: "Equal",
			inA:  []int64{1, 2, 3},
			inB:  []int64{1, 2, 3},
		},
		{
			name:  "Disjoint",
			inA:   []int64{1, 2, 3},
			inB:   []int64{5, 6, 7},
			wantA: []int64{1, 2, 3},
			wantB: []int64{5, 6, 7},
		},
		{
			name:  "Overlap",
			inA:   []int64{1, 2, 3, 4},
			inB:   []int64{3, 4, 5, 6},
			wantA: []int64{1, 2},
			wantB: []int64{5, 6},
		},
		{
			name:  "Overlap with repeated values",
			inA:   []int64{6, 4, 5, 3, 6},
			inB:   []int64{2, 1, 4, 3, 1},
			wantA: []int64{6, 5, 6},
			wantB: []int64{2, 1, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aOnly, bOnly := ComplementInt(tt.inA, tt.inB)
			if !reflect.DeepEqual(aOnly, tt.wantA) {
				t.Errorf("aOnly wrong\ngot:  %#v\nwant: %#v\n", aOnly, tt.wantA)
			}
			if !reflect.DeepEqual(bOnly, tt.wantB) {
				t.Errorf("bOnly wrong\ngot:  %#v\nwant: %#v\n", bOnly, tt.wantB)
			}
		})
	}
}

func BenchmarkComplement_equal(b *testing.B) {
	listA := []int64{1, 2, 3}
	listB := []int64{1, 2, 3}

	for n := 0; n < b.N; n++ {
		ComplementInt(listA, listB)
	}
}

func BenchmarkComplement_disjoint(b *testing.B) {
	listA := []int64{1, 2, 3}
	listB := []int64{5, 6, 7}

	for n := 0; n < b.N; n++ {
		ComplementInt(listA, listB)
	}
}

func BenchmarkComplement_overlap(b *testing.B) {
	listA := []int64{1, 2, 3, 4}
	listB := []int64{3, 4, 5, 6}

	for n := 0; n < b.N; n++ {
		ComplementInt(listA, listB)
	}
}
