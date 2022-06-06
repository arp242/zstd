package ztime

import (
	"testing"
	"time"
)

// Now returns the current time as UTC.
//
// This can be swapped out in tests with SetNow()
//
// TODO: this shouldn't use .UTC(), at least not by default.
var Now = func() time.Time { return time.Now().UTC() }

// SetNow sets Now() and restores it when the test finishes.
//
// The date is parsed with FromString().
func SetNow(t *testing.T, s string) {
	t.Helper()

	d := FromString(s)
	Now = func() time.Time { return d }
	t.Cleanup(func() {
		Now = func() time.Time { return time.Now().UTC() }
	})
}
