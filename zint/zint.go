// Package zint implements functions for ints.
package zint

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"zgo.at/zstd/zstrconv"
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

// Split a string to a slice of integers.
func Split[T integer](s string, sep string) ([]T, error) {
	s = strings.Trim(s, " \t\n"+sep)
	if len(s) == 0 {
		return nil, nil
	}

	items := strings.Split(s, sep)
	ret := make([]T, len(items))
	for i := range items {
		val, err := zstrconv.ParseInt[T](strings.TrimSpace(items[i]), 10)
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
func Fields[T integer](s string) ([]T, error) {
	sf := strings.Fields(s)
	nf := make([]T, len(sf))
	for i, f := range sf {
		n, err := zstrconv.ParseInt[T](f, 0)
		if err != nil {
			return nil, fmt.Errorf("zint.Fields: parsing entry %d in %q: %w", i, s, err)
		}
		nf[i] = n
	}

	return nf, nil
}
