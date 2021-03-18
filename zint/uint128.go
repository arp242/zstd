package zint

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

// Uint128 is an unsigned 128-bit integer.
//
// Mostly intended to store UUIDs; storing it as a []byte takes up 40 bytes,
// whereas storing it in two uint64s takes up 16 bytes.
//
// It's stored in so-called "big endian" format; that is, the most significant
// bit is on the right.
type Uint128 [2]uint64

// NewUint128 creates a new uint128 from a [16]byte.
func NewUint128(b []byte) (Uint128, error) {
	var i Uint128
	return i, i.New(b)
}

func ParseUint128(s string, base int) (Uint128, error) {
	var i Uint128
	return i, i.Parse(s, base)
}

func (i Uint128) String() string { return i.Format(16) }
func (i Uint128) IsZero() bool   { return i[0] == 0 && i[1] == 0 }

// Format according to the given base.
//
// TODO: this is not really printin a number, but just printing the 2 numbers
// side-by-side, rather than actually adding up the bits.
func (i Uint128) Format(base int) string {
	return strconv.FormatUint(i[0], base) + "-" + strconv.FormatUint(i[1], base)
}

// New sets this uint128 from a [16]byte.
func (i *Uint128) New(b []byte) error {
	if len(b) != 16 {
		return fmt.Errorf("wrong length: %d; need 16", len(b))
	}

	_ = b[15] // bounds check hint to compiler; see golang.org/issue/14808
	i[0] = uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
	i[1] = uint64(b[15]) | uint64(b[14])<<8 | uint64(b[13])<<16 | uint64(b[12])<<24 |
		uint64(b[11])<<32 | uint64(b[10])<<40 | uint64(b[9])<<48 | uint64(b[8])<<56
	return nil
}

func (i *Uint128) Parse(str string, base int) error {
	d := strings.Index(str, "-")
	if d == -1 {
		return fmt.Errorf("*Uint128.Parse: invalid format")
	}

	l, err := strconv.ParseUint(str[d+1:], base, 64)
	if err != nil {
		return fmt.Errorf("Uint128.Parse: L: %w", err)
	}
	h, err := strconv.ParseUint(str[:d], base, 64)
	if err != nil {
		return fmt.Errorf("Uint128.Parse: H: %w", err)
	}
	i[1], i[0] = l, h
	return nil
}

func (i Uint128) Bytes() []byte {
	b := make([]byte, 16)
	_ = b[15] // bounds check hint to compiler; see golang.org/issue/14808
	b[0] = byte(i[0] >> 56)
	b[1] = byte(i[0] >> 48)
	b[2] = byte(i[0] >> 40)
	b[3] = byte(i[0] >> 32)
	b[4] = byte(i[0] >> 24)
	b[5] = byte(i[0] >> 16)
	b[6] = byte(i[0] >> 8)
	b[7] = byte(i[0])

	b[8] = byte(i[1] >> 56)
	b[9] = byte(i[1] >> 48)
	b[10] = byte(i[1] >> 40)
	b[11] = byte(i[1] >> 32)
	b[12] = byte(i[1] >> 24)
	b[13] = byte(i[1] >> 16)
	b[14] = byte(i[1] >> 8)
	b[15] = byte(i[1])
	return b
}

// Value determines what to store in the DB.
func (i Uint128) Value() (driver.Value, error) { return i.Bytes(), nil }

// Scan converts the data from the DB.
func (i *Uint128) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	b, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("Uint128.Scan: must be []byte, not %T", v)
	}
	if len(b) == 0 {
		return nil
	}
	return i.New(b)
}

// MarshalText converts the data to a human readable representation.
func (i Uint128) MarshalText() ([]byte, error) { return []byte(i.Format(16)), nil }

// UnmarshalText parses text in to the Go data structure.
func (i *Uint128) UnmarshalText(v []byte) error {
	if len(v) == 0 {
		return nil
	}
	return i.Parse(string(v), 16)
}
