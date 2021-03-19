// Package zfmt implements additional formatting functions.
package zfmt

import (
	"fmt"
	"math/bits"
	"reflect"
	"regexp"
	"strconv"
)

// Binary returns a nicely formatted binary representation of a number.
func Binary(c interface{}) string {
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
