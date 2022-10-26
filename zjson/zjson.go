// Package zjson provides functions for working with JSON.
package zjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Int for APIs that return numbers as strings.
type Int int64

// Marshal in to JSON.
func (i Int) MarshalJSON() ([]byte, error) {
	return append(append([]byte(`"`), strconv.FormatInt(int64(i), 10)...), '"'), nil
}

// Unmarshal a string timestamp as an int.
func (i *Int) UnmarshalJSON(v []byte) error {
	ii, err := strconv.ParseInt(string(bytes.Trim(v, `"`)), 10, 64)
	*i = Int(ii)
	return err
}

// Timestamp for APIs that return dates as a numeric Unix timestamp.
type Timestamp struct{ time.Time }

// Marshal in to JSON.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.Unix(), 10)), nil
}

// Unmarshal a Unix timestamp as a date.
func (t *Timestamp) UnmarshalJSON(v []byte) error {
	t.Time = time.Time{}

	vv := string(v)
	if vv == "null" || vv == "undefined" || vv == "" || vv == "0" {
		return nil
	}

	n, err := strconv.ParseInt(vv, 10, 64)
	if err != nil {
		return fmt.Errorf("Timestamp.UnmarshalJSON %q: %w", vv, err)
	}
	if n > 0 { // Make sure that IsZero() works.
		t.Time = time.Unix(n, 0).UTC()
	}
	return nil
}

// UnmarshalTo unmarshals the JSON in data to a new instance of target.
//
// Target can be any type that json.Unmarshal can unmarshal to, and doesn't need
// to be a pointer. The returned value is always a pointer.
func UnmarshalTo(data []byte, target reflect.Type) (any, error) {
	if target == nil {
		return nil, errors.New("zjson.UnmarshalTo: target is nil")
	}
	for target.Kind() == reflect.Ptr {
		target = target.Elem()
	}

	t := reflect.New(target).Interface()
	return t, json.Unmarshal(data, t)
}

// MustUnmarshalTo behaves like UnmarshalTo but will panic on errors.
func MustUnmarshalTo(data []byte, target reflect.Type) any {
	t, err := UnmarshalTo(data, target)
	if err != nil {
		panic(err)
	}
	return t
}

// MustMarshal behaves like json.Marshal but will panic on errors.
func MustMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// MustMarshalIndent behaves like json.MarshalIndent but will panic on errors.
func MustMarshalIndent(v any, prefix, indent string) []byte {
	b, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		panic(err)
	}
	return b
}

// MustUnmarshal behaves like json.Unmarshal but will panic on errors.
func MustUnmarshal(data []byte, v any) {
	err := json.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}

// Indent a json string by unmarshalling it and marshalling it with
// MarshalIndent.
//
// The data will be unmarshalled in to v, which must be a pointer. Example:
//
//	Indent(`{"a": "b"}`, &map[string]string{}, "", "  ")
func Indent(data []byte, v any, prefix, indent string) ([]byte, error) {
	err := json.Unmarshal(data, v)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(v, prefix, indent)
}

// MustIndent behaves like Indent but will panic on errors.
func MustIndent(data []byte, v any, prefix, indent string) []byte {
	b, err := Indent(data, v, prefix, indent)
	if err != nil {
		panic(err)
	}
	return b
}

// MustMarshalString is like MustMarshal, but returns a string.
func MustMarshalString(v any) string {
	return string(MustMarshal(v))
}

// MustMarshalIndentString is like MustMarshalIndent, but returns a string.
func MustMarshalIndentString(v any, prefix, indent string) string {
	return string(MustMarshalIndent(v, prefix, indent))
}

// IndentString is like Indent, but returns a string.
func IndentString(data []byte, v any, prefix, indent string) (string, error) {
	b, err := Indent(data, v, prefix, indent)
	return string(b), err
}

// MustIndentString is like MustIndent, but returns a string.
func MustIndentString(data []byte, v any, prefix, indent string) string {
	return string(MustIndent(data, v, prefix, indent))
}
