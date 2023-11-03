// Package zstring implements functions for strings.
//
// All functions work correctly on Unicode codepoints/runes, but usually *don't*
// work on unicode clusters. That is, things like emojis composed of multiple
// codepoints and combining characters aren't dealt with unless explicitly
// mentioned otherwise.
package zstring

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Reverse a string.
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// Fields slices s to all substrings separated by sep. Leading/trailing
// whitespace and empty elements will be removed.
//
// e.g. "a;b", "a; b", "  a  ; b", and "a; b;" will all result in ["a", "b"].
func Fields(s, sep string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	f := strings.Split(s, sep)
	var rm []int
	for i := range f {
		f[i] = strings.TrimSpace(f[i])
		if f[i] == "" {
			rm = append(rm, i)
		}
	}
	for _, i := range rm {
		f = append(f[:i], f[i+1:]...)
	}
	return f
}

// Sub returns a substring starting at start and ending at end.
//
// Unlike regular string slicing this operates on runes/UTF-8 codepoints, rather
// than bytes.
func Sub(s string, start, end int) string {
	var (
		nchars    int
		startbyte = -1
	)
	for bytei := range s {
		if nchars == start {
			startbyte = bytei
		}
		if nchars == end {
			return s[startbyte:bytei]
		}
		nchars++
	}
	if startbyte == -1 {
		return ""
	}
	return s[startbyte:]
}

// ElideLeft returns the "n" left characters of the string.
//
// If the string is shorter than "n" it will return the first "n" characters of
// the string with "…" appended. Otherwise the entire string is returned as-is.
func ElideLeft(s string, n int) string {
	ss := Sub(s, 0, n)
	if len(s) != len(ss) {
		return ss + "…"
	}
	return s
}

// ElideRight returns the "n" right characters of the string.
//
// If the string is shorter than "n" it will return the first "n" characters of
// the string with "…" appended. Otherwise the entire string is returned as-is.
func ElideRight(s string, n int) string {
	ss := Sub(Reverse(s), 0, n)
	if len(s) != len(ss) {
		return "…" + Reverse(ss)
	}
	return s
}

// ElideCenter returns the "n" characters of the string.
//
// If the string is shorter than "n" it will return the first n/2 characters and
// last n/2 characters of the string with "…" inserted in the centre. Otherwise
// the entire string is returned as-is.
func ElideCenter(s string, n int) string {
	cc := utf8.RuneCountInString(s)
	if n >= cc {
		return s
	}

	var start string
	if n%2 == 0 {
		start = Sub(s, 0, n/2)
	} else {
		start = Sub(s, 0, n/2+1)
	}
	return start + "…" + Sub(s, cc-n/2, cc)
}

// UpperFirst transforms the first character to upper case, leaving the rest of
// the casing alone.
func UpperFirst(s string) string {
	if len(s) <= 1 {
		return strings.ToUpper(s)
	}
	for _, c := range s {
		sc := string(c)
		return strings.ToUpper(sc) + s[len(sc):]
	}
	return ""
}

// LowerFirst transforms the first character to lower case, leaving the rest of
// the casing alone.
func LowerFirst(s string) string {
	if len(s) <= 1 {
		return strings.ToLower(s)
	}
	for _, c := range s {
		sc := string(c)
		return strings.ToLower(sc) + s[len(sc):]
	}
	return ""
}

// GetLine gets the nth line \n-denoted line from a string.
//
// Line indexing starts at 1.
func GetLine(in string, n int) string {
	// Would probably be faster to use []byte and find the Nth \n character, but
	// this is "fast enough"™ for now.
	arr := strings.SplitN(in, "\n", n+1)
	if len(arr) <= n-1 {
		return ""
	}
	return arr[n-1]
}

// Filter a list.
//
// The function will be called for every item and those that return false will
// not be included in the return value.
func Filter(list []string, fun func(string) bool) []string {
	var ret []string
	for _, e := range list {
		if fun(e) {
			ret = append(ret, e)
		}
	}

	return ret
}

// FilterEmpty is a filter for Filter() to remove empty entries.
//
// An entry is considered "empty" if it's "" or contains only whitespace.
func FilterEmpty(e string) bool { return strings.TrimSpace(e) != "" }

// AlignLeft left-aligns a string, filling up any remaining width with spaces.
func AlignLeft(s string, n int) string {
	l := utf8.RuneCountInString(s)
	if l >= n {
		return s
	}
	return s + strings.Repeat(" ", n-l)
}

// AlignRight right-aligns a string, filling up any remaining width with spaces.
func AlignRight(s string, n int) string {
	l := utf8.RuneCountInString(s)
	if l >= n {
		return s
	}
	return strings.Repeat(" ", n-l) + s
}

// AlignCenter centre-aligns a string, filling up any remaining width with spaces.
func AlignCenter(s string, n int) string {
	if s == "" {
		return strings.Repeat(" ", n)
	}

	l := utf8.RuneCountInString(s)
	if l >= n {
		return s
	}

	pad := strings.Repeat(" ", (n-l)/2)
	if n%2 == 0 {
		return pad + s + pad + " "
	}
	return pad + s + pad
}

// Upto slices the string up to the first occurrence of sep. This is a shortcut
// for:
//
//	if i := strings.Index(s, sep); i > -1 {
//	  s = s[:i]
//	}
func Upto(s string, sep string) string {
	i := strings.Index(s, sep)
	if i == -1 {
		return s
	}
	return s[:i]
}

