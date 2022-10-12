// Package zslices implements functions for collections.
package zslices

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"
	"sync"
	"time"

	// Exception to the "stdlib-only" rule, since this should be in stdlib soon.
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
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
		if slices.Contains(list, s) {
			return true
		}
	}
	return false
}

// UniqSort removes duplicate entries from list; the list will be sorted.
func UniqSort[T constraints.Ordered](list []T) []T {
	slices.Sort(list)
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
func IsUniq[T constraints.Ordered](list []T) bool {
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

// Max gets the highest value from a list.
func Max[T constraints.Ordered](list []T) T {
	var max T
	for _, n := range list {
		if n > max {
			max = n
		}
	}
	return max
}

// Min gets the lowest value from a list.
func Min[T constraints.Ordered](list []T) T {
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
			if slices.Contains(o, setItem) {
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
		if slices.Contains(b, v) {
			c = append(c, v)
		}
	}
	return c
}

// SameElements reports if the two slices have the same elements.
//
// This is similar to slices.Equal, but don't take order in to account.
func SameElements[T constraints.Ordered](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	// TODO: there is probably a better way of doing this; dropping the sort
	//       would also allow constraints.Ordered to be comparable.
	aCp := slices.Clone(a)
	bCp := slices.Clone(b)
	slices.Sort(aCp)
	slices.Sort(bCp)
	return slices.Equal(aCp, bCp)
}
