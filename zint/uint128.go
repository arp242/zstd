package zint

import (
	"fmt"
	"strconv"
)

// Uint128 is an unsigned 128-bit integer.
//
// Mostly intended to store UUIDs; storing it as a []byte takes up 40 bytes,
// whereas storing it in two uint64s takes up 16 bytes.
type Uint128 struct{ H, L uint64 }

// NewUint128 creates a new uint128 from a [16]byte.
func NewUint128(b []byte) (Uint128, error) {
	var i Uint128
	return i, i.New(b)
}

func (u Uint128) String() string { return u.Format(16) }
func (i Uint128) IsZero() bool   { return i.H == 0 && i.L == 0 }
func (u Uint128) Format(base int) string {
	return strconv.FormatUint(u.H, base) + strconv.FormatUint(u.L, base)
}

// New sets this uint128 from a [16]byte.
func (i *Uint128) New(b []byte) error {
	if len(b) != 16 {
		return fmt.Errorf("wrong length: %d; need 16", len(b))
	}

	_ = b[15] // bounds check hint to compiler; see golang.org/issue/14808
	i.H = uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
	i.L = uint64(b[15]) | uint64(b[14])<<8 | uint64(b[13])<<16 | uint64(b[12])<<24 |
		uint64(b[11])<<32 | uint64(b[10])<<40 | uint64(b[9])<<48 | uint64(b[8])<<56
	return nil
}

func (i Uint128) Bytes() []byte {
	b := make([]byte, 16)
	_ = b[15] // bounds check hint to compiler; see golang.org/issue/14808
	b[0] = byte(i.H >> 56)
	b[1] = byte(i.H >> 48)
	b[2] = byte(i.H >> 40)
	b[3] = byte(i.H >> 32)
	b[4] = byte(i.H >> 24)
	b[5] = byte(i.H >> 16)
	b[6] = byte(i.H >> 8)
	b[7] = byte(i.H)
	b[8] = byte(i.L >> 56)
	b[9] = byte(i.L >> 48)
	b[10] = byte(i.L >> 40)
	b[11] = byte(i.L >> 32)
	b[12] = byte(i.L >> 24)
	b[13] = byte(i.L >> 16)
	b[14] = byte(i.L >> 8)
	b[15] = byte(i.L)
	return b
}

// Value determines what to store in the DB.
//func (i Uint128) Value() (driver.Value, error) {
func (i Uint128) Value() (interface{}, error) { return i.Bytes(), nil }

// Scan converts the data from the DB.
func (i *Uint128) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	b, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("Uint128.Scan: must be []byte, not %T", v)
	}
	return i.New(b)
}
