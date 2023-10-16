package zint

import (
	"strconv"
	"strings"
)

// Join a slice of int64 to a comma separated string with the given separator.
//
// Deprecated: use Join()
func Join64(ints []int64, sep string) string {
	s := make([]string, len(ints))
	for i := range ints {
		s[i] = strconv.FormatInt(ints[i], 10)
	}
	return strings.Join(s, sep)
}

// Uniq removes duplicate entries from the list. The list will be sorted.
//
// Deprecated: use zslices.Uniq
func Uniq(list []int64) []int64 {
	var unique []int64
	seen := make(map[int64]struct{})
	for _, l := range list {
		if _, ok := seen[l]; !ok {
			seen[l] = struct{}{}
			unique = append(unique, l)
		}
	}
	return unique
}

// Contains reports whether i is within the list.
//
// Deprecated: use slices.Contains
func Contains(list []int, i int) bool {
	for _, item := range list {
		if item == i {
			return true
		}
	}
	return false
}

// Contains64 reports whether i is within the list.
//
// Deprecated: use slices.Contains
func Contains64(list []int64, i int64) bool {
	for _, item := range list {
		if item == i {
			return true
		}
	}
	return false
}

// MinInt gets the lowest of two numbers.
//
// Deprecated: use Min()
func MinInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// MaxInt gets the highest of two numbers.
//
// Deprecated: use Max()
func MaxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// Difference returns a new slice with elements that are in "set" but not in
// "others".
//
// Deprecated: use zslices.Difference
func Difference(set []int64, others ...[]int64) []int64 {
	out := []int64{}
	for _, setItem := range set {
		found := false
		for _, o := range others {
			if Contains64(o, setItem) {
				found = true
				break
			}
		}

		if !found {
			out = append(out, setItem)
		}
	}

	return out
}
