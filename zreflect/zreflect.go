// Package zreflect implements functions for reflection.
package zreflect

import (
	"reflect"
	"strings"

	"zgo.at/zstd/internal/exp/slices"
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

// Fields gets all fields for the struct t, as a slice of names and values, in
// the order they are in the struct.
//
// If tag is not an empty string it will return the value of the tag as the
// field name, or the struct field name if it's not set. Tags with a value of
// "-" will be skipped.
//
// Fields will be skipped if the option given in noskip is set.
//
// It will panic if t is not a struct.
//
// For example:
//
//	t := struct {
//		One   string `db:"one"`
//		Two   string `db:"two,noinsert"`
//		Three int
//	}{"xxx", "yyy", 42}
//
//	Fields(t, "db", "noinsert")
//
// Will return:
//
//	[]string{"one", "Three"}
//	[]any{"xxx", 42}
func Fields(t any, tagname, skip string) (names []string, vals []any) {
	var (
		values = reflect.ValueOf(t)
		types  = reflect.TypeOf(t)
	)
	for values.Kind() == reflect.Ptr {
		values = values.Elem()
		types = types.Elem()
	}

	if tagname == "" && skip != "" {
		panic("zreflect.Fields: setting skip without tagname doesn't make much sense")
	}
	if values.Kind() != reflect.Struct {
		panic("zreflect.Fields: not a struct")
	}

	n := values.NumField()
	names = make([]string, 0, n)
	vals = make([]any, 0, n)
	for i := 0; n > i; i++ {
		t := types.Field(i)

		if t.Type.Kind() == reflect.Struct && t.Anonymous { // Embedded struct
			en, ev := Fields(values.Field(i).Interface(), tagname, skip)
			names = append(names, en...)
			vals = append(vals, ev...)
			continue
		}

		name := t.Name
		if tagname != "" {
			tag := t.Tag.Get(tagname)
			if tag == "-" {
				continue
			}
			tname, opt, _ := strings.Cut(tag, ",")
			if opt != "" && slices.Contains(strings.Split(opt, ","), skip) {
				continue
			}
			if tname != "" {
				name = tname
			}
		}

		names = append(names, name)
		vals = append(vals, values.Field(i).Interface())
	}
	return names, vals
}

// Note these can be made a bit faster by re-implementing them, but without
// adding the information we don't need:
//
//    BenchmarkFields-2        2369833              5082 ns/op            1496 B/op         33 allocs/op
//    BenchmarkNames-2         2690127              4440 ns/op            1024 B/op         31 allocs/op
//    BenchmarkValues-2        2330048              5151 ns/op            1472 B/op         32 allocs/op
//
// But for now, not really worth the effort.

// Names is like [Fields], but only returns the names.
func Names(t any, tagname, skip string) []string {
	names, _ := Fields(t, tagname, skip)
	return names
}

// Values is like [Fields], but only returns the values.
func Values(t any, tagname, skip string) []any {
	_, vals := Fields(t, tagname, skip)
	return vals
}
