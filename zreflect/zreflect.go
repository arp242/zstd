// Package zreflect implements functions for reflection.
package zreflect

import (
	"reflect"
	"strings"
)

// Tag splits the tag in to the tag name and options.
func Tag(field reflect.StructField, tag string) (string, []string) {
	t, ok := field.Tag.Lookup(tag)
	if !ok {
		return "", nil
	}

	sp := strings.Split(t, ",")
	if len(sp) > 1 {
		return sp[0], sp[1:]
	}

	return t, nil
}
