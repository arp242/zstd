// Package zstring implements functions for booleans.
package zbool

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// Bool converts various types to a boolean.
//
// It's always stored as an integer in the database (the only cross-platform way
// in SQL).
//
// Supported types:
//
//	bool
//	int* and float*     0 or 1
//	[]byte and string   "1", "true", "on", "0", "false", "off"
//	nil                 defaults to false
type Bool bool

// Scan converts the data from the DB.
func (b *Bool) Scan(src any) error {
	if b == nil {
		return fmt.Errorf("zdb.Bool: not initialized")
	}

	switch v := src.(type) {
	default:
		return fmt.Errorf("zdb.Bool: unsupported type %T", src)
	case nil:
		*b = false
	case bool:
		*b = Bool(v)
	case int:
		*b = v != 0
	case int8:
		*b = v != 0
	case int16:
		*b = v != 0
	case int32:
		*b = v != 0
	case int64:
		*b = v != 0
	case uint:
		*b = v != 0
	case uint8:
		*b = v != 0
	case uint16:
		*b = v != 0
	case uint32:
		*b = v != 0
	case uint64:
		*b = v != 0
	case float32:
		*b = v != 0
	case float64:
		*b = v != 0

	case []byte, string:
		var text string
		raw, ok := v.([]byte)
		if !ok {
			text = v.(string)
		} else if len(raw) == 1 {
			// Handle the bit(1) column type.
			*b = raw[0] == 1
			return nil
		} else {
			text = string(raw)
		}

		switch strings.TrimSpace(strings.ToLower(text)) {
		case "true", "1", "on":
			*b = true
		case "false", "0", "off":
			*b = false
		default:
			return fmt.Errorf("zdb.Bool: invalid value %q", text)
		}
	}

	return nil
}

// Value converts a bool type into a number to persist it in the database.
func (b Bool) Value() (driver.Value, error) {
	if b {
		return int64(1), nil
	}
	return int64(0), nil
}

// MarshalJSON converts the data to JSON.
func (b Bool) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%t", b)), nil
}

// UnmarshalJSON converts the data from JSON.
func (b *Bool) UnmarshalJSON(text []byte) error {
	switch string(text) {
	case "true", "1", "on":
		*b = true
		return nil
	case "false", "0", "off":
		*b = false
		return nil
	default:
		return fmt.Errorf("zdb.Bool: unknown value: %s", text)
	}
}

// MarshalText converts the data to a human readable representation.
func (b Bool) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%t", b)), nil
}

// UnmarshalText parses text in to the Go data structure.
func (b *Bool) UnmarshalText(text []byte) error {
	if b == nil {
		return fmt.Errorf("zdb.Bool: not initialized")
	}

	switch strings.Trim(strings.TrimSpace(strings.ToLower(string(text))), `"`) {
	case "true", "1", "on":
		*b = true
	case "false", "0", "off":
		*b = false
	default:
		return fmt.Errorf("zdb.Bool: invalid value %q", text)
	}

	return nil
}

func (b Bool) Bool() bool { return bool(b) }
