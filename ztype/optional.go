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
type Optional[V any] struct {
	v  V
	ok bool
}

// NewOptional creates a new Optional for the given value.
func NewOptional[V any](v V) Optional[V] {
	return Optional[V]{v, true}
}

func (o Optional[V]) String() string {
	if o.ok {
		return fmt.Sprintf("%v", o.v)
	}
	return "<not set>"
}

// Get the value and a flag indicating if it was set.
//
// The returned value is undefined if it's not set; it can be the type zero
// value, nil, or anything else.
func (o Optional[V]) Get() (V, bool) {
	return o.v, o.ok
}

// Set a value.
func (o *Optional[V]) Set(v V) {
	o.v, o.ok = v, true
}

// Unset this optional.
func (o *Optional[V]) Unset() {
	var v V
	o.v, o.ok = v, false
}

func (o *Optional[V]) Scan(v any) error {
	if v == nil {
		o.ok = false
		return nil
	}

	o.v, o.ok = v.(V)
	if !o.ok {
		return fmt.Errorf("Optional.Scan: unable to scan %#v", v)
	}

	return nil
}

func (o Optional[V]) Value() (driver.Value, error) {
	if !o.ok {
		return nil, nil
	}
	return o.v, nil
}

func (o Optional[V]) MarshalJSON() ([]byte, error) {
	if !o.ok {
		// TODO: we want to be able to also marshal this to the zero type.
		//
		// https://github.com/guregu/null has a separate "zero" package for it. Meh; not
		// a brilliant solution.
		return []byte("null"), nil
	}

	return json.Marshal(o.v)
}

func (o *Optional[V]) UnmarshalJSON(data []byte) error {
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
