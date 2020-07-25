// Package ztime implements functions for date and time.
package ztime

import (
	"fmt"
	"time"
)

var (
	day    = 24 * time.Hour
	day100 = 24 * time.Hour * 100
)

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