// From slices the string from first occurrence of sep. This is a shortcut for:
//
//	if i := strings.Index(s, sep); i > -1 {
//	  s = s[i+len(sep):]
//	}
func From(s string, sep string) string {
	i := strings.Index(s, sep)
	if i == -1 {
		return s
	}
	return s[i+len(sep):]
}

// IndexN finds the nth occurrence of a string.
//
// n starts at 1; returns -1 if there is no nth occurrence of this string.
func IndexN(s, find string, n uint) int {
	if s == "" || find == "" {
		return -1
	}
	if n == 0 {
		n = 1
	}
	n--

	var (
		off    int
		nfound uint
	)
	for i := strings.Index(s[off:], find); i != -1; i = strings.Index(s[off:], find) {
		if nfound == n {
			return off + i
		}
		nfound++
		off += i + len(find)
	}
	return -1
}

// IndexAll finds all occurrences of the string "find".
func IndexAll(s, find string) []int {
	if s == "" || find == "" {
		return nil
	}
	var (
		found = make([]int, 0, 2)
		pos   int
	)
	for {
		p := strings.Index(s[pos:], find)
		if p == -1 {
			break
		}
		found = append(found, pos+p)
		pos += p + 1
	}
	return found
}

// IndexPairs finds the position of all start/end pairs.
//
// Nested pairs are not supported.
//
// The return value is from last match to first match; this makes it easier to
// manipulate the string based on the indexes.
func IndexPairs(str, start, end string) [][]int {
	r := make([][]int, 0, 4)

	var pos int
	for {
		s := strings.Index(str[pos:], start)
		if s == -1 {
			break
		}
		e := strings.Index(str[pos+s:], end)
		if e == -1 {
			break
		}

		r = append(r, []int{pos + s, pos + s + e})
		pos = pos + s + e
	}
	if len(r) == 0 {
		return nil
	}

	for i := len(r)/2 - 1; i >= 0; i-- {
		opp := len(r) - 1 - i
		r[i], r[opp] = r[opp], r[i]
	}
	return r
}

// ReplacePairs replaces everything starting with start and ending with end with
// the return value of the callback.
func ReplacePairs(str, start, end string, f func(int, string) string) string {
	pairs := IndexPairs(str, start, end)
	for i, m := range pairs {
		if m[0] > 0 && str[m[0]] == str[m[0]-1] {
			str = str[:m[0]] + str[m[0]+1:]
			continue
		}
		str = str[:m[0]] + f(i, str[m[0]:m[1]+1]) + str[m[1]+1:]
	}
	return str
}

// Ident adds n spaces of indentation to every line.
func Indent(s string, n int) string {
	var (
		buf    strings.Builder
		indent = strings.Repeat(" ", n)
	)
	buf.Grow(len(s) + n*2)
	buf.WriteString(indent)
	for _, c := range s {
		buf.WriteRune(c)
		if c == '\n' {
			buf.WriteString(indent)
		}
	}

	// TODO: may be faster with bytes.Buffer? Can set the length on that.
	if s[len(s)-1] == '\n' {
		s = buf.String()
		return s[:len(s)-len(indent)]
	}
	return buf.String()
}

// HasSuffixes tests whether the string s ends with any of the suffixes.
//
// Identical to:
//
//	strings.HasSuffix(s, "one") || strings.HasSuffix(s, "two")
func HasSuffixes(s string, suffixes ...string) bool {
	for _, suf := range suffixes {
		h := strings.HasSuffix(s, suf)
		if h {
			return true
		}
	}
	return false
}

// HasPrefixes tests whether the string s starts with any of the prefixes.
//
// Identical to:
//
//	strings.HasPrefix(s, "one") || strings.HasPrefix(s, "two")
func HasPrefixes(s string, prefixes ...string) bool {
	for _, pre := range prefixes {
		h := strings.HasPrefix(s, pre)
		if h {
			return true
		}
	}
	return false
}

// TrimSuffixes returns s without the provided trailing suffixes strings.
//
// Identical to:
//
//	s = strings.TrimSuffix(s, "one")
//	s = strings.TrimSuffix(s, "two")
func TrimSuffixes(s string, suffixes ...string) string {
	for _, suf := range suffixes {
		s = strings.TrimSuffix(s, suf)
	}
	return s
}

// TrimPrefixes returns s without the provided leading prefixes strings.
//
// Identical to:
//
//	s = strings.TrimPrefix(s, "one")
//	s = strings.TrimPrefix(s, "two")
func TrimPrefixes(s string, prefixes ...string) string {
	for _, pre := range prefixes {
		s = strings.TrimPrefix(s, pre)
	}
	return s
}

// IsASCII reports if this string looks like it's plain 7-bit ASCII.
func IsASCII(s string) bool {
	for _, c := range s {
		if c > 0x7f {
			return false
		}
	}
	return true
}

// HasUpper reports if s has at least one upper-case character.
func HasUpper(s string) bool {
	for _, c := range s {
		if unicode.IsUpper(c) {
			return true
		}
	}
	return false
}

// Unwrap a string: single newlines become a space, whereas two or more are
// preserved.
//
// Removes newlines at the start and end of the string, but leaves all other
// spacing intact (including before and after newlines).
func Unwrap(s string) string {
	var (
		b bytes.Buffer
		n = 0
	)
	b.Grow(len(s))
	for _, c := range s {
		if c == '\n' {
			n++
			continue
		}
		if b.Len() > 0 {
			if n == 1 {
				b.WriteByte(' ')
			} else if n >= 2 {
				b.WriteByte('\n')
				b.WriteByte('\n')
			}
		}
		n = 0
		b.WriteRune(c)
	}
	return b.String()
}
