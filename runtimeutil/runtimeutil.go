// Package runtimeutil provides utilities to interface with the Go runtime.
package runtimeutil

import (
	"os"
	"strings"
)

// IsTest reports if we're running a go test command.
func IsTest() bool {
	for _, a := range os.Args[1:] {
		if strings.HasPrefix(a, "-test.") {
			return true
		}
	}
	return false
}
