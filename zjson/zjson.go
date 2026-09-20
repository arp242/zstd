// Package zjson provides functions for working with JSON.
package zjson

import (
	"encoding/json/v2"
	"fmt"
)

// MustMarshal behaves like json.Marshal but will panic on errors.
func MustMarshal(in any, opts ...json.Options) []byte {
	b, err := json.Marshal(in, opts...)
	if err != nil {
		panic(fmt.Errorf("zjson.MustMarshal: %w", err))
	}
	return b
}

// MustUnmarshal behaves like json.Unmarshal but will panic on errors.
func MustUnmarshal(in []byte, out any, opts ...json.Options) {
	err := json.Unmarshal(in, out, opts...)
	if err != nil {
		panic(fmt.Errorf("zjson.MustUnmarshal: %w", err))
	}
}
