// Package zbyte implements functions for byte slices.
package zbyte

import "bytes"

// Binary reports if this looks like binary data.
//
// Something is considered binary if it contains a NULL byte in the first 8000
// bytes.
//
// This is the same check as git uses; see buffer_is_binary.
func Binary(b []byte) bool {
	for i := range ElideLeft(b, 8000) {
		if b[i] == 0 {
			return true
		}
	}
	return false
}

// ElideLeft returns the "n" left bytes.
func ElideLeft(b []byte, n int) []byte {
	if len(b) > n {
		b = b[:n]
	}
	return b
}

// ElideRight returns the "n" right bytes.
func ElideRight(b []byte, n int) []byte {
	if len(b) > n {
		b = b[len(b)-n:]
	}
	return b
}

// HasSuffixes tests whether the byte slice b ends with any of the suffixes.
//
// Identical to:
//
//	bytes.HasSuffix(s, "one") || bytes.HasSuffix(s, "two")
func HasSuffixes(b []byte, suffixes ...[]byte) bool {
	for _, suf := range suffixes {
		if bytes.HasSuffix(b, suf) {
			return true
		}
	}
	return false
}

// HasPrefixes tests whether the byte slice b starts with any of the prefixes.
//
// Identical to:
//
//	bytes.HasPrefix(s, "one") || bytes.HasPrefix(s, "two")
func HasPrefixes(b []byte, prefixes ...[]byte) bool {
	for _, pre := range prefixes {
		if bytes.HasPrefix(b, pre) {
			return true
		}
	}
	return false
}
