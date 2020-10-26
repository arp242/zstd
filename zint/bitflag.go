package zint

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

func (f Bitflag8) Has(flag Bitflag8) bool { return f&flag != 0 }
func (f *Bitflag8) Set(flag Bitflag8)     { *f = *f | flag }
func (f *Bitflag8) Clear(flag Bitflag8)   { *f = *f &^ flag }
func (f *Bitflag8) Toggle(flag Bitflag8)  { *f = *f ^ flag }

func (f Bitflag16) Has(flag Bitflag16) bool { return f&flag != 0 }
func (f *Bitflag16) Set(flag Bitflag16)     { *f = *f | flag }
func (f *Bitflag16) Clear(flag Bitflag16)   { *f = *f &^ flag }
func (f *Bitflag16) Toggle(flag Bitflag16)  { *f = *f ^ flag }

func (f Bitflag32) Has(flag Bitflag32) bool { return f&flag != 0 }
func (f *Bitflag32) Set(flag Bitflag32)     { *f = *f | flag }
func (f *Bitflag32) Clear(flag Bitflag32)   { *f = *f &^ flag }
func (f *Bitflag32) Toggle(flag Bitflag32)  { *f = *f ^ flag }

func (f Bitflag64) Has(flag Bitflag64) bool { return f&flag != 0 }
func (f *Bitflag64) Set(flag Bitflag64)     { *f = *f | flag }
func (f *Bitflag64) Clear(flag Bitflag64)   { *f = *f &^ flag }
func (f *Bitflag64) Toggle(flag Bitflag64)  { *f = *f ^ flag }
