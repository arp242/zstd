package zstring

// Pointer provides a more convenient way to deal with pointers.
//
//   s := "x"
//   p = &s                     →  p = NewPtr("x").P
//
//   if p != nil && p == "v" {  →  if p.Value() == "v" {
//
//   v := "<nil>"
//   if p != nil { v = *p }     →  fmt.Println(p)
type Pointer struct{ P *string }

// NewPointer creates a new Pointer instance set to a pointer of s.
func NewPointer(s string) Pointer { return Pointer{&s} }

// Set Pointer to the pointer of s.
func (p *Pointer) Set(s string) { p.P = &s }

// String gets the string value, or "<nil>" if the pointer is nil.
func (p Pointer) String() string {
	if p.P == nil {
		return "<nil>"
	}
	return string(*p.P)
}

// Value gets the pointer value, or "" if the pointer is nil.
func (p Pointer) Value() string {
	if p.P == nil {
		return ""
	}
	return string(*p.P)
}
