package zstring

import (
	"math/rand"
	"sort"
	"strings"
	"time"
)

// Uniq removes duplicate entries from list; the list will be sorted.
//
// Deprecated: use zslices.Uniq
func Uniq(list []string) []string {
	sort.Strings(list)
	var last string
	l := list[:0]
	for _, str := range list {
		if str != last {
			l = append(l, str)
		}
		last = str
	}
	return l
}

// Contains reports whether str is within the list.
//
// Deprecated: use slices.Contains
func Contains(list []string, str string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}

// ContainsAny reports whether any of the strings are in the list
//
// Deprecated: use zslices.ContainsAny
func ContainsAny(list []string, strs ...string) bool {
	for _, s := range strs {
		if Contains(list, s) {
			return true
		}
	}
	return false
}

// Repeat returns a slice with the string s repeated n times.
//
// Deprecated: use zslices.Repeat
func Repeat(s string, n int) (r []string) {
	for i := 0; i < n; i++ {
		r = append(r, s)
	}
	return r
}

// Choose chooses a random item from the list.
//
// Deprecated: use zslices.Choose.
func Choose(l []string) string {
	if len(l) == 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	return l[rand.Intn(len(l))]
}

// Difference returns a new slice with elements that are in "set" but not in
// "others".
//
// Deprecated: use zslices.Difference
func Difference(set []string, others ...[]string) []string {
	out := []string{}
	for _, setItem := range set {
		found := false
		for _, o := range others {
			if Contains(o, setItem) {
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

// Split2 splits a string with strings.SplitN(.., 2) and returns the result.
//
// This makes some string splits a bit more elegant:
//
//	key, value := zstring.Split2(line, "=")
//
// Deprecated: use strings.Cut
func Split2(str, sep string) (string, string) {
	s := strings.SplitN(str, sep, 2)
	if len(s) == 1 {
		return s[0], ""
	}
	return s[0], s[1]
}

// Split3 splits a string with strings.SplitN(.., 3) and returns the result.
//
// Deprecated: never used it.
func Split3(str, sep string) (string, string, string) {
	s := strings.SplitN(str, sep, 3)
	if len(s) < 3 {
		m := make([]string, 3)
		copy(m, s)
		s = m
	}
	return s[0], s[1], s[2]
}

// Split4 splits a string with strings.SplitN(.., 4) and returns the result.
//
// Deprecated: never used it.
func Split4(str, sep string) (string, string, string, string) {
	s := strings.SplitN(str, sep, 4)
	if len(s) < 4 {
		m := make([]string, 4)
		copy(m, s)
		s = m
	}
	return s[0], s[1], s[2], s[3]
}

// Remove all values from a list.
//
// The return value indicates if this value was found at all.
//
// Deprecated: use zslices.Remove
func Remove(l *[]string, name string) bool {
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
