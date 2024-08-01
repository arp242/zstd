// Package zmap implements generic functions for maps.
package zmap

import "sort"

// cmp.Ordered, added in Go 1.21.
type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// KeysOrdered returns the sorted keys of the map.
func KeysOrdered[M ~map[K]V, K ordered, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	sort.Slice(r, func(i, j int) bool { return r[i] < r[j] })
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

// Values returns the values of the map.
func Values[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}
