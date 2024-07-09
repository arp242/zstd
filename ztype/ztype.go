// Package ztype adds extra types.
package ztype

import "reflect"

// Ptr gets a pointer to t.
func Ptr[T any](t T) *T { return &t }

// PtrOrNil gets a pointer to t, or nil if t is the zero value.
func PtrOrNil[T any](t T) *T {
	var zero T
	if reflect.DeepEqual(t, zero) {
		return nil
	}
	return &t
}

// Deref dereferences the pointer v, returning dv if it's nil.
func Deref[T any](v *T, dv T) T {
	if v == nil {
		return dv
	}
	return *v
}
