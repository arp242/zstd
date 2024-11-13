// Package zstrconv implements conversions to and from string representations.
package zstrconv

import (
	"strconv"
)

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// ParseInt parses an integer. The bitSize is infered from the type parameter.
//
// Other than that, works like [strconv.ParseInt] or [strconv.ParseUint].
func ParseInt[T integer](s string, base int) (T, error) {
	var (
		zero   T
		sz     = 0
		signed = true
	)
	switch any(zero).(type) {
	case int:
	case int8:
		sz = 8
	case int16:
		sz = 16
	case int32:
		sz = 32
	case int64:
		sz = 64
	case uint, uintptr:
		signed = false
	case uint8:
		sz, signed = 8, false
	case uint16:
		sz, signed = 16, false
	case uint32:
		sz, signed = 32, false
	case uint64:
		sz, signed = 64, false
	}

	if signed {
		n, err := strconv.ParseInt(s, base, sz)
		return T(n), err
	}
	n, err := strconv.ParseUint(s, base, sz)
	return T(n), err
}

// MustParseInt works like [ParseInt], but will panic on errors.
func MustParseInt[T integer](s string, base int) T {
	n, err := ParseInt[T](s, base)
	if err != nil {
		panic(err)
	}
	return n
}
