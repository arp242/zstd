// Package zfmt implements additional formatting functions.
package zfmt

import (
	"fmt"
	"math/bits"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Binary returns a nicely formatted binary representation of a number.
func Binary(c any) string {
	t := reflect.TypeOf(c)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		c = reflect.ValueOf(c).Elem().Interface()
	}

	l := ""
	switch t.Kind() {
	case reflect.Int8, reflect.Uint8:
		l = "8"
	case reflect.Int16, reflect.Uint16:
		l = "16"
	case reflect.Int32, reflect.Uint32:
		l = "32"
	case reflect.Int64, reflect.Uint64:
		l = "64"
	case reflect.Int, reflect.Uint:
		l = strconv.Itoa(bits.UintSize)
	default:
		panic(fmt.Sprintf("zfmt.Binary: not a number but %T: %[1]v", c, c))
	}

	b := fmt.Sprintf("%0"+l+"b", c)
	if b[0] == '-' {
		b = "1" + b[1:]
	}

	reverse := func(s string) string {
		r := []rune(s)
		for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
			r[i], r[j] = r[j], r[i]
		}
		return string(r)
	}
	reBin := regexp.MustCompile(`([01])([01])([01])([01])([01])([01])([01])([01])`)
	return reverse(reBin.ReplaceAllString(reverse(b), `$1$2$3${4}_$5$6$7$8 `))[1:]
}

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Number returns a formatted representation of n, using thousandsSep to
// separate thousands.
//
// If thousandsSep is '.' it will use ',' as the fraction separator, otherwise
// it will default to '.'.
func Number[T number](n T, thousandsSep rune) string {
	var s string
	switch any(n).(type) {
	case int, int8, int16, int32, int64:
		s = strconv.FormatInt(int64(n), 10)
	case uint, uint8, uint16, uint32, uint64:
		s = strconv.FormatUint(uint64(n), 10)
	case float32, float64:
		s = strconv.FormatFloat(float64(n), 'f', -1, 64)
	}

	s, d, _ := strings.Cut(s, ".")
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	var out []rune
	for i := range b {
		if i > 0 && i%3 == 0 && thousandsSep > 1 {
			out = append(out, thousandsSep)
		}
		out = append(out, rune(b[i]))
	}

	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}

	s = string(out)
	if d != "" {
		if thousandsSep == '.' {
			s += "," + d
		} else {
			s += "." + d
		}
	}
	return s
}
