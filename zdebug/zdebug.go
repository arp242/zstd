// Package zdebug implements functions useful when debugging programs.
package zdebug

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"zgo.at/zstd/zruntime"
)

// Stack gets a stack trace.
//
// Unlike debug.Stack() the output is much more concise: every frame is a single
// line with the package/function name and file location printed in aligned
// columns.
func Stack(filterFun ...string) []byte {
	var (
		callers = zruntime.Callers(filterFun...)
		rows    = make([][]any, 0, len(callers))
		width   = 0
	)
	for _, f := range callers {
		loc := filepath.Base(f.File) + ":" + strconv.Itoa(f.Line)
		if len(loc) > width {
			width = len(loc)
		}
		rows = append(rows, []any{loc, f.Function})
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
func PrintStack(filterFun ...string) {
	fmt.Fprint(os.Stderr, string(Stack(filterFun...)))
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
