// Package uint128 adds a uint128 type.
//
// Mostly intended to store UUIDs; storing it as a []byte takes up 40 bytes,
// whereas storing it in two uint64s takes up 16 bytes.
package uint128

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

type Uint128 struct{ H, L uint64 }

func New(b []byte) (Uint128, error) {
	var i Uint128
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &i)
	if err != nil {
		return i, fmt.Errorf("uint128.New: %w", err)
	}
	return i, nil
}

func Parse(s string, base int) (Uint128, error) {
	var i Uint128
	return i, i.Parse(s, base)
}

func (i Uint128) IsZero() bool {
	return i.H == 0 && i.L == 0
}

func (i Uint128) Bytes() ([]byte, error) {
	w := new(bytes.Buffer)
	err := binary.Write(w, binary.BigEndian, i)
	if err != nil {
		return nil, fmt.Errorf("Uint128.Bytes: %w", err)
	}
	return w.Bytes(), nil
}

func (u Uint128) Format(base int) string {
	if u.H == 0 {
		return strconv.FormatUint(u.L, base)
	}
	return strconv.FormatUint(u.H, base) + strconv.FormatUint(u.L, base)
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
