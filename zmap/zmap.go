// Package zmap implements generic functions for maps.
package zmap

import (
	"cmp"
	"slices"
)

// KeysOrdered returns the sorted keys of the map.
func KeysOrdered[M ~map[K]V, K cmp.Ordered, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	slices.Sort(r)
	return r
}

// LongestKey returns the longest key in this map and the unsorted list of all
// keys.
func LongestKey[M ~map[string]V, V any](m M) ([]string, int) {
	var (
		l = 0
		r = make([]string, 0, len(m))
	)
	for k := range m {
		r = append(r, k)
		if ll := len(k); ll > l {
			l = ll
		}
	}
	return r, l
}
