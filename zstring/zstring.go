// Package zstring implements functions for strings.
//
// All functions work correctly on UTF-8 characters unless mentioned otherwise.
package zstring

import (
	"math/rand"
	"sort"
	"strings"
	"time"
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
// Unlike regular string slicing this operates on runes/UTF-8 characters, rather
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
func GetLine(in string, n int) string {
	// Would probably be faster to use []byte and find the Nth \n character, but
	// this is "fast enough"™ for now.
	arr := strings.SplitN(in, "\n", n+1)
	if len(arr) <= n-1 {
		return ""
	}
	return arr[n-1]
}

// Uniq removes duplicate entries from list; the list will be sorted.
func Uniq(list []string) []string {
	sort.Strings(list)
	var last string
	l := list[:0]
	for _, str := range list {
		if str != last {
			l = append(l, str)
		}
		last = str
	}
	return l
}

// Contains reports whether str is within the list.
func Contains(list []string, str string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}

// Repeat returns a slice with the string s repeated n times.
func Repeat(s string, n int) (r []string) {
	for i := 0; i < n; i++ {
		r = append(r, s)
	}
	return r
}

// Choose chooses a random item from the list.
func Choose(l []string) string {
	if len(l) == 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	return l[rand.Intn(len(l))]
}

// Filter a list. The function will be called for every item and those that
// return false will not be included in the return value.
func Filter(list []string, fun func(string) bool) []string {
	var ret []string
	for _, e := range list {
		if fun(e) {
			ret = append(ret, e)
		}
	}

	return ret
}

// FilterEmpty can be used as an argument for Filter() and will return false if
// e is empty or contains only whitespace.
func FilterEmpty(e string) bool { return strings.TrimSpace(e) != "" }

// Difference returns a new slice with elements that are in "set" but not in
// "others".
func Difference(set []string, others ...[]string) []string {
	out := []string{}
	for _, setItem := range set {
		found := false
		for _, o := range others {
			if Contains(o, setItem) {
				found = true
				break
			}
		}

		if !found {
			out = append(out, setItem)
		}
	}

	return out
}

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

// Split2 splits a string with strings.SplitN(.., 2) and returns the result.
//
// This makes some string splits a bit more elegant:
//
//   key, value := zstring.Split2(line, "=")
func Split2(str, sep string) (string, string) {
	s := strings.SplitN(str, sep, 2)
	if len(s) == 1 {
		return s[0], ""
	}
	return s[0], s[1]
}

// Split3 splits a string with strings.SplitN(.., 3) and returns the result.
func Split3(str, sep string) (string, string, string) {
	s := strings.SplitN(str, sep, 3)
	if len(s) < 3 {
		m := make([]string, 3)
		copy(m, s)
		s = m
	}
	return s[0], s[1], s[2]
}

// Split4 splits a string with strings.SplitN(.., 4) and returns the result.
func Split4(str, sep string) (string, string, string, string) {
	s := strings.SplitN(str, sep, 4)
	if len(s) < 4 {
		m := make([]string, 4)
		copy(m, s)
		s = m
	}
	return s[0], s[1], s[2], s[3]
}

// Upto slices the string up to the first occurrence of sep. This is a shortcut
// for:
//
//   if i := strings.Index(s, sep); i > -1 {
//     s = s[:i]
//   }
func Upto(s string, sep string) string {
	i := strings.Index(s, sep)
	if i == -1 {
		return s
	}
	return s[:i]
}

// From slices the string from first occurrence of sep. This is a shortcut for:
//
//   if i := strings.Index(s, sep); i > -1 {
//     s = s[i+len(sep):]
//   }
func From(s string, sep string) string {
	i := strings.Index(s, sep)
	if i == -1 {
		return s
	}
	return s[i+len(sep):]
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

// Tabwidth gets the display width of a string, taking tabs in to account.
//
// This does *not* handle various unicode aspects (i.e. graphmeme clusters,
// display width).
func TabWidth(s string) int {
	l := utf8.RuneCountInString(s)
	if strings.Index(s, "\t") == -1 {
		return l
	}

	split := strings.Split(s, "\t")
	for _, ss := range split[:len(split)-1] {
		l += 7 - utf8.RuneCountInString(ss)
	}
	return l
}
