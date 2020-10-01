package zfmt

import (
	"fmt"
	"regexp"
)

// Binary returns the binary representation of a number.
func Binary(c interface{}) string {
	l := ""
	switch c.(type) {
	case int8, uint8:
		l = "8"
	case int16, uint16:
		l = "16"
	case int32, uint32:
		l = "32"
	case int64, uint64, int, uint:
		l = "64"
	default:
		panic(fmt.Sprintf("not a number: %T: %[1]q", c, c))
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

	return fmt.Sprintf("%s", reverse(reBin.ReplaceAllString(reverse(b), `$1$2$3${4}_$5$6$7$8 `)))[1:]
}
