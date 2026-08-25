package zslice

import (
	"cmp"
	"slices"
)

// Repeat returns a slice with the value v repeated n times.
//
// Deprecated: use slices.Repeat()
//
//go:fix inline
func Repeat[T any](s T, n int) []T { return slices.Repeat([]T{s}, n) }

// Max gets the highest value from a list.
//
// Deprecated: use slices.Max()
//
//go:fix inline
func Max[T cmp.Ordered](list []T) T { return slices.Max(list) }

// Min gets the lowest value from a list.
//
// Deprecated: use slices.Min()
//
//go:fix inline
func Min[T cmp.Ordered](list []T) T { return slices.Min(list) }
