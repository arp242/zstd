// Package zslice implements generic functions for slices.
package zslice

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"
	"sort"
	"sync"
	"time"
)

// Choose a random item from the list.
func Choose[T any](list []T) T {
	if len(list) == 0 {
		var t T
		return t
	}

	m := big.NewInt(int64(len(list)))
	n, err := rand.Int(rand.Reader, m)
	if err != nil {
		panic(fmt.Errorf("zcollect.Choose: %w", err))
	}
	return list[n.Int64()]
}

var (
	randSource     *mrand.Rand
	randSourceOnce sync.Once
)

// Shuffle randomizes the order of values.
//
// This uses math/rand, and is not "true random".
func Shuffle[T any](list []T) {
	if len(list) < 2 {
		return
	}

	randSourceOnce.Do(func() {
		randSource = mrand.New(mrand.NewSource(time.Now().UnixNano()))
	})

	randSource.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })
}

// ContainsAny reports whether any of the strings are in the list
func ContainsAny[T comparable](list []T, find ...T) bool {
	for _, s := range find {
		if contains(list, s) {
			return true
		}
	}
	return false
}

// UniqSort removes duplicate entries from list; the list will be sorted.
func UniqSort[T ordered](list []T) []T {
	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
	var last T
	l := list[:0]
	for _, str := range list {
		if str != last {
			l = append(l, str)
		}
		last = str
	}
	return l
}

// Uniq removes duplicate entry from the list.
//
// The order will be preserved, and the first item will be kept. This is slower
// than UniqSort().
func Uniq[T comparable](list []T) []T {
	var unique []T
	seen := make(map[T]struct{})
	for _, l := range list {
		if _, ok := seen[l]; !ok {
			seen[l] = struct{}{}
			unique = append(unique, l)
		}
	}
	return unique
}

// IsUniq reports if the list contains unique values.
func IsUniq[T ordered](list []T) bool {
	return len(list) == len(UniqSort(list))
}

// Repeat returns a slice with the value v repeated n times.
func Repeat[T any](s T, n int) []T {
	r := make([]T, 0, n)
	for i := 0; i < n; i++ {
		r = append(r, s)
	}
	return r
}

// Remove all values from a list.
//
// The return value indicates if this value was found at all.
func Remove[T comparable](l *[]T, name T) bool {
	found := false
	ll := *l
	for i := len(ll) - 1; i >= 0; i-- {
		if ll[i] == name {
			ll = append(ll[:i], ll[i+1:]...)
			found = true
		}
	}
	*l = ll
	return found
}

// RemoveIndexes removes all the given indexes.
//
// The indexes is expected to be sorted from lowest to highest.
//
// Will panic on out of bounds.
func RemoveIndexes[T any](l *[]T, indexes ...int) {
	ll := *l
	for i := len(indexes) - 1; i >= 0; i-- {
		ll = append(ll[:indexes[i]], ll[indexes[i]+1:]...)
	}
	*l = ll
}

// Max gets the highest value from a list.
func Max[T ordered](list []T) T {
	var max T
	for _, n := range list {
		if n > max {
			max = n
		}
	}
	return max
}

// Min gets the lowest value from a list.
func Min[T ordered](list []T) T {
	var min T
	for _, n := range list {
		if n < min {
			min = n
		}
	}
	return min
}

// Difference returns a new slice with elements that are in "set" but not in
// "others".
func Difference[T comparable](set []T, others ...[]T) []T {
	out := make([]T, 0, len(set))
	for _, setItem := range set {
		found := false
		for _, o := range others {
			if contains(o, setItem) {
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

// Intersect returns a new slice with elements that are in both "a" and "b".
func Intersect[T comparable](a, b []T) []T {
	c := make([]T, 0, len(a))
	for _, v := range a {
		if contains(b, v) {
			c = append(c, v)
		}
	}
	return c
}

// SameElements reports if the two slices have the same elements.
//
// This is similar to [slices.Equal], but doesn't take order in to account.
func SameElements[T ordered](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	// TODO: there is probably a better way of doing this; dropping the sort
	//       would also allow cmp.Ordered to be comparable.
	aCp := clone(a)
	bCp := clone(b)
	sort.Slice(aCp, func(i, j int) bool { return aCp[i] < aCp[j] })
	sort.Slice(bCp, func(i, j int) bool { return bCp[i] < bCp[j] })
	return equal(aCp, bCp)
}

// Copy a slice.
//
// This is like [slices.Clone], but you can set a len and cap for the new slice;
// this can be larger than the src slice to prevent copying the array if you're
// appending more items later, or lower if you want to copy and truncate the
// array.
//
// Like make(), this will panic if len > cap.
func Copy[T any](src []T, len, cap int) []T {
	dst := make([]T, len, cap)
	copy(dst, src)
	return dst
}

// AppendCopy is like append(), but ensures the new value is always a copy.
//
// The len and cap will always be set to exactly the len and cap of the new
// array.
func AppendCopy[T any](s []T, app T, more ...T) []T {
	n := Copy(s, len(s), len(s)+len(more)+1)
	return append(append(n, app), more...)
}

// Longest gets the longest string value in this list.
func Longest(list []string) int {
	l := 0
	for _, s := range list {
		if ll := len(s); ll > l {
			l = ll
		}
	}
	return l
}

// LongestFunc gets the longest string value in this list.
//
// for example to get the longest email address from a []mail.Address:
//
//	zslice.LongestFunc(addrs, func(s mail.Address) string { return s.Address })
func LongestFunc[T any](list []T, f func(T) string) int {
	l := 0
	for _, s := range list {
		if ll := len(f(s)); ll > l {
			l = ll
		}
	}
	return l
}

// LastIndex returns the index of the last occurrence of v in s, or -1 if not
// present.
func LastIndex[S ~[]E, E comparable](s S, v E) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == v {
			return i
		}
	}
	return -1
}

// Go 1.19 compat stuff.

type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

func contains[S ~[]E, E comparable](s S, v E) bool {
	for i := range s {
		if v == s[i] {
			return true
		}
	}
	return false
}

func equal[S ~[]E, E comparable](s1, s2 S) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func clone[S ~[]E, E any](s S) S {
	// The s[:0:0] preserves nil in case it matters.
	return append(s[:0:0], s...)
}
