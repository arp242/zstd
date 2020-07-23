package zint

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"strconv"
)

// Uint128 is an unsigned 128-bit integer.
//
// Mostly intended to store UUIDs; storing it as a []byte takes up 40 bytes,
// whereas storing it in two uint64s takes up 16 bytes.
type Uint128 struct{ H, L uint64 }

func NewUint128(b []byte) (Uint128, error) {
	var i Uint128
	return i, i.New(b)
}

func ParseUint128(s string, base int) (Uint128, error) {
	var i Uint128
	return i, i.Parse(s, base)
}

func (u Uint128) String() string { return u.Format(10) }

func (i Uint128) IsZero() bool {
	return i.H == 0 && i.L == 0
}

func (i Uint128) Bytes() ([]byte, error) {
	w := new(bytes.Buffer)
	err := binary.Write(w, binary.BigEndian, i)
	if err != nil {
		return nil, fmt.Errorf("Uint128.Bytes: %w", err)
	}
	if i.H == 0 {
		w.Write(make([]byte, 16))
	}
	return w.Bytes(), nil
}

// TODO: rename and implement fmt.Formatter?
func (u Uint128) Format(base int) string {
	if u.H == 0 {
		return strconv.FormatUint(u.L, base)
	}
	return strconv.FormatUint(u.H, base) + strconv.FormatUint(u.L, base)
}

func (i *Uint128) New(b []byte) error {
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, i)
	if err != nil {
		return fmt.Errorf("uint128.New: %w", err)
	}
	return nil
}

func (i *Uint128) Parse(str string, base int) error {
	var s int
	switch base { // TODO: figure out how to make this generic.
	case 16:
		s = 16
	case 10:
		s = 19
	default:
		return fmt.Errorf("Uint128.Parse: unsupported base: %d", base)
	}
	if len(str) < s {
		return fmt.Errorf("Uint127.Parse: len(%q)=%d, need %d", str, len(str), s)
	}

	l, err := strconv.ParseUint(str[s:], base, 64)
	if err != nil {
		return fmt.Errorf("Uint128.Parse: L: %w", err)
	}
	h, err := strconv.ParseUint(str[:s], base, 64)
	if err != nil {
		return fmt.Errorf("Uint128.Parse: H: %w", err)
	}
	i.L, i.H = l, h
	return nil
}

// Value determines what to store in the DB.
func (i Uint128) Value() (driver.Value, error) {
	return i.Bytes()
}

// Scan converts the data from the DB.
func (i *Uint128) Scan(v interface{}) error {
	if v == nil {
		return nil
	}

	var err error
	switch vv := v.(type) {
	case string:
		err = i.Parse(vv, 10)
	case []byte:
		// TODO: or New?
		//err = i.Parse(string(vv), 10)
		err = i.New(vv)
	case uint64:
		i.L = vv
	case int64:
		i.L = uint64(vv)
	case int:
		i.L = uint64(vv)
	case uint:
		i.L = uint64(vv)
	default:
		err = fmt.Errorf("Uint128.Scan: unknown type: %T", v)
	}
	return err
}

// MarshalText converts the data to a human readable representation.
func (i Uint128) MarshalText() ([]byte, error) {
	return []byte(i.Format(16)), nil
}

// UnmarshalText parses text in to the Go data structure.
func (i *Uint128) UnmarshalText(v []byte) error {
	return i.Parse(string(v), 16)
}
