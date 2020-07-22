package zruntime

import (
	"fmt"
	"runtime"
	"testing"
)

func TestMemStats(t *testing.T) {
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
