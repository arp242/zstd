package zstring

import (
	"database/sql/driver"
	"fmt"
)

// Ptr provides a more convenient way to deal with pointers.
//
//   s := "x"
//   p = &s                     →  p = zstring.NewPtr("x").P
//
//   if p != nil && p == "v" {  →  if p.Val() == "v" {
//
//   v := "<nil>"
//   if p != nil { v = *p }     →  fmt.Println(p)
type Ptr struct{ P *string }

// NewPtr creates a new Ptr instance set to a pointer of s.
func NewPtr(s string) Ptr { return Ptr{&s} }

// Set to the pointer of s.
func (p *Ptr) Set(s string) { p.P = &s }

// String gets the string value, or "<nil>" if the pointer is nil.
func (p Ptr) String() string {
	if p.P == nil {
		return "<nil>"
	}
	return string(*p.P)
}

// Val gets the pointer value, or "" if the pointer is nil.
func (p Ptr) Val() string {
	if p.P == nil {
		return ""
	}
	return string(*p.P)
}

// Scan converts the data from the DB.
func (p *Ptr) Scan(src interface{}) error {
	if p == nil {
		return fmt.Errorf("zstring.Ptr: not initialized")
	}

	switch v := src.(type) {
	default:
		return fmt.Errorf("zstring.Ptr: unsupported type %T", src)
	case nil:
	case []byte:
		tt := string(v)
		p.P = &tt
	case string:
		p.P = &v
		//*b = false
	}

	return nil
}

// Value converts a bool type into a number to persist it in the database.
func (p Ptr) Value() (driver.Value, error) {
	return p.P, nil
}

// MarshalJSON converts the data to JSON.
func (b Ptr) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%s", b)), nil
}

// UnmarshalJSON converts the data from JSON.
func (p *Ptr) UnmarshalJSON(text []byte) error {
	tt := string(text)
	switch tt {
	case "null", "undefined":
		return nil
	default:
		p.P = &tt
		return nil
	}
}

// MarshalText converts the data to a human readable representation.
func (p Ptr) MarshalText() ([]byte, error) {
	return []byte(p.Val()), nil
}

// UnmarshalText parses text in to the Go data structure.
func (p *Ptr) UnmarshalText(text []byte) error {
	if p == nil {
		return fmt.Errorf("zstring.Ptr: not initialized")
	}

	tt := string(text)
	p.P = &tt
	return nil
}
