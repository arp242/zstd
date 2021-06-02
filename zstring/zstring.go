// Package zstring implements functions for strings.
//
// All functions work correctly on Unicode codepoints/runes, but usually *don't*
// work on unicode clusters. That is, things like emojis composed of multiple
// codepoints and combining characters aren't dealt with unless explicitly
// mentioned otherwise.
package zstring

import (
	"bytes"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
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

// ContainsAny reports whether any of the strings are in the list
func ContainsAny(list []string, strs ...string) bool {
	for _, s := range strs {
		if Contains(list, s) {
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

const nbsp = 0xa0

// WordWrap word wraps at n columns and prefixes subsequent lines with "prefix".
//
// Note the prefix is excluded from the length calculations; so if you want to
// wrap at 80 with a prefix of "> ", then you should wrap at 78.
//
// Adapted from: https://github.com/mitchellh/go-wordwrap
func WordWrap(text, prefix string, lim int) string {
	var (
		init                             = make([]byte, 0, len(text))
		buf                              = bytes.NewBuffer(init)
		wordBuf, spaceBuf                bytes.Buffer
		current, wordBufLen, spaceBufLen int
	)
	for _, char := range text {
		switch {
		case char == '\n':
			if wordBuf.Len() == 0 {
				if current+spaceBufLen > lim {
					current = 0
				} else {
					current += spaceBufLen
					spaceBuf.WriteTo(buf)
				}
				spaceBuf.Reset()
				spaceBufLen = 0
			} else {
				current += spaceBufLen + wordBufLen
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				spaceBufLen = 0
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
				wordBufLen = 0
			}
			buf.WriteRune(char)
			current = 0
		case unicode.IsSpace(char) && char != nbsp:
			if spaceBuf.Len() == 0 || wordBuf.Len() > 0 {
				current += spaceBufLen + wordBufLen
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				spaceBufLen = 0
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
				wordBufLen = 0
			}

			spaceBuf.WriteRune(char)
			spaceBufLen++
		default:
			wordBuf.WriteRune(char)
			wordBufLen++

			if current+wordBufLen+spaceBufLen > lim && wordBufLen < lim {
				buf.WriteRune('\n')
				buf.WriteString(prefix)
				current = 0
				spaceBuf.Reset()
				spaceBufLen = 0
			}
		}
	}

	if wordBuf.Len() == 0 {
		if current+spaceBufLen <= lim {
			spaceBuf.WriteTo(buf)
		}
	} else {
		spaceBuf.WriteTo(buf)
		wordBuf.WriteTo(buf)
	}

	return buf.String()
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

// DisplayWidth gets the display width of a string, taking tabs and escape
// sequences in to account.
//
// This does *not* handle various unicode aspects (i.e. graphmeme clusters,
// display width).
func DisplayWidth(s string) int {
	l := utf8.RuneCountInString(s)

	// Tabs are not a fixed width, but go to the nearest multiple of 8.
	split := strings.Split(s, "\t")
	for _, ss := range split[:len(split)-1] {
		l += 7 - utf8.RuneCountInString(ss)
	}

	// Escape sequences.
	for _, esc := range IndexAll(s, "\x1b") {
		i := 1
		for _, c := range s[esc:] {
			if c == 'm' {
				break
			}
			i++
		}
		l -= i
	}

	// TODO: Maybe also find a list of common unprintable things?

	return l
}

// HasSuffixes tests whether the string s ends with any of the suffixes.
//
// Identical to:
//
//   strings.HasSuffix(s, "one") || strings.HasSuffix(s, "two")
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
//   strings.HasPrefix(s, "one") || strings.HasPrefix(s, "two")
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
//   s = strings.TrimSuffix(s, "one")
//   s = strings.TrimSuffix(s, "two")
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
//   s = strings.TrimPrefix(s, "one")
//   s = strings.TrimPrefix(s, "two")
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

// Remove all values from a list.
//
// The return value indicates if this value was found at all.
func Remove(l *[]string, name string) bool {
	found := false
	ll := *l
	for i := len(ll) - 1; i >= 0; i-- {
		if ll[i] == name {
			ll = append(ll[:i], ll[i+1:]...)
			found = true
		}
	}
	*l = ll
	return found
}

// String converts a value to a string.
//
// This works for all built-in primitives, []byte, and values that implement
// fmt.Stringer.
func String(v interface{}) string {
	if v == nil {
		return ""
	}

	// Using a type switch like this is a bit ugly, but it avoids allocs and
	// seems to be the fastest (reflect isn't *that* much slower, but it's a bit
	// slower and can panic).
	switch vv := v.(type) {
	default:
		return "<zstring.ToString: unsupported type: " + reflect.TypeOf(v).String() + ">"

	case interface{ String() string }:
		return vv.String()

	case string:
		return vv
	case *string:
		return *vv
	case []byte:
		return string(vv)
	case *[]byte:
		return string(*vv)

	case bool:
		return strconv.FormatBool(vv)
	case *bool:
		return strconv.FormatBool(*vv)

	case float32:
		return strconv.FormatFloat(float64(vv), 'f', -1, 32)
	case *float32:
		return strconv.FormatFloat(float64(*vv), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(vv, 'f', -1, 64)
	case *float64:
		return strconv.FormatFloat(*vv, 'f', -1, 64)

	case complex64:
		return strconv.FormatComplex(complex128(vv), 'f', -1, 64)
	case *complex64:
		return strconv.FormatComplex(complex128(*vv), 'f', -1, 64)
	case complex128:
		return strconv.FormatComplex(vv, 'f', -1, 64)
	case *complex128:
		return strconv.FormatComplex(*vv, 'f', -1, 64)

	case int:
		return strconv.FormatInt(int64(vv), 10)
	case *int:
		return strconv.FormatInt(int64(*vv), 10)
	case int8:
		return strconv.FormatInt(int64(vv), 10)
	case *int8:
		return strconv.FormatInt(int64(*vv), 10)
	case int16:
		return strconv.FormatInt(int64(vv), 10)
	case *int16:
		return strconv.FormatInt(int64(*vv), 10)
	case int32:
		return strconv.FormatInt(int64(vv), 10)
	case *int32:
		return strconv.FormatInt(int64(*vv), 10)
	case int64:
		return strconv.FormatInt(vv, 10)
	case *int64:
		return strconv.FormatInt(*vv, 10)

	case uint:
		return strconv.FormatUint(uint64(vv), 10)
	case *uint:
		return strconv.FormatUint(uint64(*vv), 10)
	case uint8:
		return strconv.FormatUint(uint64(vv), 10)
	case *uint8:
		return strconv.FormatUint(uint64(*vv), 10)
	case uint16:
		return strconv.FormatUint(uint64(vv), 10)
	case *uint16:
		return strconv.FormatUint(uint64(*vv), 10)
	case uint32:
		return strconv.FormatUint(uint64(vv), 10)
	case *uint32:
		return strconv.FormatUint(uint64(*vv), 10)
	case uint64:
		return strconv.FormatUint(vv, 10)
	case *uint64:
		return strconv.FormatUint(*vv, 10)
	}
}
