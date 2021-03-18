// Package zdebug implements functions useful when debugging programs.
package zdebug

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"zgo.at/zstd/zstring"
)

type Caller struct {
	Function string
	File     string
	Line     int
}

// Callers gets a list of callers.
func Callers() []Caller {
	var (
		pc     = make([]uintptr, 50)
		n      = runtime.Callers(2, pc)
		frames = runtime.CallersFrames(pc[:n])
	)

	ret := make([]Caller, 0, n)
	for f, more := frames.Next(); more; f, more = frames.Next() {
		if zstring.Contains([]string{
			"runtime.goexit", "testing.tRunner",
			"zgo.at/zstd/zdebug.Stack",
			"zgo.at/zstd/zdebug.PrintStack",
		}, f.Function) {
			continue
		}
		ret = append(ret, Caller{
			File:     f.File,
			Line:     f.Line,
			Function: f.Function,
		})
	}
	return ret
}

// Stack gets a stack trace.
//
// Unlike debug.Stack() the output is much more concise: every frame is a single
// line with the package/function name and file location printed in aligned
// columns.
func Stack() []byte {
	var (
		callers = Callers()
		rows    = make([][]interface{}, 0, len(callers))
		width   = 0
	)
	for _, f := range Callers() {
		loc := filepath.Base(f.File) + ":" + strconv.Itoa(f.Line)
		if len(loc) > width {
			width = len(loc)
		}
		rows = append(rows, []interface{}{loc, f.Function})
	}

	buf := new(bytes.Buffer)
	f := fmt.Sprintf("\t%%-%ds   %%s\n", width)
	for _, r := range rows {
		fmt.Fprintf(buf, f, r...)
	}
	return buf.Bytes()
}

// PrintStack prints a stack trace to stderr.
//
// Unlike debug.PrintStack() the output is much more concise: every frame is a
// single line with the package/function name and file location printed in
// aligned columns.
func PrintStack() {
	fmt.Fprint(os.Stderr, string(Stack()))
}

// Loc gets a location in the stack trace.
//
// Use 0 for the current location; 1 for one up, etc.
func Loc(n int) string {
	_, file, line, ok := runtime.Caller(n + 1)
	if !ok {
		file = "???"
		line = 0
	}

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	return fmt.Sprintf("%v:%v", file, line)
}
