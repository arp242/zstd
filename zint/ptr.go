package zint

import "strconv"

// Ptr provides a more convenient way to deal with pointers.
//
//   s := "x"
//   p = &s                     →  p = NewPtr("x").P
//
//   if p != nil && p == "v" {  →  if p.Value() == "v" {
//
//   v := "<nil>"
//   if p != nil { v = *p }     →  fmt.Println(p)
type Ptr struct{ P *int }

// NewPtr creates a new Ptr instance set to a pointer of s.
func NewPtr(s int) Ptr { return Ptr{&s} }

// Set to the pointer of s.
func (p *Ptr) Set(s int) { p.P = &s }

// String gets the string value, or "<nil>" if the pointer is nil.
func (p Ptr) String() string {
	if p.P == nil {
		return "<nil>"
	}
	return strconv.FormatInt(int64(*p.P), 10)
}

// Value gets the pointer value, or 0 if the pointer is nil.
func (p Ptr) Value() int {
	if p.P == nil {
		return 0
	}
	return int(*p.P)
}

// Ptr64 provides a more convenient way to deal with pointers.
//
//   s := "x"
//   p = &s                     →  p = NewPtr("x").P
//
//   if p != nil && p == "v" {  →  if p.Value() == "v" {
//
//   v := "<nil>"
//   if p != nil { v = *p }     →  fmt.Println(p)
type Ptr64 struct{ P *int64 }

// NewPtr64 creates a new Ptr64 instance set to a pointer of s.
func NewPtr64(s int64) Ptr64 { return Ptr64{&s} }

// Set to the pointer of s.
func (p *Ptr64) Set(s int64) { p.P = &s }

// String gets the string value, or "<nil>" if the pointer is nil.
func (p Ptr64) String() string {
	if p.P == nil {
		return "<nil>"
	}
	return strconv.FormatInt(int64(*p.P), 10)
}

// Value gets the pointer value, or 0 if the pointer is nil.
func (p Ptr64) Value() int64 {
	if p.P == nil {
		return 0
	}
	return int64(*p.P)
}
