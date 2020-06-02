// Package zbyte implements functions for byte slices.
package zbyte

import "unicode/utf8"

// Binary reports if this looks like binary data.
func Binary(b []byte) bool {
	for i := range b {
		if (b[i] <= 0x1f && b[i] != 0x09 && b[i] != 0x10 && b[i] != 0x13) || b[i] == 0xff {
			return true
		}
	}
	return !utf8.Valid(b)
}
