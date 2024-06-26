// Package zint implements functions for ints.
package zint

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Int with various methods to make conversions easier; useful especially in
// templates etc.
type Int int

func (s Int) String() string   { return strconv.FormatInt(int64(s), 10) }
func (s Int) Int() int         { return int(s) }
func (s Int) Int64() int64     { return int64(s) }
func (s Int) Float32() float32 { return float32(s) }
func (s Int) Float64() float64 { return float64(s) }

// Join a slice of ints to a comma separated string with the given separator.
func Join[T integer](ints []T, sep string) string {
	s := make([]string, len(ints))
	for i := range ints {
		s[i] = strconv.FormatInt(int64(ints[i]), 10)
	}
	return strings.Join(s, sep)
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

// Range creates an []int counting at "start" up to (and including) "end".
func Range(start, end int) []int {
	rng := make([]int, end-start+1)
	for i := 0; i < len(rng); i++ {
		rng[i] = start + i
	}
	return rng
}

// DivideCeil divides two integers and rounds up, rather than down (which is
// what happens when you do int64/int64).
func DivideCeil(count int64, pageSize int64) int64 {
	return int64(math.Ceil(float64(count) / float64(pageSize)))
}

// ToIntSlice converts any []int type to an []int64.
func ToIntSlice(v any) ([]int64, bool) {
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
func ToUintSlice(v any) ([]uint64, bool) {
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
