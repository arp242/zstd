package zruntime

import (
	"fmt"
	"runtime"
	"testing"
)

// TODO: actually test this!
func TestMemStats(t *testing.T) {
	t.Skip()

	var m MemStats
	m.Read()
	m.Print(0)
	m.Read()
	m.Read()
	m.Read()
	m.Read()
	m.Read()
	m.Read()
	m.Read()
	m.Read()
	m.Read()
	m.Read()
	m.Read()
	m.Read()
	m.Read()
	fmt.Println()
	runtime.GC()
	m.Print(0)
}
