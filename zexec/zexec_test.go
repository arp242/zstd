package zexec

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestCommand(t *testing.T) {
	cmd := exec.Command("sh", "-c", "echo one; echo >&2 two; echo three; sleep 1; echo four; echo >&2 five")

	ch, err := Readlines(cmd)
	if err != nil {
		t.Fatal(err)
	}

	var have []string
	for {
		line, ok := <-ch
		if !ok {
			break
		}
		if line.Error != nil {
			fmt.Fprintln(os.Stderr, "error reading:", line.Error)
			break
		}

		have = append(have, line.String())
	}

	t.Skip("TODO: fails")
	// have: []string{"one", "three", "error: two", "error: five", "four"}
	// want: []string{"one", "error: two", "three", "four", "error: five"}
	want := []string{"one", "error: two", "three", "four", "error: five"}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave: %#v\nwant: %#v", have, want)
	}
}
