// Package zcollect implements functions for collections.
package zcollect

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Choose chooses a random item from the list.
func Choose[T any](list []T) T {
	if len(list) == 0 {
		var t T
		return t
	}

	m := big.NewInt(int64(len(list)))
	n, err := rand.Int(rand.Reader, m)
	if err != nil {
		panic(fmt.Errorf("zcollect.Choose: %w", err))
	}
	return list[n.Int64()]
}
