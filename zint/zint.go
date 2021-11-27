// Package zint implements functions for ints.
package zint

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Int with various methods to make conversions easier; useful especially in
// templates etc.
type Int int

func (s Int) String() string   { return strconv.FormatInt(int64(s), 10) }
func (s Int) Int() int         { return int(s) }
func (s Int) Int64() int64     { return int64(s) }
func (s Int) Float32() float32 { return float32(s) }
func (s Int) Float64() float64 { return float64(s) }

// Join a slice of ints to a comma separated string with the given separator.
func Join(ints []int, sep string) string {
	s := make([]string, len(ints))
	for i := range ints {
		s[i] = strconv.Itoa(ints[i])
	}
	return strings.Join(s, sep)
}

// Join a slice of int64 to a comma separated string with the given separator.
func Join64(ints []int64, sep string) string {
	s := make([]string, len(ints))
	for i := range ints {
		s[i] = strconv.FormatInt(ints[i], 10)
	}
	return strings.Join(s, sep)
}

// Uniq removes duplicate entries from the list. The list will be sorted.
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

// Split a string to a slice of []int64.
func Split(s string, sep string) ([]int64, error) {
	s = strings.Trim(s, " \t\n"+sep)
	if len(s) == 0 {
		return nil, nil
	}

	items := strings.Split(s, sep)
	ret := make([]int64, len(items))
	for i := range items {
		val, err := strconv.ParseInt(strings.TrimSpace(items[i]), 10, 64)
		if err != nil {
			return nil, err
		}
		ret[i] = val
	}

	return ret, nil
}

// Contains reports whether i is within the list.
func Contains(list []int, i int) bool {
	for _, item := range list {
		if item == i {
			return true
		}
	}
	return false
}

// Contains64 reports whether i is within the list.
func Contains64(list []int64, i int64) bool {
	for _, item := range list {
		if item == i {
			return true
		}
	}
	return false
}

// Range creates an []int counting at "start" up to (and including) "end".
func Range(start, end int) []int {
	rng := make([]int, end-start+1)
	for i := 0; i < len(rng); i++ {
		rng[i] = start + i
	}
	return rng
}

// Fiter a list.
//
// The function will be called for every item and those that return false will
// not be included in the return value.
func Filter(list []int64, fun func(int64) bool) []int64 {
	var ret []int64
	for _, e := range list {
		if fun(e) {
			ret = append(ret, e)
		}
	}

	return ret
}

// FilterEmpty is a filter for Filter() to remove empty entries.
func FilterEmpty(e int64) bool { return e != 0 }

// Min gets the lowest of two numbers.
func Min(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

// Max gets the highest of two numbers.
func Max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

// MinInt gets the lowest of two numbers.
func MinInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// MaxInt gets the highest of two numbers.
func MaxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// NonZero returns the first argument that is not 0. It will return 0 if all
// arguments are 0.
func NonZero(a, b int64, c ...int64) int64 {
	if a != 0 {
		return a
	}
	if b != 0 {
		return b
	}
	for i := range c {
		if c[i] != 0 {
			return c[i]
		}
	}
	return 0
}

// DivideCeil divides two integers and rounds up, rather than down (which is
// what happens when you do int64/int64).
func DivideCeil(count int64, pageSize int64) int64 {
	return int64(math.Ceil(float64(count) / float64(pageSize)))
}

// Differencereturns a new slice with elements that are in "set" but not in
// "others".
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

// ToIntSlice converts any []int type to an []int64.
func ToIntSlice(v interface{}) ([]int64, bool) {
	var r []int64
	switch vv := v.(type) {
	case []int64:
		r = vv
	case []uint64:
		vvv := make([]int64, len(vv))
		for i := range vv {
			vvv[i] = int64(vv[i])
		}
		r = vvv

	case []int8:
		vvv := make([]int64, len(vv))
		for i := range vv {
			vvv[i] = int64(vv[i])
		}
		r = vvv
	case []int16:
		vvv := make([]int64, len(vv))
		for i := range vv {
			vvv[i] = int64(vv[i])
		}
		r = vvv
	case []int32:
		vvv := make([]int64, len(vv))
		for i := range vv {
			vvv[i] = int64(vv[i])
		}
		r = vvv
	case []int:
		vvv := make([]int64, len(vv))
		for i := range vv {
			vvv[i] = int64(vv[i])
		}
		r = vvv
	case []uint8:
		vvv := make([]int64, len(vv))
		for i := range vv {
			vvv[i] = int64(vv[i])
		}
		r = vvv
	case []uint16:
		vvv := make([]int64, len(vv))
		for i := range vv {
			vvv[i] = int64(vv[i])
		}
		r = vvv
	case []uint32:
		vvv := make([]int64, len(vv))
		for i := range vv {
			vvv[i] = int64(vv[i])
		}
		r = vvv
	case []uint:
		vvv := make([]int64, len(vv))
		for i := range vv {
			vvv[i] = int64(vv[i])
		}
		r = vvv
	}

	return r, r != nil
}

// ToUintSlice converts any []int type to an []uint64.
func ToUintSlice(v interface{}) ([]uint64, bool) {
	var r []uint64
	switch vv := v.(type) {
	case []int64:
		vvv := make([]uint64, len(vv))
		for i := range vv {
			vvv[i] = uint64(vv[i])
		}
		r = vvv
	case []uint64:
		r = vv

	case []int8:
		vvv := make([]uint64, len(vv))
		for i := range vv {
			vvv[i] = uint64(vv[i])
		}
		r = vvv
	case []int16:
		vvv := make([]uint64, len(vv))
		for i := range vv {
			vvv[i] = uint64(vv[i])
		}
		r = vvv
	case []int32:
		vvv := make([]uint64, len(vv))
		for i := range vv {
			vvv[i] = uint64(vv[i])
		}
		r = vvv
	case []int:
		vvv := make([]uint64, len(vv))
		for i := range vv {
			vvv[i] = uint64(vv[i])
		}
		r = vvv
	case []uint8:
		vvv := make([]uint64, len(vv))
		for i := range vv {
			vvv[i] = uint64(vv[i])
		}
		r = vvv
	case []uint16:
		vvv := make([]uint64, len(vv))
		for i := range vv {
			vvv[i] = uint64(vv[i])
		}
		r = vvv
	case []uint32:
		vvv := make([]uint64, len(vv))
		for i := range vv {
			vvv[i] = uint64(vv[i])
		}
		r = vvv
	case []uint:
		vvv := make([]uint64, len(vv))
		for i := range vv {
			vvv[i] = uint64(vv[i])
		}
		r = vvv
	}

	return r, r != nil
}

// RoundToPowerOf2 rounds up to the nearest power of 2.
//
// https://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2
func RoundToPowerOf2(n uint64) uint64 {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	n++
	return n
}

// Fields splits a strings with strings.Fields() and parses each entry as an
// integer.
func Fields(s string) ([]int64, error) {
	sf := strings.Fields(s)
	nf := make([]int64, len(sf))
	for i, f := range sf {
		n, err := strconv.ParseInt(f, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("zint.Fields: parsing entry %d in %q: %w", i, s, err)
		}
		nf[i] = n
	}

	return nf, nil
}
