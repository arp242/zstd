// Package ztime implements functions for date and time.
package ztime

import (
	"fmt"
	"math"
	"strings"
	"time"
)

var (
	day    = 24 * time.Hour
	day100 = 24 * time.Hour * 100
)

// Takes records how long the execution of f takes.
func Takes(f func()) time.Duration {
	s := time.Now()
	f()
	return time.Now().Sub(s)
}

// DurationAs formats a duration as the given time unit.
//
// Use Round() if you want to limit the precision; for example:
//
//   DurationAs(d.Round(time.Microsecond), time.Millisecond)
func DurationAs(d, as time.Duration) string {
	f := float64(d) / float64(as)
	_, frac := math.Modf(f)
	if frac == 0 {
		return fmt.Sprintf("%.0f", f)
	}

	s := strings.TrimRight(fmt.Sprintf("%f", f), "0")
	if s == "0." {
		s = fmt.Sprintf("%.999f", f)
		for i, c := range s {
			if c != '0' && c != '.' {
				s = s[:i+1]
				break
			}
		}
	}
	return s
}

// FormatDuration formats duration as a concise string. It's less accurate than
// Duration.String(), but shorter.
func FormatDuration(d time.Duration) string {
	switch {
	case d >= day100:
		h := int(d.Round(time.Hour).Hours())
		return fmt.Sprintf("%dd", h/24)
	case d > day:
		h := int(d.Round(time.Hour).Hours())
		return fmt.Sprintf("%dd%dh", h/24, h%24)
	case d >= 10*time.Minute:
		s := d.Round(time.Minute).String()
		return s[:len(s)-2]
	}

	return d.Round(time.Second).String()
}

// LeapYear reports if this year is a leap year according to the Gregorian
// calendar.
func LeapYear(t time.Time) bool {
	y := t.Year()
	return y%4 == 0 && (y%100 != 0 || y%400 == 0)
}

// DaysInMonth gets the number of days for the month.
func DaysInMonth(t time.Time) int {
	m := t.Month()
	if m == 2 && LeapYear(t) {
		return 29
	}
	switch m {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	default:
		return 30
	}
}

// LastInMonth reports if the current day is the last day in this month.
func LastInMonth(t time.Time) bool { return t.Day() == DaysInMonth(t) }

func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
