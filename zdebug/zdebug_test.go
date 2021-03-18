package zdebug

import (
	"fmt"
	"testing"
)

func TestCallers(t *testing.T) {
	func() {
		for _, c := range Callers() {
			fmt.Println(c)
		}
	}()
}

func TestPrintStack(t *testing.T) {
	func() {
		PrintStack()
	}()
}

func TestLoc(t *testing.T) {
	func() {
		fmt.Println(Loc(0))
		fmt.Println(Loc(1))
		fmt.Println(Loc(2))
	}()
}
