package zmap

import (
	"cmp"
	"maps"
	"slices"
)

// Keys returns an unsorted list of keys of the map.
//
// Deprecated: use slices.Collect(maps.Keys())
//
//go:fix inline
func Keys[M ~map[K]V, K cmp.Ordered, V any](m M) []K { return slices.Collect(maps.Keys(m)) }

// Values returns the values of the map.
//
// Deprecated: use slices.Collect(maps.Values())
//
//go:fix inline
func Values[M ~map[K]V, K comparable, V any](m M) []V { return slices.Collect(maps.Values(m)) }
