package zint

import "strconv"

type (
	// Bitflag8 is an uint8 with some extra methods to make bitmask/flag
	// manipulation a bit more convenient.
	Bitflag8 uint8

	// Bitflag16 is an uint16 with some extra methods to make bitmask/flag
	// manipulation a bit more convenient.
	Bitflag16 uint16

	// Bitflag32 is an uint32 with some extra methods to make bitmask/flag
	// manipulation a bit more convenient.
	Bitflag32 uint32

	// Bitflag64 is an uint64 with some extra methods to make bitmask/flag
	// manipulation a bit more convenient.
	Bitflag64 uint64
)

func (f Bitflag8) Has(flag Bitflag8) bool           { return f&flag != 0 }
func (f *Bitflag8) Set(flag Bitflag8)               { *f = *f | flag }
func (f *Bitflag8) Clear(flag Bitflag8)             { *f = *f &^ flag }
func (f *Bitflag8) Toggle(flag Bitflag8)            { *f = *f ^ flag }
func (f *Bitflag8) UnmarshalJSON(text []byte) error { return f.UnmarshalText(text) }
func (f *Bitflag8) UnmarshalText(text []byte) error {
	i, err := strconv.ParseUint(string(text), 10, 8)
	if err != nil {
		return err
	}
	*f |= Bitflag8(i)
	return nil
}

func (f Bitflag16) Has(flag Bitflag16) bool          { return f&flag != 0 }
func (f *Bitflag16) Set(flag Bitflag16)              { *f = *f | flag }
func (f *Bitflag16) Clear(flag Bitflag16)            { *f = *f &^ flag }
func (f *Bitflag16) Toggle(flag Bitflag16)           { *f = *f ^ flag }
func (f *Bitflag16) UnmarshalJSON(text []byte) error { return f.UnmarshalText(text) }
func (f *Bitflag16) UnmarshalText(text []byte) error {
	i, err := strconv.ParseUint(string(text), 10, 16)
	if err != nil {
		return err
	}
	*f |= Bitflag16(i)
	return nil
}

func (f Bitflag32) Has(flag Bitflag32) bool          { return f&flag != 0 }
func (f *Bitflag32) Set(flag Bitflag32)              { *f = *f | flag }
func (f *Bitflag32) Clear(flag Bitflag32)            { *f = *f &^ flag }
func (f *Bitflag32) Toggle(flag Bitflag32)           { *f = *f ^ flag }
func (f *Bitflag32) UnmarshalJSON(text []byte) error { return f.UnmarshalText(text) }
func (f *Bitflag32) UnmarshalText(text []byte) error {
	i, err := strconv.ParseUint(string(text), 10, 32)
	if err != nil {
		return err
	}
	*f |= Bitflag32(i)
	return nil
}

func (f Bitflag64) Has(flag Bitflag64) bool          { return f&flag != 0 }
func (f *Bitflag64) Set(flag Bitflag64)              { *f = *f | flag }
func (f *Bitflag64) Clear(flag Bitflag64)            { *f = *f &^ flag }
func (f *Bitflag64) Toggle(flag Bitflag64)           { *f = *f ^ flag }
func (f *Bitflag64) UnmarshalJSON(text []byte) error { return f.UnmarshalText(text) }
func (f *Bitflag64) UnmarshalText(text []byte) error {
	i, err := strconv.ParseUint(string(text), 10, 64)
	if err != nil {
		return err
	}
	*f |= Bitflag64(i)
	return nil
}
