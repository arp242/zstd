// Package ztest contains helper functions that are useful for writing tests.
package ztest

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// ErrorContains checks if the error message in out contains the text in
// want.
//
// This is safe when out is nil. Use an empty string for want if you want to
// test that err is nil.
func ErrorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}

// Read data from a file.
func Read(t *testing.T, paths ...string) []byte {
	t.Helper()

	path := filepath.Join(paths...)
	file, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read %v: %v", path, err)
	}
	return file
}

// TempFile creates a new temporary file and returns the path and a clean
// function to remove it.
//
//  f, clean := TempFile("some\ndata")
//  defer clean()
func TempFile(t *testing.T, data string) (string, func()) {
	t.Helper()

	fp, err := os.CreateTemp(os.TempDir(), "gotest")
	if err != nil {
		t.Fatalf("test.TempFile: could not create file in %v: %v", os.TempDir(), err)
	}

	defer func() {
		err := fp.Close()
		if err != nil {
			t.Fatalf("test.TempFile: close: %v", err)
		}
	}()

	_, err = fp.WriteString(data)
	if err != nil {
		t.Fatalf("test.TempFile: write: %v", err)
	}

	return fp.Name(), func() {
		err := os.Remove(fp.Name())
		if err != nil {
			t.Errorf("test.TempFile: cannot remove %#v: %v", fp.Name(), err)
		}
	}
}

// NormalizeIndent removes tab indentation from every line.
//
// This is useful for "inline" multiline strings:
//
//   cases := []struct {
//       string in
//   }{
//       `
//	 	    Hello,
//	 	    world!
//       `,
//   }
//
// This is nice and readable, but the downside is that every line will now have
// two extra tabs. This will remove those two tabs from every line.
//
// The amount of tabs to remove is based only on the first line, any further
// tabs will be preserved.
func NormalizeIndent(in string) string {
	indent := 0
	for _, c := range strings.TrimLeft(in, "\n") {
		if c != '\t' {
			break
		}
		indent++
	}

	r := ""
	for _, line := range strings.Split(in, "\n") {
		r += strings.Replace(line, "\t", "", indent) + "\n"
	}

	return strings.TrimSpace(r)
}

// R recovers a panic and cals t.Fatal().
//
// This is useful especially in subtests when you want to run a top-level defer.
// Subtests are run in their own goroutine, so those aren't called on regular
// panics. For example:
//
//   func TestX(t *testing.T) {
//       clean := someSetup()
//       defer clean()
//
//       t.Run("sub", func(t *testing.T) {
//           panic("oh noes")
//       })
//   }
//
// The defer is never called here. To fix it, call this function in all
// subtests:
//
//   t.Run("sub", func(t *testing.T) {
//       defer test.R(t)
//       panic("oh noes")
//   })
//
// See: https://github.com/golang/go/issues/20394
func R(t *testing.T) {
	t.Helper()
	r := recover()
	if r != nil {
		t.Fatalf("panic recover: %v", r)
	}
}

// SP makes a new String Pointer.
func SP(s string) *string { return &s }

// I64P makes a new Int64 Pointer.
func I64P(i int64) *int64 { return &i }

var inlines map[string]struct {
	inlined bool
	line    string
}

// MustInline issues a t.Error() if the Go compiler doesn't report that this
// function can be inlined.
//
// The first argument must the the full package name (i.e. "zgo.at/zstd/zint"),
// and the rest are function names to test:
//
//   ParseUint128         Regular function
//   Uint128.IsZero       Method call
//   (*Uint128).Parse     Pointer method
//
// The results are cached, so running it multiple times is fine.
//
// Inspired by the test in cmd/compile/internal/gc/inl_test.go
func MustInline(t *testing.T, pkg string, funs ...string) {
	t.Helper()

	if inlines == nil {
		getInlines(t)
	}

	for _, f := range funs {
		f = pkg + " " + f
		l, ok := inlines[f]
		if !ok {
			t.Errorf("unknown function: %q", f)
		}
		if !l.inlined {
			t.Errorf(l.line)
		}
	}
}

func getInlines(t *testing.T) {
	out, err := exec.Command("go", "list", "-f={{.Module.Path}}|{{.Module.Dir}}").CombinedOutput()
	if err != nil {
		t.Errorf("ztest.MustInline: %s\n%s", err, string(out))
		return
	}
	out = out[:len(out)-1]
	i := bytes.IndexRune(out, '|')
	mod, dir := string(out[:i]), string(out[i+1:])

	cmd := exec.Command("go", "build", "-gcflags=-m -m", mod+"/...")
	cmd.Dir = dir
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Errorf("ztest.MustInline: %s\n%s", err, string(out))
		return
	}

	inlines = make(map[string]struct {
		inlined bool
		line    string
	})

	var pkg string

	add := func(line string, i int, in bool) {
		fname := strings.TrimSuffix(line[i:i+strings.IndexRune(line[i:], ':')], " as")
		inlines[pkg+" "+fname] = struct {
			inlined bool
			line    string
		}{in, mod + line}
	}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "# ") {
			pkg = line[2:]
		}
		if i := strings.Index(line, ": can inline "); i > -1 {
			add(line, i+13, true)
		}
		if i := strings.Index(line, ": inline call to "); i > -1 {
			add(line, i+17, true)
		}
		if i := strings.Index(line, ": cannot inline "); i > -1 {
			add(line, i+16, false)
		}
	}
}
