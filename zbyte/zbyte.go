// Package zbyte implements functions for byte slices.
package zbyte

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
