package zint

import "strconv"

// Pointer provides a more convenient way to deal with pointers.
//
//   s := "x"
//   p = &s                     →  p = NewPtr("x").P
//
//   if p != nil && p == "v" {  →  if p.Value() == "v" {
//
//   v := "<nil>"
//   if p != nil { v = *p }     →  fmt.Println(p)
type Pointer struct{ P *int }

// NewPointer creates a new Pointer instance set to a pointer of s.
func NewPointer(s int) Pointer { return Pointer{&s} }

// Set Pointer to the pointer of s.
func (p *Pointer) Set(s int) { p.P = &s }

// String gets the string value, or "<nil>" if the pointer is nil.
func (p Pointer) String() string {
	if p.P == nil {
		return "<nil>"
	}
	return strconv.FormatInt(int64(*p.P), 10)
}

// Value gets the pointer value, or 0 if the pointer is nil.
func (p Pointer) Value() int {
	if p.P == nil {
		return 0
	}
	return int(*p.P)
}

// Pointer64 provides a more convenient way to deal with pointers.
//
//   s := "x"
//   p = &s                     →  p = NewPtr("x").P
//
//   if p != nil && p == "v" {  →  if p.Value() == "v" {
//
//   v := "<nil>"
//   if p != nil { v = *p }     →  fmt.Println(p)
type Pointer64 struct{ P *int64 }

// NewPointer64 creates a new Pointer64 instance set to a pointer of s.
func NewPointer64(s int64) Pointer64 { return Pointer64{&s} }

// Set Pointer64 to the pointer of s.
func (p *Pointer64) Set(s int64) { p.P = &s }

// String gets the string value, or "<nil>" if the pointer is nil.
func (p Pointer64) String() string {
	if p.P == nil {
		return "<nil>"
	}
	return strconv.FormatInt(int64(*p.P), 10)
}

// Value gets the pointer value, or 0 if the pointer is nil.
func (p Pointer64) Value() int64 {
	if p.P == nil {
		return 0
	}
	return int64(*p.P)
}
