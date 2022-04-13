// Package zstd is a collection of extensions to Go's stdlib.
package zstd

// Ptr gets a pointer to t.
func Ptr[T any](t T) *T { return &t }

// Deref dereferences the pointer v, returning dv if it's nil.
func Deref[T any](v *T, dv T) T {
	if v == nil {
		return dv
	}
	return *v
}
