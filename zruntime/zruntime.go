// Package zruntime provides utilities to interface with the Go runtime.
package zruntime

import (
	"os"
	"strings"
)

// Test reports if we're running a go test command.
func Test() bool {
	for _, a := range os.Args[1:] {
		if strings.HasPrefix(a, "-test.") {
			return true
		}
	}
	return false
}

// TestVerbose reports if the test was started with the -v flag.
func TestVerbose() bool {
	for _, a := range os.Args[1:] {
		if a == "-test.v=true" {
			return true
		}
	}
	return false
}
