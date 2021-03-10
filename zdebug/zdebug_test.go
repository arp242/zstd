package zdebug

import (
	"fmt"
	"testing"
)

func TestPrintStack(t *testing.T) {
	PrintStack()
}

func TestLoc(t *testing.T) {
	fmt.Println(Loc(0))
	fmt.Println(Loc(1))
}
