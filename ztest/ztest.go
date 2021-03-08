// Package ztest contains helper functions that are useful for writing tests.
package ztest

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"zgo.at/zstd/zint"
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

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYz"

// Replace pieces of text with a placeholder string.
//
// This is use to test output which isn't stable, for example because it
// contains times:
//
//   ztest.Replace("Time: 1161 seconds", `Time: (\d+) s`)
//
// Will result in "Time: AAAA seconds".
//
// The number of replacement characters is equal to the input, unless the
// pattern contains "+" or "*" in which case it's always replaced by three
// characters.
func Replace(s string, patt ...string) string {
	type x struct {
		start, end int
		varWidth   bool
	}
	var where []x

	// Collect what to replace first so we can order things sensibly from A → B
	// → C → D, etc.
	for _, p := range patt {
		varWidth := false
		if i := strings.IndexAny(p, "+*"); i >= 0 {
			varWidth = i == 0 || p[i-1] != '\\'
		}

		for _, m := range regexp.MustCompile(p).FindAllStringSubmatchIndex(s, -1) {
			off := 2
			if len(m) == 2 { // No groups, replace everything.
				off = 0
			}

			for i := off; len(m) > i; i += 2 {
				where = append(where, x{
					start:    m[i],
					end:      m[i+1],
					varWidth: varWidth,
				})
			}
		}
	}

	sort.Slice(where, func(i, j int) bool { return where[i].start > where[j].start })
	charidx := zint.MinInt(len(where), len(chars)) - 1
	for _, w := range where {
		l := 3
		if !w.varWidth {
			l = w.end - w.start
		}

		s = s[:w.start] + strings.Repeat(string(chars[charidx]), l) + s[w.end:]
		charidx--
		if charidx < 0 {
			charidx = len(chars) - 1
		}
	}
	return s
}

// Read data from a file.
func Read(t *testing.T, paths ...string) []byte {
	t.Helper()

	path := filepath.Join(paths...)
	file, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ztest.Read: cannot read %v: %v", path, err)
	}
	return file
}

// TempFile creates a new temporary file and returns the path.
func TempFile(t *testing.T, data string) string {
	t.Helper()

	tmpdir := t.TempDir()
	fp, err := os.CreateTemp(tmpdir, "ztest.*")
	if err != nil {
		t.Fatalf("ztest.TempFile: could not create file in %v: %v", tmpdir, err)
	}

	defer func() {
		err := fp.Close()
		if err != nil {
			t.Fatalf("ztest.TempFile: close: %v", err)
		}
	}()

	_, err = fp.WriteString(data)
	if err != nil {
		t.Fatalf("ztest.TempFile: write: %v", err)
	}

	t.Cleanup(func() {
		err := os.Remove(fp.Name())
		if err != nil {
			t.Errorf("ztest.TempFile: cannot remove %#v: %v", fp.Name(), err)
		}
	})

	return fp.Name()
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
