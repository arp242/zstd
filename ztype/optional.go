package ztype

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Optional represents a value that may or may not exist.
//
// The zero value represents a non-existent value.
//
//	type Strukt struct {
//	    Value zstd.Optional[int]
//	}
//
//	s := Struct{
//	  Value: zstd.NewOptional[]
//	}
//	if v, ok := s.Value.Get(); ok {
//	}
type Optional[Value any] struct {
	v  Value
	ok bool
}

// NewOptional creates a new Optional for the given value.
func NewOptional[Value any](v Value) Optional[Value] {
	return Optional[Value]{v, true}
}

func (o Optional[Value]) String() string {
	if o.ok {
		return fmt.Sprintf("%v", o.v)
	}
	return "<not set>"
}

// Get the value and a flag indicating if it was set.
//
// The returned value is undefined if it's not set; it can be the type zero
// value, nil, or anything else.
func (o Optional[Value]) Get() (Value, bool) {
	return o.v, o.ok
}

// Set a value.
func (o *Optional[Value]) Set(v Value) {
	o.v = v
	o.ok = true
}

// Unset this optional.
func (o *Optional[Value]) Unset() {
	var v Value
	o.v = v
	o.ok = false
}

func (o *Optional[Value]) Scan(v any) error {
	if v == nil {
		o.ok = false
		return nil
	}

	o.v, o.ok = v.(Value)
	if !o.ok {
		return fmt.Errorf("Optional.Scan: unable to scan %#v", v)
	}

	return nil
}

func (o Optional[Value]) Value() (driver.Value, error) {
	if !o.ok {
		return nil, nil
	}
	return o.v, nil
}

func (o Optional[Value]) MarshalJSON() ([]byte, error) {
	if !o.ok {
		return []byte("null"), nil
	}

	return json.Marshal(o.v)
}

func (o *Optional[Value]) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case "null", "undefined":
		return nil
	}

	err := json.Unmarshal(data, &o.v)
	if err != nil {
		return err
	}

	o.ok = true
	return nil
}
